package handle

import (
	"errors"
	g "gin-blog/internal/global"
	"log/slog"
	"net/http"

	model "gin-blog/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/gin-contrib/sessions"
	"gorm.io/gorm"
)

// 定义响应结构体
type Response [T any]struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data T 			 `json:"data"`
}

// HTTP code + code + msg + data
func ReturnHttpResponse  (c *gin.Context,httpcode int,code int, msg string,data any) {
	c.JSON(httpcode,Response[any]{
		Code : code,
		Msg  : msg,
		Data : data,
	})
}

// Result + data
func ReturnResponse  (c *gin.Context,r g.Result,data any){
	ReturnHttpResponse(c,http.StatusOK,r.GetCode(),r.GetMsg(),data)
}

// data 
func ReturnSuccess(c *gin.Context, data any) {
	ReturnResponse(c,g.OkReresult,data)
}
func ReturnFail (c *gin.Context,data any) {
	ReturnResponse(c,g.FailResult,data)
}

//预料中的错误  = 业务错误 + 系统错误  在业务层处理，返回200 http状态码
//意外的错误  = 触发panic， 在中间件中被捕获，返回500 http状态码
//data是错误数据（可以是error和string）， error是业务错误
func ReturnError (c *gin.Context,r g.Result,data any){
	slog.Info("[FUNC-RETURN-ERROR]] :" + r.Msg)

	var val string

	if data != nil {
		switch v := data.(type){
		case error:
			val = v.Error()
		
		case string:
			val = v
		}
	}

	c.AbortWithStatusJSON(
		http.StatusOK,
		Response[any]{
			Code : r.Code,
			Msg  : r.Msg,
			Data : val,
		},
	)
	
}


// 从上下文中获gorm.DB
func GetDB(c *gin.Context) *gorm.DB {
	return c.MustGet(g.CTX_DB).(*gorm.DB)
}

// 从上下文中获redis.Client
func GetRDB(c *gin.Context) *redis.Client {
	return c.MustGet(g.CTX_RDB).(*redis.Client)
}

// 分页获取数据
type PageList [T any] struct {
	Page  int 	 `json:"page"`
	Size  int 	 `json:"size"`
	Total int64  `json:"total"`
	Data  []T  	 `json:"page_data"`	
}

type PageQuery struct{
	Page  int `form:"page"`
	Size  int `form:"size"`
	Keyword string `form:"keyword"`
}

// 1. 从Context中获取UserInfo，如果获取到说明context中已经保存了
// 2. 从Session中获取UserInfo Uid(通过解析token获取)
// 3. 根据id从数据库中获取  然后在设置在gin context中
func CurrentUserAuth(c *gin.Context) (*model.UserAuth,error){
	// 1
	var key string =  g.CTX_USER_AUTH
	if cache, exist := c.Get(key); exist && cache != nil{
		slog.Debug("[FUNC-GET-USER-AUTH] : " + cache.(*model.UserAuth).Username)
		return cache.(*model.UserAuth),nil
	}
	
	// 2
	 s := sessions.Default(c)
	 id := s.Get(key)
	 
	// 3
	if id != nil{
		db := GetDB(c)
		userauth,err := model.GetUserAuthInfoById(db,id.(int))
		
		c.Set(key,userauth)
		if err != nil{
			return nil,err
		}
	}
	return nil,errors.New("session中没有 user_auth_id")
}
