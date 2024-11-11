package middleware

import (
	g "gin-blog/internal/global"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

// 把gorm.DB 注入到gin.Context中
// 后续开发只需要 c.MustGet(g.CTX_DB).(*gorm.DB) 即可拿到db
func WithGormDB(db *gorm.DB) gin.HandlerFunc{
	return func(c *gin.Context){
		c.Set(g.CTX_DB,db)
		c.Next()
	}
}

// 把redis.Client 注入到gin.Context中
// 后续开发只需要 c.MustGet(g.CTX_RDB).(*redis.Client) 即可拿到redis
func WithRDB( rdb *redis.Client) gin.HandlerFunc{
	return func (c *gin.Context){
		c.Set(g.CTX_RDB,rdb)
		c.Next()
	}
}

// 设置存储session的位置
func WithCookieStore(name,secret string) gin.HandlerFunc{
	store := cookie.NewStore([]byte(secret)) //使用密钥加密和解密session id
	store.Options(sessions.Options{Path:"/",MaxAge: 600}) // 所有域名下的都可以访问,同时设置session的有效期为10分钟 
	return sessions.Sessions(name,store)
}

func WithMemStore(name,secret string) gin.HandlerFunc{
	store := memstore.NewStore([]byte(secret))
	store.Options(sessions.Options{Path:"/",MaxAge: 600})
	return sessions.Sessions(name,store)
}


func CORS() gin.HandlerFunc{
	return cors.New(cors.Config{
		AllowOrigins: 				[]string{"*"},
		AllowMethods: 				[]string{"GET","POST","PUT","DELETE","OPTIONS","PATCH"},
		AllowHeaders:				[]string{"Origin","Authorization","Content-Type","X-Requested-With"},
		ExposeHeaders:              []string{"Content-Type"}, // 服务端返回数据后，允许客户端访问的响应头
		AllowCredentials: 			true, //允许发送cookie等验证信息
		AllowOriginFunc: 			func(origin string) bool{return true},
		MaxAge: 					24 * time.Hour, // 缓存预检请求，提高服务性能
	})
}
