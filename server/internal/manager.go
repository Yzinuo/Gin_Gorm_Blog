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
	base.GET("/email/verify",userAuthAPI.VerifyCode)
	base.GET("/logout",userAuthAPI.Logout)
	base.POST("/report",blogInfoAPI.Report)
	base.GET("/config",blogInfoAPI.GetConfigMap)
	base.PATCH("/config",blogInfoAPI.UpdateBlogInfo)
}

// 注册admin handler
func RegisterAdminHandler(r *gin.Engine){
	auth := r.Group("/api")

	// 对于这组中间键,使用JWT中间件鉴权,注意使用顺序
	auth.Use(middleware.JWTAuth())  // JWT验证
	auth.Use(middleware.PermissionCheck())// 验证是否有权限
	auth.Use(middleware.OperationLog()) // 记录操作
	auth.Use(middleware.ListenOnline()) // 刷新用户登录信息,如果是黑名单禁止访问
	
	auth.GET("/home",blogInfoAPI.GetBlogInfo)// 后台首页信息获取
	auth.POST("/upload",uploadAPI.UploadFile) //文件上传

	// 设置模块
	setting := auth.Group("/setting")
	{
		setting.GET("/about",blogInfoAPI.GetAbout)
		setting.PUT("/about",blogInfoAPI.UpdateAbout)		
	}

	// 用户模块
	user := auth.Group("/user")
	{
		user.GET("/list",userAPI.GetList)
		user.PUT("",userAPI.Update)  // 更新用户信息
		user.PUT("/disable",userAPI.UpdateDisable) // 修改用户禁用状态
		user.PUT("/current/password",userAPI.UpdateCurrentPasswordReq) // 修改当前用户密码
		user.GET("/info",userAPI.GetInfo) // 获取当前用户信息
		user.GET("/online",userAPI.GetOnlineList) // 获取在线用户
		user.POST("/offline/:id",userAPI.ForceOffline) // 强制用户下线
	}

	category := auth.Group("category")
	{
		category.GET("/list",categoryAPI.GetList) //更具关键词挑选出分类以及它文章的数量
		category.POST("",categoryAPI.SaveOrUPdate)
		category.DELETE("",categoryAPI.Delete)
		category.GET("/option",categoryAPI.GetCategoryOption) // 得到分类的列表
	}
	article := auth.Group("article")
	{
		article.GET("list",articleAPI.GetList) //查询文章
		article.POST("",articleAPI.SavaOrUpdate)
		article.PUT("top",articleAPI.UPdateTOP)
		article.GET("/:id",articleAPI.GetDetail) //文章的详细
		article.PUT("/soft-delete",articleAPI.SoftDelArticle)
		article.DELETE("",articleAPI.DeleteArticle)
		article.POST("/export",articleAPI.Export)
		article.POST("/import",articleAPI.Import) // 导入文章
	}
	comment := auth.Group("comment")
	{
		comment.GET("/list",commentAPI.GetList)
		comment.PUT("/review",commentAPI.UpdateReview) // 更新回复
		comment.DELETE("",commentAPI.Delete)
	}
	tag := auth.Group("tag")
	{
		tag.GET("/list",tagAPI.GetList)
		tag.POST("",tagAPI.SaveOrUpdate)
		tag.DELETE("",tagAPI.Delete)
		tag.GET("/option",tagAPI.GetOption)
	}
	message := auth.Group("message")
	{
		message.GET("/list",messageAPI.GetList)
		message.DELETE("",messageAPI.Delte)
		message.PUT("/review",messageAPI.UpdateReview) // 审核留言
	}
	link := auth.Group("link")
	{
		link.GET("/list",linkAPI.GetList)
		link.DELETE("",linkAPI.Delete)
		link.POST("",linkAPI.SaveOrUpdateLink)
	}
	resource := auth.Group("resource")
	{
		resource.GET("/list",resourceAPI.GetTreeList)  // 树形的资源列表
		resource.POST("",resourceAPI.SaveOrUpdate)
		resource.DELETE("/:id",resourceAPI.Delete)
		resource.PUT("/anonymous",resourceAPI.UpdateAnonymous) //修改匿名访问
		resource.GET("/option",resourceAPI.GetOption)
	}
	menu := auth.Group("menu")
	{
		menu.GET("list",menuAPI.GetTreeList)
		menu.POST("",menuAPI.SaveOrUPdateMenu)
		menu.DELETE("/:id",menuAPI.Delete)
		menu.GET("/user/list",menuAPI.GetUserMenu) // 用户自己的菜单
		menu.GET("/option",menuAPI.GetOption)
	}
	role := auth.Group("role")
	{
		role.GET("/list",roleAPI.GetList)
		role.POST("",roleAPI.SaveOrUpdate)
		role.DELETE("",roleAPI.Delete)
		role.GET("/option",roleAPI.GetOption)
	}
	operationLog := auth.Group("operation")
	{
		operationLog.GET("/log/list",operationLogAPI.GetList)
		operationLog.DELETE("",operationLogAPI.Delete)
	}
	page := auth.Group("/page")
	{
		page.GET("/list",pageAPI.GetList)
		page.POST("",pageAPI.SaveAndUpdate)
		page.DELETE("",pageAPI.Delete)
	}
}

//注册博客展示 front handler
func RegisterFrontHandler(r *gin.Engine){
	// 注册front handler
	base := r.Group("api/front")

	base.GET("/home",frontAPI.GetHomeInfo) //获得博客首页信息
	base.GET("/page",pageAPI.GetList) // 前台页面
	base.GET("/about",blogInfoAPI.GetAbout)

	article := base.Group("/article")
	{
		article.GET("/list",frontAPI.GetArticleList) // 前台文章列表
		article.GET("/:id",frontAPI.GetArticleInfo) // 前台文章详情
		article.GET("/archive",frontAPI.GetArchiveList)//归档
		article.GET("/search",frontAPI.SearchArticle)
	}
	category := base.Group("/category")
	{
		category.GET("/list",frontAPI.GetCategoryList)  // 得到分类列表
	}
	tag := base.Group("/tag")
	{
		tag.GET("/list",frontAPI.GetTagList)
	}
	link := base.Group("/link")
	{
		link.GET("/list",frontAPI.GetLinkList)
	}
	message := base.Group("/message")
	{
		message.GET("/list",frontAPI.GetMessageList)
	}
	comment := base.Group("/comment")
	{
		comment.GET("/list",frontAPI.GetCommentList)
		comment.GET("/replies/:comment_id",frontAPI.GetReplyListByCommentId)
	}

	base.Use(middleware.JWTAuth()) //需要登录才能进行的操作，在前端中使用下列API还会有额外的逻辑：检查是否有用户信息，没有则跳出登录框
	{
		base.POST("/upload",uploadAPI.UploadFile) // 上传文件
		base.GET("/user/info",userAPI.GetInfo)  // 获得用户信息
		base.PUT("/user/info",userAPI.UpdateCurrent) // 更新用户信息
		base.POST("/message",frontAPI.SaveMessage)   //新增留言
		base.POST("/comment",frontAPI.SaveComment)   // 新增评论
		//本身是要用POST方法的，但是因为handler内部接收数据的细节，使用GET
		base.GET("/comment/like/:comment_id",frontAPI.LikeComment) // 前台点赞评论
		base.GET("/article/like/:article_id",frontAPI.LikeArticle) // 前台点赞文章
	}

}