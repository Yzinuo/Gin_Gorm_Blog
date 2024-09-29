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


	ARTICLE_LIKE_COUNT = "article_like_count" // 文章点赞数
	ARTICLE_VIEW_COUNT = "article_view_count" // 文章浏览数
	ARTICLE_COMMENT_COUNT = "article_comment_count" // 文章评论数
)

// config key
const (
	CONFIG_ARTICLE_COVER = "article_cover" // 文章封面
)