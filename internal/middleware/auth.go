package middleware

import (
	"errors"
	"fmt"
	g "gin-blog/internal/global"
	"gin-blog/internal/handle"
	"gin-blog/internal/model"
	"gin-blog/internal/utils/jwt"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 利用JWT实现快速身份认证，实现API的权限控制 和 登录状态的记录(session 的用户数据由jwt生成，如果session中没有用户信息，在currentauth中算没有登录)
// 中间件为了实现延迟执行,所以一般都是闭包
func JWTAuth() gin.HandlerFunc{
	return func(c *gin.Context){
		slog.Debug("[middleware JWTAuth]]  JUST Do Jwt for API")
		db := c.MustGet(g.CTX_DB).(*gorm.DB)

		// 如果当前请求的方法受系统管理才需要做JWT认证,否则直接跳过认证
		// 不受系统管理的方法一般是:不依赖于服务器端框架或系统提供的管理功能，而是由开发者手动实现的方法
		url,method := c.FullPath()[4:],c.Request.Method
		resouce,err := model.GetResource(db,url,method)
		if err != nil{
			if errors.Is(err,gorm.ErrRecordNotFound){
				slog.Debug("[middleware JWTAuth]] 该请求不受系统管理")
				c.Set("skip_auth_Check",true)
				c.Next()
				c.Set("skip_auth_Check",false)
				return
			}
			handle.ReturnError(c,g.ErrDbOp,err)
			return
		}
		if resouce.Anonymous {
			slog.Debug("[middleware JWTAuth]] 该请求的方法 %s %s是匿名,无需jwt认证",method,url)
			c.Set("skip_auth",true)
			c.Next()
			c.Set("skip_auth",false)
		}

		// 获取Toekn
		authorization := c.Request.Header.Get("Authorization")
		if authorization == ""{
			handle.ReturnError(c,g.ErrTokenNotExist,nil)
			return
		}
		//正确的格式为: Bearer Token
		splitauth := strings.Split(authorization," ")
		//正常来说必须是2 出现异常现象 被攻击.
		if len(splitauth)!=2 || splitauth[0] != "Bearer"{
			handle.ReturnError(c,g.ErrTokenType ,nil)
			return
		}
		token := splitauth[1]
		
		//解析token ,从token中获取信息并设置
		claims,err := jwt.ParseToken(g.Conf.JWT.Secret,token)
		if err != nil {
			handle.ReturnError(c,g.ErrTokenWrong,nil)
			return
		}
		if time.Now().Unix() > claims.ExpiresAt.Unix(){
			handle.ReturnError(c,g.ErrTokenRuntime,nil)
			return
		}

		user,err := model.GetUserInfoById(db,claims.UserId)
		if err != nil {
			handle.ReturnError(c,g.ErrUserNotExist,nil)
			return
		}

		sess := sessions.Default(c)
		sess.Set(g.CTX_USER_AUTH ,user)
		sess.Save()

		c.Set(g.CTX_USER_AUTH,user)
	}
}

func PermissionCheck() gin.HandlerFunc{
	return func(c *gin.Context) {
		if c.GetBool("skip_auth_Check"){
			c.Next()
			return
		}

		db := c.MustGet(g.CTX_DB).(*gorm.DB)
		auth,_ := handle.CurrentUserAuth(c)
		if auth == nil{
			handle.ReturnError(c,g.ErrUserNotExist,nil)
			return
		}
		
		if auth.IsSuper{
			slog.Debug("[middleware PermissionCheck]] 该用户是超级管理员,无需权限验证")
			c.Next()
			return
		}


		url,method := c.FullPath()[4:],c.Request.Method
		slog.Debug(fmt.Sprintf("[middleware PermissionCheck]] 开始对%s url: %s  method:%s 进行权限认证"),
					auth.Username,url,method)
		for _,role := range auth.Roles{
			pass,err := model.CheckRoleAuth(db,url,method,role.ID)
			if err != nil {
				handle.ReturnError(c,g.ErrDbOp,err)
				return
			}

			if !pass {
				handle.ReturnError(c,g.ErrPermission ,nil)
				return
			}
		}
		
		slog.Debug("[middleware PermissionCheck]] 该用户有访问权限")
		c.Next()
	}
}
