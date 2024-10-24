package global

// gin | session
const (
	CTX_DB = "_db_filed"
	CTX_RDB =  "_rdb_filed"
	CTX_USER_AUTH = "_user_auth_filed"
)

// redis key
const (
	PAGE = "page" // 页面封面
	CONFIG = "config" // 博客配置
	
	VIEW_COUNT   = "view_count"  // 文章浏览数
	VISITOR_AREA = "visitor_area" // 参观的地域

	KEY_UNIQUE_VISITOR_SET = "unique_visitor" // 唯一用户记录 set

	OFFLINE_USER = "offline_user:" // 强制下线用户
	ONLINE_USER  = "online_user:"  // 在线用户

	ARTICLE_USER_LIKE_SET = "article_user_like:" // 文章点赞 Set
	ARTICLE_LIKE_COUNT = "article_like_count" // 文章点赞数
	ARTICLE_VIEW_COUNT = "article_view_count" // 文章浏览数

	COMMENT_USER_LIKE_SET = "comment_user_like:" // 评论点赞 Set
	COMMENT_LIKE_COUNT    = "comment_like_count" // 评论点赞数
	ARTICLE_COMMENT_COUNT = "article_comment_count" // 文章评论数
)

// config key
const (
	CONFIG_ARTICLE_COVER = "article_cover" // 文章封面
	CONFIG_ABOUT             = "about"
	CONFIG_IS_COMMENT_REVIEW = "is_comment_review" // 评论默认是否需要审核
)