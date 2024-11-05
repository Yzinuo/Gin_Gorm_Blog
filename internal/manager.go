// 注册handler，把handl中的handle函数和对应的route path对应。
package ginblog

import (
	"gin-blog/internal/handle"
	"gin-blog/internal/middleware"

	"github.com/gin-gonic/gin"
)

var (
	// 后台管理系统接口

	categoryAPI     handle.Category     // 分类
	tagAPI          handle.Tag          // 标签
	articleAPI      handle.Article      // 文章
	userAPI         handle.User         // 用户
	userAuthAPI     handle.Auth     // 用户账号
	commentAPI      handle.Comment      // 评论
	uploadAPI       handle.Upload       // 文件上传
	messageAPI      handle.Message      // 留言
	linkAPI         handle.Link         // 友情链接
	roleAPI         handle.Role         // 角色
	resourceAPI     handle.Resource     // 资源
	menuAPI         handle.Menu         // 菜单
	blogInfoAPI     handle.BlogInfo     // 博客设置
	operationLogAPI handle.OperationLog // 操作日志
	pageAPI         handle.Page         // 页面

	// 博客前台接口

	frontAPI handle.Front // 博客前台接口汇总
)

// 使用外观设计模式一口气全部完成注册
func RegisterAllHandler(r *gin.Engine){
	RegisterBaseHandler(r)
	RegisterAdminHandler(r)
	RegisterFrontHandler(r)
}

// 注册admin front通用的handler
func RegisterBaseHandler(r *gin.Engine){
	base := r.Group("/api")
	
	base.POST("/login",userAuthAPI.Login)
	base.POST("/register",userAuthAPI.Register)
	base.GET("/logout",userAuthAPI.Logout)
	base.POST("/report",blogInfoAPI.Report)
	base.GET("/config",blogInfoAPI.GetConfigMap)
	base.PATCH("/config",blogInfoAPI.UpdateBlogInfo)
}

// 注册admin handler
func RegisterAdminHandler(r *gin.Engine){
	auth := r.Group("api")

	// 对于这组中间键,使用JWT中间件鉴权,注意使用顺序
	auth.Use(middleware.JWTAuth())  // JWT验证
	auth.Use(middleware.PermissionCheck())// 验证是否有权限
	auth.Use(middleware.OperationLog()) // 记录操作
	auth.Use(middleware.ListenOnline()) // 刷新用户登录信息,如果是黑名单禁止访问
	
}

//注册博客展示 front handler
func RegisterFrontHandler(r *gin.Engine){
	// 注册front handler
	r.GET("/front", FrontHandler)
	r.GET("/front/ping", FrontHandler)
}