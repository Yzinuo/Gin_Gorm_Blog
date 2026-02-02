// 负责获取博客的首页信息
package handle

import (
	"context"
	"gin-blog/internal/model"
	utils "gin-blog/internal/utils"
	"log/slog"
	"strconv"
	"strings"
	"time"

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
		normalizeWebsiteRuntime(config)
		ReturnSuccess(c,config)
		return
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

	normalizeWebsiteRuntime(config)
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

func normalizeWebsiteRuntime(config map[string]string) {
	const key = "website_createtime"
	value := strings.TrimSpace(config[key])
	if value == "" {
		return
	}

	createTime, ok := parseFlexibleTime(value)
	if !ok {
		return
	}

	now := time.Now()
	if createTime.After(now) {
		config["website_runtime_seconds"] = "0"
		config["website_runtime_days"] = "0"
	} else {
		seconds := int64(now.Sub(createTime).Seconds())
		config["website_runtime_seconds"] = strconv.FormatInt(seconds, 10)
		config["website_runtime_days"] = strconv.FormatInt(seconds/86400, 10)
	}

	config["website_createtime_unix"] = strconv.FormatInt(createTime.Unix(), 10)
	config["website_createtime_rfc3339"] = createTime.Format(time.RFC3339)
}

func parseFlexibleTime(value string) (time.Time, bool) {
	if value == "" {
		return time.Time{}, false
	}

	if isDigits(value) {
		if ts, err := strconv.ParseInt(value, 10, 64); err == nil {
			switch {
			case len(value) >= 13:
				return time.UnixMilli(ts), true
			case len(value) >= 10:
				return time.Unix(ts, 0), true
			}
		}
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02",
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return t, true
		}
	}

	return time.Time{}, false
}

func isDigits(value string) bool {
	for i := 0; i < len(value); i++ {
		c := value[i]
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(value) > 0
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
		rdb.SAdd(ctx,g.KEY_UNIQUE_VISITOR_SET,uuid)
	}

	ReturnSuccess(c,nil)
}
