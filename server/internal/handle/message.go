package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type UpdateRevirewReq struct {
	Ids		 []int	`json:"ids"`
	Isreview bool   `json:"is_review"`
}

type MessageQuery struct{
	PageQuery
	Nickname string `json:"nickname"`
	IsReview  *bool  `json:"is_review"`
}

type Message struct{}


func (*Message) Delte(c *gin.Context){
	var ids []int
	if err := c.ShouldBindJSON(&ids);err!=nil{
		ReturnError(c,g.ErrRequest,err)
		return
	}

	rows,err := model.DeleteMessage(GetDB(c),ids)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	ReturnSuccess(c,rows)
}

func (*Message) GetList(c *gin.Context){
	var query MessageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	data,total,err := model.GetMessageList(GetDB(c),query.Page,query.Size,query.Nickname)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,PageList[model.Message]{
		Total: total,
		Data: data,
		Page: query.Page,
		Size: query.Size,
	})
}


func (*Message) UpdateReview(c *gin.Context) {
	var req UpdateRevirewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	rows, err := model.UpdateMessageReview(GetDB(c), req.Ids, req.Isreview)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}