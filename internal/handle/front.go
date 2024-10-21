// 这个handler主要职责是加载首页信息
package handle

import (
	"strconv"
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

type FaddCommentReq struct {
	ReplyUserId   int  		`json:"reply_user_id" form:"reply_user_id"`
	TopicId   	  int   	`json:"topic_id"  form:"topic_id"`
	Content       string 	`json:"content" form:"content"`
	ParentId      int 		`json:"parent_id" form:"parent_id"`
	Type		  int       `json:"type" form:"type"  validate:"required,min=1,max=3" label:"评论类型"`
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
	CategoryId   int	`json:"category_id"`
	TagId        int	`json:"tag_id"`
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
	var req FAddMessageReq
	db := GetDB(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	auth,_ := CurrentUserAuth(c)
	ipaddress := utils.IP.GetIpaddress(c)
	source := utils.IP.GetIPsource(ipaddress)
	Isreviewed := model.GetConfigBool(db,g.CONFIG_IS_COMMENT_REVIEW)
	info := auth.UserInfo

	message,err := model.CreatenewMessage(db,info.Nickname,info.Nickname,req.Content,ipaddress,source,req.Speed,Isreviewed)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	ReturnSuccess(c,message)
}

// 新增评论
func (*Front) SaveComment(c *gin.Context){
	var req FaddCommentReq
	if err := c.ShouldBindJSON(&req); err!=nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	db := GetDB(c)
	auth,_ := CurrentUserAuth(c)

	Isreviewed := model.GetConfigBool(db,g.CONFIG_IS_COMMENT_REVIEW)
	info := auth.UserInfo

	var data *model.Comment
	var err error
	if req.ReplyUserId == 0{
		data,err = model.AddComment(db,info.ID,req.Type,req.TopicId,req.Content,Isreviewed)
		if err != nil{
			ReturnError(c,g.ErrDbOp,err)
			return
		}
	}else {
		data,err = model.AddReplyComment(db,info.ID,req.ReplyUserId,req.ParentId,req.Content,Isreviewed)
		if err != nil{
			ReturnError(c,g.ErrDbOp,err)
			return
		}
	}

	ReturnSuccess(c,data)
}

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

