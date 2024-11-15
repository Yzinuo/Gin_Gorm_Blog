// 这个handler主要职责是加载首页信息
package handle

import (
	"strconv"
	"strings"
	"time"

	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"gin-blog/internal/utils"

	"github.com/gin-gonic/gin"
)

type Front struct{}

// 留言请求
type FAddMessageReq struct {
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar"`
	Content  string `json:"content"  binding:"required"`
	Speed    int 	`json:"speed"`
}

type FAddCommentReq struct {
	ReplyUserId int    `json:"reply_user_id" form:"reply_user_id"`
	TopicId     int    `json:"topic_id" form:"topic_id"`
	Content     string `json:"content" form:"content"`
	ParentId    int    `json:"parent_id" form:"parent_id"`
	Type        int    `json:"type" form:"type" validate:"required,min=1,max=3" label:"评论类型"`
}

type FCommentQuery struct {
	PageQuery
	ReplyUserId int 	`json:"reply_user_id" form:"reply_user_id"`
	TopicId     int 	`json:"topic_id" form:"topic_id"`
	Content     string  `json:"content" form:"content"`
	ParentId    int 	`json:"parent_id" form:"parent_id"`
	Type        int    	`json:"type" form:"type"`
}

type FArticleQuery struct {
	PageQuery
	CategoryId   int	`form:"category_id"`
	TagId        int	`form:"tag_id"`
}

type ArchiveVO struct{
	ID			int 	`json:"id"`
	Title 		string  `json:"title"`
	Created_at  time.Time `json:"created_at"`
}

type ArticleSearchVO struct{
	ID 			int 	`json:"id"`
	Title 		string 	`json:"title"`
	Content 	string 	`json:"content"`
}

// 获取网站的数据
func (*Front) GetHomeInfo (c *gin.Context){
	db := GetDB(c)
	rdb := GetRDB(c)

	data,err := model.GetFrontStatistics(db)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	data.ViewCount,_ = rdb.Get(rdbctx,g.VIEW_COUNT).Int64()
	
	ReturnSuccess(c,data)
}

func (*Front) GetTagList(c *gin.Context) {
	data,_,err:= model.GetTagList(GetDB(c),1,1000,"")
	if err != nil{
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	ReturnSuccess(c,data)
}

func (*Front) GetMessageList(c *gin.Context) {
	data,_,err:= model.GetMessageList(GetDB(c),1,1000,"")
	if err != nil{
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	ReturnSuccess(c,data)
}


func (*Front) GetCategoryList(c *gin.Context) {
	data,_,err:= model.GetCategoryList(GetDB(c),1,1000,"")
	if err != nil{
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	ReturnSuccess(c,data)
}

func (*Front) GetLinkList(c *gin.Context) {
	data,_,err:= model.GetLinkList(GetDB(c),1,1000,"")
	if err != nil{
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	ReturnSuccess(c,data)
}

//新增留言
func (*Front) SaveMessage(c *gin.Context){
	auth,err := CurrentUserAuth(c)
	if err != nil {
		ReturnError(c,g.ErrTokenRuntime,err)
		return
	}
	
	var req FAddMessageReq
	db := GetDB(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}


	ipaddress := utils.IP.GetIpaddress(c)
	source := utils.IP.GetIPsource(ipaddress)
	Isreviewed := model.GetConfigBool(db,g.CONFIG_IS_COMMENT_REVIEW)
	info := auth.UserInfo

	message,err := model.CreatenewMessage(db,info.Nickname,info.Avatar,req.Content,ipaddress,source,req.Speed,Isreviewed)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	ReturnSuccess(c,message)
}

// 新增评论
func (*Front) SaveComment(c *gin.Context){
	var req FAddCommentReq
	if err := c.ShouldBindJSON(&req); err!=nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	db := GetDB(c)
	auth,_ := CurrentUserAuth(c)

	Isreviewed := model.GetConfigBool(db,g.CONFIG_IS_COMMENT_REVIEW)

	var data *model.Comment
	var err error
	if req.ReplyUserId == 0{
		data,err = model.AddComment(db,auth.ID,req.Type,req.TopicId,req.Content,Isreviewed)
		
	}else {
		data,err = model.AddReplyComment(db,auth.ID,req.ReplyUserId,req.ParentId,req.Content,Isreviewed)
	}
	
	if err != nil{
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,data)
}

// 获取评论列表
func (*Front) GetCommentList(c *gin.Context){
	var query FCommentQuery
	if err := c.ShouldBindQuery(&query);err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	data,total,err := model.GetBlogCommentList(db,query.Page,query.Size,query.Type,query.TopicId)
	if err != nil{
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	LikeMap := rdb.HGetAll(rdbctx,g.COMMENT_USER_LIKE_SET).Val()
	for i,comment := range data {
		if len(comment.ReplyList) >= 3{
			data[i].ReplyList = data[i].ReplyList[:3]
		}
		data[i].LikeCount ,_= strconv.Atoi(LikeMap[strconv.Itoa(comment.ID)])
	}

	ReturnSuccess(c,PageList[model.CommentVO]{
		Page : query.Page,
		Size : query.Size,
		Data : data,
		Total : total,
	})
}

// 获取评论的回复列表
func (*Front) GetReplyListByCommentId(c *gin.Context){
	commentid,err := strconv.Atoi(c.Param("comment_id"))
	if err!= nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}
	db := GetDB(c)
	rdb := GetRDB(c)

	var query PageQuery
	if err = c.ShouldBindQuery(&query);err != nil{
		ReturnError(c,g.ErrRequest,err)
		return
	}

	replyList,err := model.GetCommentByid(db,query.Page,query.Size,commentid)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	data := make([]model.CommentVO,0)
	likeCountMap := rdb.HGetAll(rdbctx,g.COMMENT_LIKE_COUNT).Val()
	for _,reply := range replyList{
		likeCount,_ := strconv.Atoi(likeCountMap[strconv.Itoa(reply.ID)])
		data = append(data, model.CommentVO{
			Comment : reply,
			LikeCount: likeCount,
		})
	}
	
	ReturnSuccess(c,data)
}

// 点赞评论
func (*Front) LikeComment(c *gin.Context){
	commentid,err  := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	auth,_ := CurrentUserAuth(c)
	rdb := GetRDB(c)
	
	// 判断用户是否点赞过这条评论
	key := g.COMMENT_USER_LIKE_SET + strconv.Itoa(auth.ID)
	if rdb.SIsMember(rdbctx,key,commentid).Val(){
		// 点赞过，取消点赞
		rdb.SRem(rdbctx,key,commentid)
		rdb.HIncrBy(rdbctx,g.COMMENT_LIKE_COUNT,strconv.Itoa(commentid),-1)
	}else{
		// 没有点赞过，点赞
		rdb.SAdd(rdbctx,key,commentid)
		rdb.HIncrBy(rdbctx,g.COMMENT_LIKE_COUNT,strconv.Itoa(commentid),1)
	}

	ReturnSuccess(c,nil)
}


// 关于文章的handler

// 获取博客前端文章列表
func (*Front) GetArticleList(c *gin.Context){
	var query  FArticleQuery 
	if err := c.ShouldBindQuery(&query);err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}
	
	data ,err,_:= model.GetBlogArticleList(GetDB(c),query.Page,query.Size,query.CategoryId,query.TagId)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	ReturnSuccess(c,data)
}

// 得到文章的信息
func (*Front) GetArticleInfo(c *gin.Context){
	id,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)
	data,err :=	model.GetBlogArticle(db,id)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	// 获得六篇推荐文章 相同tag
	article := model.BlogArticleVO{Article: *data}
	article.RecommendArticles ,err = model.GetRecommandList(db,id,6)
	if err != nil{
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	//获取最新文章
	article.NewestArticles ,err = model.GetNewestList(db,5)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	// 获取上一篇和下一篇文章 文章的评论量
	article.NextArticle,err = model.GetNextArticle(db,id)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	article.LastArticle,err = model.GetLastArticle(db,id)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	article.CommentCount,err = model.GetCommentCountOfArticle(db,id)
	if err!= nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	
	// 增加文章的访问量
	rdb.ZIncrBy(rdbctx,g.ARTICLE_VIEW_COUNT,1,strconv.Itoa(id))
	
	// 获取文章的点赞量，访问量
	article.LikeCount = int64(rdb.ZScore(rdbctx,g.ARTICLE_LIKE_COUNT,strconv.Itoa(id)).Val())

	likeCount, _ := strconv.Atoi(rdb.HGet(rdbctx, g.ARTICLE_LIKE_COUNT, strconv.Itoa(id)).Val())
	article.LikeCount = int64(likeCount)

	ReturnSuccess(c,article)
}

// 文章归档  把所有文章都以id title 创建时间 形式存储。 
// 方便后续查找 有序
func (*Front) GetArchiveList(c *gin.Context){
	var query FArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	data,err ,total:= model.GetBlogArticleList(GetDB(c),query.Page,query.Size,query.CategoryId,query.TagId)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	archive := make([]ArchiveVO,0)
	for _,article := range data{
		archive = append(archive, ArchiveVO{
			ID: article.ID,
			Title: article.Title,
			Created_at: article.CreatedAt,
		})
	}

	ReturnSuccess(c,PageList[ArchiveVO]{
		Total: total,
		Data: archive,
		Page: query.Page,
		Size: query.Size,
	})
}

func (*Front) LikeArticle(c *gin.Context){
	id,err := strconv.Atoi(c.Param("article_id"))
	if err != nil{
		ReturnError(c,g.ErrRequest,err)
		return
	}

	auth,_ := CurrentUserAuth(c)
	rdb := GetRDB(c)
	
	key := g.ARTICLE_USER_LIKE_SET+strconv.Itoa(auth.ID)
	if rdb.SIsMember(rdbctx,key,id).Val(){
		// 取消点赞
		rdb.SRem(rdbctx,key,id)
		rdb.HIncrBy(rdbctx,g.ARTICLE_LIKE_COUNT,strconv.Itoa(id),-1)
	}else{
		// 点赞
		rdb.SAdd(rdbctx,key,id)
		rdb.HIncrBy(rdbctx,g.ARTICLE_LIKE_COUNT,strconv.Itoa(id),1)
	}

	ReturnSuccess(c,nil)
}

func (*Front) SearchArticle(c *gin.Context){
	result := make([]ArticleSearchVO,0)
	keyword := c.Query("keyword")
	if keyword == ""{
		ReturnError(c,g.ErrRequest,nil)
		return
	}	
	db := GetDB(c)

	articleList,err := model.List(db,[]model.Article{},"*","",
						"is_delete = 0 AND status = 1 AND (title LIKE ? OR content LIKE ?)","%"+keyword+"%","%"+keyword+"%")
	if err!= nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	// 高亮关键词 ，设置搜索结果
	for _,article := range articleList{
		// 如果标题存在 高亮标题
		strings.ReplaceAll(article.Title,keyword,"<span style='color:#f47466'>"+keyword+"</span>")
		// 高亮内容
		content := article.Content
		keywordIndex := unicodeIndex(content,keyword)
		
		// 关键词在内容中
		if keywordIndex != -1{
			preIndex,afterIndex := 0,0
			if keywordIndex > 25 {
				preIndex = keywordIndex - 25
			}
		// 防止中文出乱码
			preText := substring(content,preIndex,keywordIndex)

			keywordEndIndex := keywordIndex + unicodeLen(keyword)
			afterLength := len(content) - keywordEndIndex
			if afterLength > 175 {
				afterIndex = keywordEndIndex + 175
			} else {
				afterIndex = keywordEndIndex + afterLength
			}
			// afterText := string([]rune(content)[keywordStartIndex:afterIndex])
			afterText := substring(content, keywordIndex, afterIndex)
			// 高亮内容中的关键字
			content = strings.ReplaceAll(preText+afterText, keyword,
				"<span style='color:#f47466'>"+keyword+"</span>")
		}
		
		result = append(result,ArticleSearchVO{
			ID: article.ID,
			Title: article.Title,
			Content: content,
		})

	}
	ReturnSuccess(c,result)
}

// UTF8是变长编码， 英文1字节，中文3字节。对切片的操作都是按字节操作的，不是字符 所以内容中有中文很容易操作出错
// 把切片转换为rune切片 就能保证操作是按字符操作的，所以不会乱码
// 获取中文的字符串中字符串的实际位置，而非字节位置
func unicodeIndex (str,substr string) int{
	result := strings.Index(str,substr)
	// 存在这个字符
	if result > 0 {
		preText := []byte(str)[0:result]
		context := []rune(string(preText))
		result = len(context)
	}
	return result
}

// 获取代中文的字符串实际长度，而非字节长度
func unicodeLen(str string) int{
	context := []rune(str)
	return len(context)
}

// 解决中文获取位置不对的问题
func substring(source string, start int, end int) string{
	str := []rune(source)
	len := len(str)
	
	if start >=  end{
		return ""
	}

	if start < 0{
		start = 0
	}
	if end > len{
		end = len
	}
	if start == 0 && end == len{
		return source
	}

	var substr string
	for i:= start; i < end ; i++{
		substr = substr + string(str[i])
	}

	return substr
}