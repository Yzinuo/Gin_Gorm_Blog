package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type UPdateReviewReq struct {
	Ids 		[]int `json:"ids"`
	IsReview	bool  `json:"is_review"`
}

type CommentQuery struct{
	PageQuery 
	Nickname  string  `form:"nickname"`
	IsReview  *bool	  `form:"is_review"`	
	Type	  int 	  `form:"type"`
}

type Comment struct{}

func (*Comment) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	result := GetDB(c).Delete(&model.Comment{},"id IN ?",ids)
	if result.Error != nil {
		ReturnError(c,g.ErrDbOp,result.Error)
		return
	}
	ReturnSuccess(c,result.RowsAffected)
}

func(*Comment) UpdateReview(c *gin.Context){
	var req UPdateReviewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	maps  := map[string]any{"is_review": req.IsReview}
	
	result := GetDB(c).Model(&model.Comment{}).Where("id IN ?", req.Ids).Updates(maps)
	if result.Error != nil {
		ReturnError(c,g.ErrDbOp,result.Error)
		return
	}

	ReturnSuccess(c,nil)
}

func (*Comment) GetList(c *gin.Context){
	var req CommentQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	list,total,err :=model.GetCommentList(GetDB(c),req.IsReview,req.Page,req.Size,req.Type,req.Nickname)
	if err != nil{
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,PageList[model.Comment]{
		Page: req.Page,
		Size: req.Size,
		Total: total,
		Data: list,
	})
}

