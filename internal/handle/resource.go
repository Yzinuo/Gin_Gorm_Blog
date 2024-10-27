// resource 存储着每个方法（handler）对应的url 和使用的http方法。
// 具有树形结构
package handle

import (
	"strconv"
	"time"

	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Resource struct{}

type ResourceTreeVO struct {
	ID        int              `json:"id"`
	CreatedAt time.Time        `json:"created_at"`
	Name      string           `json:"name"`
	Url       string           `json:"url"`
	Method    string           `json:"request_method"`
	Anonymous bool             `json:"is_anonymous"`
	Children  []ResourceTreeVO `json:"children"`
}

// TODO: 使用 oneof 标签校验数据
type AddOrEditResourceReq struct {
	ID       int    `json:"id"`
	Url      string `json:"url"`
	Method   string `json:"request_method"`
	Name     string `json:"name"`
	ParentId int    `json:"parent_id"`
}

type EditAnonymousReq struct {
	ID        int  `json:"id" binding:"required"`
	Anonymous bool `json:"is_anonymous"`
}

// 更新或新建资源
func (*Resource) SaveOrUpdate(c *gin.Context) {
	var req AddOrEditResourceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	db := GetDB(c)
	err := model.SaveOrUpdateResource(db, req.ID, req.ParentId, req.Name, req.Url, req.Method)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// 编辑资源的匿名访问
func (*Resource) UpdateAnonymous(c *gin.Context) {
	var req EditAnonymousReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	err := model.UpdateResourceAnonymous(GetDB(c), req.ID, req.Anonymous)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// 删除资源， 如果资源有子资源，删除后，子资源无法处理，因此设定有子资源无法删
func (*Resource) Delete(c *gin.Context){
	resourceId,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	db := GetDB(c)
	// 检查是否有人能使用resource 有的话不能删除
	result,err := model.CheckResourceInUse(db, resourceId)
	if result{
		ReturnError(c,g.ErrResourceUsedByRole,err)
		return
	}

	// 查询是否有子资源
	resource,err := model.GetResourceById(db, resourceId)
	if err != nil && errors.Is(err,gorm.ErrRecordNotFound){
		ReturnError(c,g.ErrResourceNotExist,err)
		return
	}else if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	if resource.ParentId == 0{
		result,err = model.CheckResourceHasChild(db,resourceId)
		if result{
			ReturnError(c,g.ErrResourceHasChildren,err)
			return
		}
	}
	// 删除
	if db.Delete(&model.Resource{},resourceId).Error != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,nil)
}