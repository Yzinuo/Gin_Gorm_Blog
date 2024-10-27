package handle

import (
	"errors"
	g "gin-blog/internal/global"
	"gin-blog/internal/model"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
type Menu struct{}

type MenuTreeVO struct{
	model.Menu
	Children []MenuTreeVO `json:"children"`
}

type TreeOptionVO struct {
	ID       int            `json:"key"`
	Label    string         `json:"label"`
	Children []TreeOptionVO `json:"children"`
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
func (*Menu)GetTreeList(c *gin.Context){
	keyword := c.Query("keyword")

	menuList,err := model.GetMenuList(GetDB(c),keyword)
	if err != nil{
		ReturnError(c,g.ErrDbOp,err)
		return
	}
	ReturnSuccess(c,menus2MenuVos(menuList))
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
		ReturnError(c,g.ErrMenuUsedByRole,nil)
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

// GetOption 获取菜单的树形结构
func (*Menu)GetOption(c *gin.Context) {
	db := GetDB(c)
	menulist,err := model.GetMenuList(db,"")
	if err != nil{
		ReturnError(c,g.ErrDbOp,err)
		return
	}

	result := make([]TreeOptionVO,0)
	for _,menu := range menus2MenuVos(menulist){
		Tree := TreeOptionVO{ID: menu.ID,Label: menu.Name}
		for _,children := range menu.Children{
			Tree.Children = append(Tree.Children, TreeOptionVO{ID: children.ID,Label: children.Name})
		}
		result = append(result, Tree)	
	}

	ReturnSuccess(c,result)
}

// 生成菜单的树形结构
func menus2MenuVos(menus []model.Menu) []MenuTreeVO{
	// 先筛选出一级菜单，再筛选出子菜单
	// 树形结构即返回 result
	result := make([]MenuTreeVO,0)
	FirstLevel := getFirstLevelMenus(menus)
	Childrenmap := getMenuChildrenMap(menus)

	for _,first := range FirstLevel{
		Tree := MenuTreeVO{Menu: first}
		for _,child := range Childrenmap[first.ID]{
			Tree.Children = append(Tree.Children, MenuTreeVO{Menu: child})
		}		
		delete(Childrenmap,first.ID)
		result = append(result, Tree)
	}
	sortMenu(result)
	return result
}
// 筛选一级菜单  parentId == 0
func getFirstLevelMenus(menuList []model.Menu)[]model.Menu{
	menu := make([]model.Menu,0)
	for _,m := range menuList{
		if m.ParentId == 0{
			menu = append(menu, m)
		}
	}
	return menu
}

func getMenuChildrenMap(menus []model.Menu)map[int][]model.Menu{
	childrenmap := make(map[int][]model.Menu)

	for _,m := range menus{
		if m.ParentId != 0{
			childrenmap[m.ParentId] = append(childrenmap[m.ParentId], m)
		}
	}
	return childrenmap
}
// 排序菜单
func sortMenu(menus []MenuTreeVO){
	sort.Slice(menus,func(i,j int) bool{
		return menus[i].OrderNum < menus[j].OrderNum
	})

	for i := range menus{
		sort.Slice(menus[i].Children,func(j,k int)bool{
			return menus[i].Children[j].OrderNum < menus[i].Children[k].OrderNum
		})
	}
}