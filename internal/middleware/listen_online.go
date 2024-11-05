package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	g "gin-blog/internal/global"
	"gin-blog/internal/handle"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

// 监听在线用户,刷新用户的在线持续时间
func ListenOnline() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx := context.Background()
		rdb := c.MustGet("rdb").(*redis.Client)

		auth,err := handle.CurrentUserAuth(c)
		if err != nil {
			handle.ReturnError(c,g.ErrUserAuth,err)
		}

		onlinekey := g.ONLINE_USER+strconv.Itoa(auth.ID)
		offlinekey := g.OFFLINE_USER+strconv.Itoa(auth.ID)
		
		// 查询是否在黑名单中
		if rdb.Exists(ctx,offlinekey).Val() == 1 {
			fmt.Print("用户已被强制下线,黑名单中")
			handle.ReturnError(c,g.ErrForceOffline,err)
			c.Abort()
			return
		}

		// 刷新用户的在线时间
		rdb.Set(ctx,onlinekey,auth.ID,10*time.Minute)
		c.Next()
	}
}