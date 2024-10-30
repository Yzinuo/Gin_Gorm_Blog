package handle

import (
	"gin-blog/internal/model"
	g "gin-blog/internal/global"
	"github.com/gin-gonic/gin"
)

type  Role struct{}

type AddOrEditRoleReq struct{
	ID     	  int 			`json:"id"`
	Name      string 		`json:"name"`
	Label     string 		`json:"label"`
	IsDisable bool 	  		`json:"is_disable"`
	ResourceIds []int       `json:"resource_ids"`
	MenuIds   []int         `json:"menu_ids"`
}

// 获取各个角色的Id 和 名字
func (*Role) GetOption(c *gin.Context){
	data,err := model.GetRoleOption(GetDB(c))
	if err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	ReturnSuccess(c,data)
}

// 根据特定的查询条件获取角色列表 roleVO还需要分配Menu和Resource
func (*Role) GetList(c *gin.Context){
	var query PageQuery
	db := GetDB(c)
	if err := c.ShouldBindQuery(&query); err!= nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	data,total,err := model.GetRoleList(db,query.Page,query.Size,query.Keyword)
	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	result := make([]model.RoleVO,0)
	for _,role := range data {
		role.ResourceIds,_ = model.GetResourceIdsByRoleId(db,role.ID)
		role.MenuIds,_ = model.GetMenuIdsByRoleId(db,role.ID)
		result = append(result,role)
	}

	ReturnSuccess(c,PageList[model.RoleVO]{
		Total: total,
		Data: result,
		Page: query.Page,
		Size: query.Size,
	})
}

// 更新或添加角色
// 有id 则更新，否则添加
func (*Role) SaveOrUpdate(c *gin.Context){
	var req AddOrEditRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	db := GetDB(c)
	// 有id则更新，否则添加
	if req.ID == 0 {
		err := model.SaveRole(db,req.Name,req.Label)
		if err != nil {
			ReturnError(c,g.ErrDbOp,err)
			return
		}
	}else {
		err := model.UpdateRole(db,req.ID,req.Name,req.Label,req.IsDisable,req.ResourceIds,req.MenuIds)
		if err!= nil {
			ReturnError(c,g.ErrDbOp,err)
			return
		}
	}

	ReturnSuccess(c,nil)
}

// 删除角色
func (*Role) Delete (c *gin.Context){
	var ids []int 
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	db := GetDB(c)
	err := model.DeleteRoles(db,ids)
	if err != nil{
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	
	ReturnSuccess(c,nil)
}