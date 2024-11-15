package handle

import (
	"errors"
	g "gin-blog/internal/global"
	model "gin-blog/internal/model"
	utils "gin-blog/internal/utils"
	"gin-blog/internal/utils/jwt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

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
	Username string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=4,max=20"`
	// Code     string `json:"code"  binding:"required"`
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
	var regreq RegisterReq
	if err := c.ShouldBindJSON(&regreq); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}
	regreq.Username = utils.Format(regreq.Username)

	// 检查用户名是否存在，避免重复注册
	auth,err := model.GetUserAuthInfoByName(GetDB(c),regreq.Username)
	if err != nil {
		var flag bool = false
		if errors.Is(err,gorm.ErrRecordNotFound) {
			flag = true
		}
		if !flag{
			ReturnError(c,g.ErrDbOp,err)
			return
		}
	}

	if auth != nil {
		ReturnError(c,g.ErrUserExist,err)
		return
	}
	

	// 通过邮箱验证
	info := utils.GenEmailVerificationInfo(regreq.Username,regreq.Password)
	SetMailInfo(GetRDB(c),info,15*time.Minute) // 15分钟过期
	EmailData := utils.GetEmailData(regreq.Username,info)
	err = utils.SendEmail(regreq.Username,EmailData)
	if err != nil {
		ReturnError(c,g.ErrSendEmail,err)
		return
	}

	ReturnSuccess(c,nil)
}

func (*Auth) VerifyCode(c *gin.Context) {
    var code string
    if code = c.Query("info"); code == "" {
        returnErrorPage(c)
        return
    }

    // 验证是否有code在数据库中
    ifExist, err := GetMailInfo(GetRDB(c), code)
    if err != nil {
        returnErrorPage(c)
        return
    }
    if !ifExist {
        returnErrorPage(c)
        return
    }

    DeleteMailInfo(GetRDB(c), code)

    username, password, err := utils.ParseEmailVerificationInfo(code)
    if err != nil {
        returnErrorPage(c)
        return
    }

    // 注册用户
      _,_,_,err = model.CreateNewUser(GetDB(c), username, password)
    if err != nil {
        returnErrorPage(c)
        return
    }

    // 注册成功，返回成功页面
    c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
        <!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>注册成功</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f4f4f4;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: #fff;
                    padding: 20px;
                    border-radius: 8px;
                    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    text-align: center;
                }
                h1 {
                    color: #5cb85c;
                }
                p {
                    color: #333;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>注册成功</h1>
                <p>恭喜您，注册成功！</p>
            </div>
        </body>
        </html>
    `))
}

func returnErrorPage(c *gin.Context) {
    c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`
        <!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>注册失败</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f4f4f4;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: #fff;
                    padding: 20px;
                    border-radius: 8px;
                    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    text-align: center;
                }
                h1 {
                    color: #d9534f;
                }
                p {
                    color: #333;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>注册失败</h1>
                <p>请重试。</p>
            </div>
        </body>
        </html>
    `))
}