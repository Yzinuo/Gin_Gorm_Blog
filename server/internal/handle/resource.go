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
	rows,err :=model.DeleteResource(db,resourceId)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,rows)
}

// 获取资源树型结构
func (*Resource) GetTreeList(c *gin.Context){
	keyword := c.Query("keyword")
	
	resourcelist,err := model.GetResourceList(GetDB(c),keyword)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	ReturnSuccess(c,resource2ResourceVos(resourcelist))
}

// 获取资源列表(树形)
func (*Resource) GetOption(c *gin.Context) {
	result := make([]TreeOptionVO, 0)

	db := GetDB(c)
	resources, err := model.GetResourceList(db, "")
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	parentList := getModuleList(resources)
	childrenMap := getResourceChildrenMap(resources)

	for _, item := range parentList {
		var children []TreeOptionVO
		for _, re := range childrenMap[item.ID] {
			children = append(children, TreeOptionVO{
				ID:    re.ID,
				Label: re.Name,
			})
		}
		result = append(result, TreeOptionVO{
			ID:       item.ID,
			Label:    item.Name,
			Children: children,
		})
	}
	ReturnSuccess(c, result)
}

// 把resource列表转换成ReourceTree型结构
func resource2ResourceVos(resources []model.Resource) []ResourceTreeVO{
	//首先遍历第一级的资源，为每一个一级资源创造一个子树
	// 通过children map，丰富这个子树的children
	result := make([]ResourceTreeVO,0)
	first_level := getModuleList(resources)
	childrenMap := getResourceChildrenMap(resources)

	for _,resouce := range first_level{
		resource_tree := resource2ResourceTreeVo(resouce)
		for _,child := range childrenMap[resouce.ID]{
			resource_tree.Children = append(resource_tree.Children, resource2ResourceTreeVo(child))
		}
		result = append(result, resource_tree)
	}

	return result
}

// 把单个resource转换成Tree型结构
// 原型设计模式？
func resource2ResourceTreeVo(resource model.Resource) ResourceTreeVO{
	return ResourceTreeVO{ID: resource.ID, CreatedAt: resource.CreatedAt, Name: resource.Name, 
							Url: resource.Url, Method: resource.Method, Anonymous: resource.Anonymous}
}

// 获取一级资源  parent_id == 0
func getModuleList(resourcelist []model.Resource) []model.Resource{
	first_level := make([]model.Resource,0)
	for _,resource := range resourcelist{
		if resource.ParentId == 0{
			first_level = append(first_level,resource)
		}
	}

	return first_level
}

// 获取每个模块下的资源map
func getResourceChildrenMap(resourcelist []model.Resource) map[int][]model.Resource{
	childrenMap := make(map[int][]model.Resource,0)
	
	for _,resource := range resourcelist{
		if resource.ParentId!= 0{
			childrenMap[resource.ParentId] = append(childrenMap[resource.ParentId],resource)
		}
	}
	return childrenMap
}

