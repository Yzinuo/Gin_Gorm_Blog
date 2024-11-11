package handle

import (
	"encoding/json"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct{}

type UpdateCurrentUserReq struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Intro    string `json:"intro"`
	Website  string `json:"website"`
	Email    string `json:"email"`
}

type UpdateCurrentPasswordReq struct {
	NewPassword string `json:"new_password" binding:"required,min = 4,max =20"`
	OldPassword string `json:"old_password" binding:"required,min = 4,max =20"`
}

type UpdateUserReq struct {
	UserAuthId int    `json:"user_auth_id"`
	Nickname   string `json:"nickname"`
	RoleIds    []int  `json:"role_ids"`
}

type UpdateUserDisableReq struct {
	UserAuthId int  `json:"id"`
	IsDisable  bool `json:"is_disable"`
}

type UserQuery struct {
	PageQuery
	LoginType int8   `form:"login_type"`
	Username  string `form:"username"`
	Nickname  string `form:"nickname"`
}

type ForceOfflineReq struct {
	UserInfoId int `json:"user_info_id"`
}

// 从redis中获取用户点赞信息后，返回用户信息
func (*User) GetInfo(c *gin.Context) {
	user, err := CurrentUserAuth(c)
	if err != nil {
		ReturnError(c, g.ErrTokenRuntime, err)
		return
	}
	userInfo := model.UserInfoVO{UserInfo: *user.UserInfo}

	rdb := GetRDB(c)

	//从redis中获取用户点赞信息
	userInfo.ArticleLikeSet, err = rdb.SMembers(rdbctx, g.ARTICLE_USER_LIKE_SET+strconv.Itoa(user.ID)).Result()
	if err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}
	userInfo.CommentLikeSet, err = rdb.SMembers(rdbctx, g.COMMENT_USER_LIKE_SET+strconv.Itoa(user.ID)).Result()
	if err != nil {
		ReturnError(c, g.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, userInfo)
}

// 更新当前用户信息
func (*User) UpdateCurrent(c *gin.Context) {
	var req UpdateCurrentUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	auth, _ := CurrentUserAuth(c)
	err := model.UpdateUserInfo(GetDB(c), auth.UserInfoId, req.Nickname, req.Avatar, req.Intro, req.Website)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, nil)
}

// 更新用户信息  只能更新昵称+角色
func (*User) Update(c *gin.Context) {
	var req UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	if err := model.UpdateUserNicknameAndRole(GetDB(c), req.UserAuthId, req.Nickname, req.RoleIds); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, nil)
}

// 查询用户
func (*User) GetList(c *gin.Context) {
	var query UserQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	list, count, err := model.GetUserList(GetDB(c), query.LoginType, query.Page, query.Size, query.Nickname, query.Username)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, PageList[model.UserAuth]{
		Total: count,
		Data:  list,
		Size:  query.Size,
		Page:  query.Page,
	})
}

// 修改用户拉黑状态
func (*User) UpdateDisable(c *gin.Context) {
	var req UpdateUserDisableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	err := model.UpdateUserDisable(GetDB(c), req.UserAuthId, req.IsDisable)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, nil)
}

// 修改当前用户密码
// 首先先需要输入旧密码进行哈希认证，然后再修改密码
func (*User) UpdateCurrentPasswordReq(c *gin.Context) {
	var req UpdateCurrentPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	// 先判断旧密码是否正确
	auth, _ := CurrentUserAuth(c)
	if !utils.BcryptCheck(req.OldPassword, auth.Password) {
		ReturnError(c, g.ErrPassword, nil)
		return
	}

	// 再修改密码
	hashPassword, _ := utils.BcryptHash(req.NewPassword)
	if err := model.UpdateUserPassword(GetDB(c), auth.ID, hashPassword); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// 查询在线用户
// 先通过redis查询所有在线的用户的键表。 再循环查询用户的信息
// 根据上次登录时间顺序排序
func (*User) GetOnlineList(c *gin.Context) {
	keyword := c.Query("keyword")
	rdb := GetRDB(c)

	// 先查询所有在线用户的键表
	onlineList := make([]model.UserAuth, 0)
	keys := rdb.Keys(rdbctx, g.ONLINE_USER+"*").Val()

	for _, key := range keys {
		//查询用户信息，如果不包括keyword跳过
		var auth model.UserAuth
		userInfo := rdb.Get(rdbctx, key).Val() // redis中存储的是json数据
		json.Unmarshal([]byte(userInfo), &auth)

		if keyword != "" && !strings.Contains(auth.Username, keyword) && !strings.Contains(auth.UserInfo.Nickname, keyword) {
			continue
		}
		onlineList = append(onlineList, auth)
	}

	sort.Slice(onlineList, func(i, j int) bool {
		return onlineList[i].LastLoginTime.Unix() > onlineList[j].LastLoginTime.Unix()
	})

	ReturnSuccess(c, onlineList)
}

// 强制下线用户
// 不能离线自己。 从redis中删除在线用户的键表，添加到离线用户的键表
func (*User) ForceOffline(c *gin.Context) {
	id := c.Param("id")
	uid, err := strconv.Atoi(id)
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	auth, _ := CurrentUserAuth(c)
	if auth.ID == uid{
		ReturnError(c,g.ErrForceOfflineSelf,nil)
		return
	}

	rdb := GetRDB(c)

	// 先从redis中删除在线用户的键值
	rdb.Del(rdbctx, g.ONLINE_USER+strconv.Itoa(uid))

	// 再添加到离线用户的键值
	rdb.Set(rdbctx, g.OFFLINE_USER+strconv.Itoa(uid), uid, time.Hour)

	ReturnSuccess(c, "强制离线成功")
}
