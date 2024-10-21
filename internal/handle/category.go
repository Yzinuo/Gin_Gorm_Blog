package handle

import (
	g "gin-blog/internal/global"
	"gin-blog/internal/model"

	"github.com/gin-gonic/gin"
)

type AddorEiditReq struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Category struct{}

// 更具Page keyword关键词查询分类及其文章数量
func (*Category) GetList(c *gin.Context){
	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}
	
	data,total,err:=model.GetCategoryList(GetDB(c),query.Page,query.Size,query.Keyword)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	
	ReturnSuccess(c,PageList[model.CategoryVO]{
		Page : query.Page,
		Size : query.Size,
		Data : data,
		Total : total,
	})
}

// 保存更新
func (*Category) SaveOrUPdate(c *gin.Context){
	var req AddorEiditReq
	if err := c.ShouldBindJSON(&req); err !=nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	category,err  := model.SaveOrUpdateCategory(GetDB(c),req.Id,req.Name)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,category)
}

// 删除
func (*Category) Delete(c *gin.Context){
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	count ,err := model.Count(GetDB(c),&model.Article{},"category_id IN ? ", ids)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	if count > 0 {
		ReturnError(c,g.ErrCateHasArt,err)
		return
	}

	rows,err := model.DeleteCategory(GetDB(c),ids)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,rows)
}

// 获取文章列表
func (*Category) GetCategoryOption(c *gin.Context){
	list,err := model.GetcategoryOption(GetDB(c))
	if err!= nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,list)
}
