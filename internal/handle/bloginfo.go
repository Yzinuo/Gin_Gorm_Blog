// 负责获取博客的首页信息
package handle

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
