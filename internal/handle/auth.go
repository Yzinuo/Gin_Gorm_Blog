package handle

import (
	"errors"
	g "gin-blog/internal/global"
	model "gin-blog/internal/model"
	utils "gin-blog/internal/utils"
	"gin-blog/internal/utils/jwt"
	"log/slog"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Auth struct{}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=4,max=20"`
	Code     string `json:"code"  binding:"required"`
}

type LoginVO struct {
	model.UserInfo
	// 记录自己点赞的文章和评论
	ArticleLikeList []string `json:"article_like_list"`
	CommentLikeList []string `json:"comment_like_list"`
	Token           string   `json:"token"`
}

// 用户登陆
func (*Auth) Login(c *gin.Context) {
	var logreq LoginReq

	db := GetDB(c)
	rdb := GetRDB(c)

	if err := c.ShouldBind(&logreq); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	// 查询用户信息
	auth, err := model.GetUserAuthInfoByName(db, logreq.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrUserNotExist, err)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	//验证密码
	if !utils.BcryptCheck(logreq.Password, auth.Password) {
		ReturnError(c, g.ErrPassword, err)
		return
	}

	//更新登陆IP
	Ipaddress := utils.IP.GetIpaddress(c)
	IP := utils.IP.GetIPsourceSimpleInfo(Ipaddress)

	//获取userInfo （头像） AND 查询权限id
	userinfo, err := model.GetUserInfoById(db, auth.UserInfoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrUserNotExist, err)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	roleids, err := model.GetRoleIdsByUserId(db, auth.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, g.ErrUserNotExist, err)
			return
		}
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 从redis缓存中直接获取 点赞的文章和评论
	articleLikeSet, err := rdb.SMembers(rdbctx, g.ARTICLE_USER_LIKE_SET+strconv.Itoa(auth.ID)).Result()
	if err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}
	commentLikeSet, err := rdb.SMembers(rdbctx, g.COMMENT_USER_LIKE_SET+strconv.Itoa(auth.ID)).Result()
	if err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	//生成Token
	conf := g.Conf.JWT
	token, err := jwt.GenToken(conf.Secret, conf.Issuer, int(conf.Expire), auth.ID, roleids)
	if err != nil {
		ReturnError(c, g.ErrTokenCreate, err)
		return
	}

	// 更新用户登陆信息
	err = model.UpdateUserLoginInfo(db, auth.ID, Ipaddress, IP)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
	}

	// session更新用户信息
	slog.Info("登录成功 ！")
	session := sessions.Default(c)
	session.Set(g.CTX_USER_AUTH, auth.ID)
	session.Save()

	//为了防止在缓存中同时存在在线和强制下线两种状态，删除下线状态
	offlineuser := g.OFFLINE_USER + strconv.Itoa(auth.ID)
	rdb.Del(rdbctx, offlineuser).Result()

	ReturnSuccess(c, LoginVO{
		UserInfo:        *userinfo,
		ArticleLikeList: articleLikeSet,
		CommentLikeList: commentLikeSet,
		Token:           token,
	})

}

// 登出就是删掉context，redis，session的用户信息
func (*Auth) Logout(c *gin.Context) {
	c.Set(g.CTX_USER_AUTH, nil)

	auth, _ := CurrentUserAuth(c)
	//已经被删除
	if auth == nil {
		ReturnSuccess(c, nil)
		return
	}

	session := sessions.Default(c)
	session.Delete(g.CTX_USER_AUTH)
	session.Save()

	rdb := GetRDB(c)
	OnlineUser := g.ONLINE_USER + strconv.Itoa(auth.ID)
	rdb.Del(rdbctx, OnlineUser)

	ReturnSuccess(c, nil)
}
func (*Auth) Register(c *gin.Context) {
	ReturnSuccess(c, "注册")
}

func (*Auth) SendCode(c *gin.Context) {
	ReturnSuccess(c, "发送邮箱验证码")
}
