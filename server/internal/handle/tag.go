package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type Tag struct{}

type AddOrEditTagReq struct{
	ID  	int     	`json:"id"`
	Name    string  	`json:"name" Binding:"required"`
}

// 获取Tag的列表
func (*Tag) GetList(c *gin.Context){
	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}	
	
	data,total,err := model.GetTagList(GetDB(c),query.Page,query.Size,query.Keyword)
	if err != nil{
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,PageList[model.TagVO]{
		Data: data,
		Total: total,
		Page: query.Page,
		Size: query.Size,
	})
}

// 增加或更新tag
func (*Tag) SaveOrUpdate(c *gin.Context){
	var req AddOrEditTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}
	data,err := model.SaveOrCreateTag(GetDB(c),req.ID,req.Name)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,data)
}

// 删除tag
func (*Tag) Delete(c *gin.Context){
	var  ids []int
	if err := c.ShouldBindJSON(&ids); err!= nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}
	db := GetDB(c)
	// 删除Tag之前，需要先检查Tag下是否有文章，有文章则不能删除
	count,err := model.Count(db,model.ArticleTag{},"tag_id IN ?",ids)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	// 有文章
	if count > 0 {
		ReturnError(c,g.ErrTagHasArt,err)
		return
	}
	result := db.Delete(model.Tag{},"id IN?",ids)
	if result.Error!= nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	ReturnSuccess(c,nil)
}

// 获取Tag的 id-name 列表
func (*Tag) GetOption(c *gin.Context){
	data,err := model.GetTagOption(GetDB(c))
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,data)
}