package handle

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
type Menu struct{}

type MenuTreeVO struct{
	model.Menu
	Children []MenuTreeVO `json:"children"`
}

// 生成当前用户后台管理界面的菜单
func (*Menu) GetUserMenu(c *gin.Context){
	db := GetDB(c)
	auth,_ := CurrentUserAuth(c)
	menu := []model.Menu{}
	var err error

	if auth.IsSuper {
		menu,err = model.GetAllMenuList(db)	
	}else{
		menu,err = model.GetMenuListByUserId(db,auth.ID)
	}

	if err != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	
	ReturnSuccess(c,menus2MenuVos(menu))
}

// 关键词查找对应菜单
func (*Menu)GetMenuList(c *gin.Context){
	keyword := c.Query("keyword")

	menuList,err := model.GetMenuList(GetDB(c),keyword)
	if err != nil{
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	ReturnSuccess(c,menuList)
}

// 新增或修改菜单
func (*Menu) SaveOrUPdateMenu(c *gin.Context){
	var req model.Menu
	if err := c.ShouldBindJSON(&req); err != nil{
		ReturnError(c,g.ErrRequest,err)
		return
	}

	err := model.SaveOrUpdateMenu(GetDB(c),&req)
	if  err != nil {
		ReturnError(c,g.ErrDbOp,err)
	}
	ReturnSuccess(c,nil)
}

// 删除菜单
func (*Menu)Delete(c *gin.Context){
	db := GetDB(c)

	menuId,err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c,g.ErrRequest,err)
		return
	}

	// 查询菜单是否能被使用
	use,_ := model.CheckMenuInUse(db,menuId)
	if use {
		ReturnError(c,g.ErrRequest,nil)
		return
	}

	menu,err := model.GetMenuById(db,menuId)
	if err != nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			ReturnError(c,g.ErrMenuNotExist,err)
			return
		}
		ReturnError(c,g.ErrDbOp,err)
	}
	// 菜单是否有子菜单 有的话不能删除
	if menu.ParentId == 0 {
		if has,_ := model.CheckMenuHasChidren(db,menuId); has {
			ReturnError(c,g.ErrMenuHasChildren,nil)
			return
		}
	}

	 result := db.Delete(&model.Menu{},menuId)
	 if result.Error != nil {
		ReturnError(c,g.ErrDbOp,err)
		return
	 }

	 ReturnSuccess(c,nil)
}

// 生成菜单的树形结构
func menus2MenuVos(menus []model.Menu)
// 筛选一级菜单  parentId == 0
func getFirstLevelMenus(menuList []model.Menu)

func getMenuChildrenMap(menus []model.Menu)
// 排序菜单
func sortMenu(menus []MenuTreeVO)