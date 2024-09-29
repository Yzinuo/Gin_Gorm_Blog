package handle

import (
	g "gin-blog/internal/global"
	model "gin-blog/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

type Article struct {}

type AddOrEditeArticleReq struct {
	ID		int		`json:"id"`
	Title   string	`json:"title"  binding:"required"` 
	Desc 	string  `json:"desc"`
	Content string  `json:"content" binding:"required"`
	Img	string  `json:"imag"`
	Status  int     `json:"status"  binding:"required,min=1,max = 3"`
	Type 	int 	`json:"type"    binding:"required, min=1,max=3"`	
	IsTop   bool 	`json:"is_top"`
	OriginalUrl string `json:"original_url"`
	
	Tagnames  []string `json:"tagnames"`	
	CategoryName string `json:"category_name"`
}

// 注意是软删除
type SoftDelArticleReq struct {
	IDs 	[]int 	`json:"ids" binding:"required"`	
	IsDel   bool 	`json:"is_del"`
}

type QueryArticle struct {
	PageQuery
	Title 	   string `form:"title"`	
	Type  	   int    `form:"type"`
	TagId 	   int    `form:"tag_id"`
	CategoryId int    `form:"category_id"`
	Status 	   int    `form:"status"`
	IsDelete   *bool   `form:"is_delete"`
}

type UpdateArticleTopReq struct {
	ID		int		`json:"id"`
	IsTop   bool 	`json:"is_top"`
}

type ArticleVO struct{
	model.Article
	
	LikeCount int  `gorm:"-" json:"like_count"`
	ViewCount int  `gorm:"-" json:"view_count"`
	CommentCount int  `gorm:"-" json:"comment_count"`
}

// 增或改文章
func (*Article)SavaOrUpdateArticle (c *gin.Context) {
	req := AddOrEditeArticleReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}
	userauth,_ := CurrentUserAuth(c)
	db := GetDB(c)

	if req.Img == "" {
		req.Img,_  = model.GetValueByKey(db,g.CONFIG_ARTICLE_COVER)
	}

	if req.Type == 0 {
		req.Type = 1  // 默认原创
	}
	
	article := model.Article{
		Model: model.Model{ID:req.ID},
		Title: req.Title,
		Desc: req.Desc,
		Content: req.Content,
		Img: req.Img,
		Status: req.Status,
		Type: req.Type,
		IsTop: req.IsTop,
		OriginalUrl: req.OriginalUrl,
	
		UserId: userauth.ID,
	}

	err := model.SvaeOrUpdateArticle(db,&article,req.CategoryName,req.Tagnames)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	
	ReturnSuccess(c,article)
}

// 软删除文章
func (*Article) SoftDelArticle(c *gin.Context) {
	delreq := SoftDelArticleReq{}

	if err := c.ShouldBind(&delreq); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	rows,err := model.UpdateArticleSoftDlete(GetDB(c),delreq.IDs,delreq.IsDel)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,rows)
}

// 删除文章
func (*Article) DeleteArticle(c *gin.Context) {
	var ids []int
	
	if err := c.ShouldBind(&ids); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	rows,err := model.DeleteArticle(GetDB(c),ids)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,rows)
}

// 查询文章
func (*Article)QureyAticle(c *gin.Context,rdb *redis.Client) {
	Queryreq := QueryArticle{}
	if err := c.ShouldBind(&Queryreq); err != nil{
		ReturnError(c,g.ErrRequest,err)
		return
	}

	articles,rows,err := model.GetArticleList(GetDB(c),Queryreq.Page,Queryreq.Size,Queryreq.Title,Queryreq.IsDelete,
													Queryreq.Status,Queryreq.Type,Queryreq.CategoryId,Queryreq.TagId)
	if err!= nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	
	likeCountMap := rdb.HGetAll(c,g.ARTICLE_LIKE_COUNT).Val()
	viewCount := rdb.ZRangeWithScores(c,g.ARTICLE_VIEW_COUNT,0,-1).Val()

	var viewCountMap = make(map[int]int)
	Data := make([]ArticleVO,0)
	
	for _,article := range viewCount{
		key ,_:= strconv.Atoi(article.Member.(string))
		viewCountMap[key] = int(article.Score)
	}

	for _,article := range articles {
		LikeCount,_ :=   strconv.Atoi(likeCountMap[strconv.Itoa(article.ID)])
		Data = append(Data, 
			ArticleVO{
				Article: article,
				LikeCount: LikeCount,
				ViewCount: viewCountMap[article.ID],
			},
		)	
	}

	ReturnSuccess(c,PageList[ArticleVO]{
		Page: Queryreq.Page,
		Size: Queryreq.Size,
		Total: rows,
		Data: Data,
	})
}


