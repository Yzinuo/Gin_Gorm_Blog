package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type Link struct{}

// 增加或修改友链
type AddOrEditLinkReq struct{
	ID		int 		`json:"id"`
	Name	string		`json:"name"`
	Avatar  string 		`json:"avatar"`
	Address string		`json:"address"`
	Intro   string      `json:"intro"`
}

// 查询友链
func  (*Link) GetList(c *gin.Context){
	var query PageQuery
	if err := c.ShouldBindJSON(&query); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	data,total,err := model.GetLinkList(GetDB(c),query.Page,query.Size,query.Keyword)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,PageList[model.FriendLink]{
		Data: data,
		Total: total,
		Page: query.Page,
		Size: query.Size,
	})
}

// 增加或修改友链
func (*Link) SaveOrUpdateLink(c *gin.Context){
	var req AddOrEditLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	data,err := model.SvaeOrCreateLink(GetDB(c),req.ID,req.Name,req.Avatar,req.Address,req.Intro)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,data)
}

// 删除友链
func (*Link) Delete(c *gin.Context){
	var ids []int
	if err := c.ShouldBindJSON(&ids);err != nil{
		ReturnError(c,g.ErrRequest,err)
		return
	}

	err := GetDB(c).Delete(&model.FriendLink{},"id IN ?",ids).Error
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,nil)
}