// 负责获取博客的首页信息
package handle

import (
	"context"
	"gin-blog/internal/model"
	utils "gin-blog/internal/utils"
	"log/slog"
	"strings"

	g "gin-blog/internal/global"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

type BlogInfo struct{}

type BlogInfoHomeVO struct {
	ArticleCount  int `json:"article_count"`
	ViewCount  int `json:"view_count"`
	UserCount	  int `json:"user_count"`
	MessageCount  int `json:"message_count"`
}

type AboutReq struct {
	Content string `json:"content"`
}


func(*BlogInfo) GetConfigMap(c *gin.Context){
	db := GetDB(c)
	rdb := GetRDB(c)
	
	config,err := GetConfigCache(rdb)
	if err != nil {
		ReturnError(c,g.ErrRedisOp,err)
		return
	}

	if len(config) > 0 {
		slog.Debug("get config from redis cache")
		ReturnSuccess(c,config)
	}

	config,err = model.GetConfigMap(db)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	if err := AddConfigCache(rdb,config); err != nil {
		ReturnError(c,g.ErrRedisOp,err)
		return
	}

	ReturnSuccess(c,config)
}

func (*BlogInfo) UpdateBlogInfo(c *gin.Context ){
	var m map[string]string
	if err := c.ShouldBindJSON(&m);err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}
	
	err := model.UpdateConfigMap(GetDB(c),m)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	DelConfigCache(GetRDB(c))
	ReturnSuccess(c,nil)
}

func (* BlogInfo) GetBlogInfo (c *gin.Context){
	db := GetDB(c)
	rdb := GetRDB(c)

	articlecount,err := model.Count(db,&model.Article{},"status = ? AND is_delete = ?",1,0)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	usercount,err := model.Count(db,&model.UserInfo{})
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	messagecount,err := model.Count(db,&model.Message{})
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	viewcount,err  := rdb.Get(rdbctx,g.VIEW_COUNT).Int()
	if err != nil && err != redis.Nil{
		ReturnError(c,g.ErrRedisOp,err)
		return
	}

	ReturnSuccess(c,BlogInfoHomeVO{
		ArticleCount: int(articlecount),
		ViewCount: int(viewcount),
		UserCount: int(usercount),
		MessageCount: int(messagecount),
	})
}

func (*BlogInfo) GetAbout(c *gin.Context) {
	ReturnSuccess(c, model.GetValueByKey(GetDB(c), g.CONFIG_ABOUT))
}

func (*BlogInfo) UpdateAbout (c *gin.Context){
	var req AboutReq
	
	if err := c.ShouldBindJSON(&req); err !=nil{
		ReturnError(c,g.ErrRequest,err)
		return
	}
	
	if err := model.FindOrCreateConfig(GetDB(c),g.CONFIG_ABOUT,req.Content); err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	ReturnSuccess(c,nil)
}

// 用户登录 rdb记录信息
func (*BlogInfo) Report(c *gin.Context){
	rdb := GetRDB(c)
	
	ipaddress := utils.IP.GetIpaddress(c)
	useragent := utils.IP.GetUserAgent(c)
	brower := useragent.Name + " " + useragent.Version.String()
	OS := useragent.OS + " "+ useragent.OSVersion.String()

	uuid := utils.MD5(ipaddress+brower+OS)
	ctx := context.Background()

	// 判断有无记录过
	if !rdb.SIsMember(ctx,g.KEY_UNIQUE_VISITOR_SET,uuid).Val() {
		ipsource := utils.IP.GetIPsource(ipaddress)
		if ipsource != ""{
			source := strings.Split(ipsource,"|")
			rdb.HIncrBy(ctx,g.VISITOR_AREA,strings.ReplaceAll(source[2],"省",""),1)
		} else {
			rdb.HIncrBy(ctx,g.VISITOR_AREA,"未知",1)
		}
		// 记录一下为记录的用户
		rdb.Incr(ctx,g.VIEW_COUNT)
		rdb.SAdd(ctx,g.KEY_UNIQUE_VISITOR_SET,uuid)
	}

	ReturnSuccess(c,nil)
}
