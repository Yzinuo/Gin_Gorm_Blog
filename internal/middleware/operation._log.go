// gin 中获取回复内容,对于返回的内容进行包装处理
package middleware

import (
	"bytes"
	"gin-blog/internal/handle"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"
	"io"
	"log/slog"
	"strings"
	g "gin-blog/internal/global"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var optMap = map[string]string{
	"Article":      "文章",
	"BlogInfo":     "博客信息",
	"Category":     "分类",
	"Comment":      "评论",
	"FriendLink":   "友链",
	"Menu":         "菜单",
	"Message":      "留言",
	"OperationLog": "操作日志",
	"Resource":     "资源权限",
	"Role":         "角色",
	"Tag":          "标签",
	"User":         "用户",
	"Page":         "页面",
	"Login":        "登录",

	"POST":   "新增或修改",
	"PUT":    "修改",
	"DELETE": "删除",
}

// 代理模式
type DecResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer // 缓存
}

func (w DecResponseWriter) Write(data []byte) (int,error){
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w DecResponseWriter) WriteString(s string) (int,error){
	w.body.Write([]byte(s))
	return w.ResponseWriter.WriteString(s)
}

func GetOptionMap(key string)string{
	return optMap[key]
}

//// "gin-blog/api/v1.(*Resource).Delete-fm" => "Resource
func getOptResource (handleName string) string{
	s := strings.Split(handleName,".")[1]
	return s[2:len(s)-1]
}

func OperationLog() gin.HandlerFunc{
	return func (c *gin.Context){
		// 如果是get方法(太多)和文件上传操作(body太长),不记录
		if c.Request.Method != "GET" || !strings.Contains(c.Request.RequestURI,"upload") {// RequstURI:  /search?q=golang HTTP/1.1
			blw := DecResponseWriter{
				ResponseWriter: c.Writer,
				body: bytes.NewBufferString(""),
			}
			
			//自定义gin的Writer 和 request body
			c.Writer = blw
			body,_ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			
			auth,_ := handle.CurrentUserAuth(c)
			ipAddress := utils.IP.GetIpaddress(c)
			ipSource := utils.IP.GetIPsource(ipAddress)
			moduleName := getOptResource(c.HandlerName())

			operationLog := model.OperationLog{
				OptModule:     moduleName, // TODO: 优化
				OptType:       GetOptionMap(c.Request.Method),
				OptUrl:        c.Request.RequestURI,
				OptMethod:     c.HandlerName(),
				OptDesc:       GetOptionMap(c.Request.Method) + moduleName, // TODO: 优化
				RequestParam:  string(body),
				RequestMethod: c.Request.Method,
				UserId:        auth.UserInfoId,
				Nickname:      auth.UserInfo.Nickname,
				IpAddress:     ipAddress,
				IpSource:      ipSource,
			}
			c.Next()
			operationLog.ResponseData = blw.body.String()
			
			db := c.MustGet("db").(*gorm.DB)
			if result := db.Create(&operationLog).Error; result != nil {
				slog.Error("记录操作日志失败",result)
				handle.ReturnError(c, g.ErrDbOp, result)
				return
			}
		}else {
			c.Next()
			return
		}

	}
}
