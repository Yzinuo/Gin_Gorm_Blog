package model

import (
	"encoding/json"
	"gin-blog/internal/utils"
	"log/slog"
	"strconv"

	"time"

	"gorm.io/gorm"
)

// 权限控制 7 张表（4 模型 + 3 关联）
type UserAuth struct {
	Model
	Username      string     `gorm:"unique;type:varchar(50)" json:"username"` // 邮箱地址
	Password      string     `gorm:"type:varchar(100)" json:"-"`
	LoginType     int        `gorm:"type:tinyint(1);comment:登录类型" json:"login_type"`
	IpAddress     string     `gorm:"type:varchar(20);comment:登录IP地址" json:"ip_address"`
	IpSource      string     `gorm:"type:varchar(50);comment:IP来源" json:"ip_source"`
	LastLoginTime *time.Time `json:"last_login_time"`
	IsDisable     bool       `json:"is_disable"`
	IsSuper       bool       `json:"is_super"` // 超级管理员只能后台设置

	UserInfoId int       `json:"user_info_id"`
	UserInfo   *UserInfo `json:"info"` // 自动推断UserInfoId是外键
	Roles      []*Role   `json:"roles" gorm:"many2many:user_auth_role"`
}

// 转换从json模式 []byte 代表任意形式的json
func (u *UserAuth) MarshalBinary() (data []byte, err error){
	return json.Marshal(u)
}

type Role struct{
	Model
	Name		string		`gorm:"unique" json:"name"`
	Label		string		`gorm:"unique" json:"label"`
	IsDisable	bool		`json:"is_disable"`

	Resources 	[]Resource	`gorm:"many2many:role_resource" json:"resource"`
	Menus		[]Menu		`gorm:"many2many:role_menu" json:"menus"`
	Users		[]UserAuth	`gorm:"many2many:user_auth_role" json:"users"`
}

// 存储着各个模块下的各个path对应的http method （博客后台）
type Resource struct{
	Model
	Name		string		`gorm:"unique;type:varchar(50)" json:"name"`
	ParentId	int			`json:"parent_id"` // 标志着哪个模块下的方法
	Url			string		`gorm:"type:varchar(255)" json:"url"`
	Method		string		`gorm:"type:varchar(10)" json:"request_method"`
	Anonymous	bool		`json:"is_anonymous"` // 允不允许匿名访问
	Roles		[]*Role		`gorm:"many2many:role_resource" json:"roles"`
}

/*
菜单设计:

目录: catalogue === true
  - 如果是目录, 作为单独项, 不展开子菜单（例如 "首页", "个人中心"）
  - 如果不是目录, 且 parent_id 为 0, 则为一级菜单, 可展开子菜单（例如 "文章管理" 下有 "文章列表", "文章分类", "文章标签" 等子菜单）
  - 如果不是目录, 且 parent_id 不为 0, 则为二级菜单

隐藏: hidden
  - 隐藏则不显示在菜单栏中

外链: external, external_link
  - 如果是外链, 如果设置为外链, 则点击后会在新窗口打开
*/

type Menu struct {
	Model
	ParentId     int    `json:"parent_id"`
	Name         string `gorm:"uniqueIndex:idx_name_and_path;type:varchar(20)" json:"name"` // 菜单名称
	Path         string `gorm:"uniqueIndex:idx_name_and_path;type:varchar(50)" json:"path"` // 路由地址
	Component    string `gorm:"type:varchar(50)" json:"component"`                          // 组件路径
	Icon         string `gorm:"type:varchar(50)" json:"icon"`                               // 图标
	OrderNum     int8   `json:"order_num"`                                                  // 排序
	Redirect     string `gorm:"type:varchar(50)" json:"redirect"`                           // 重定向地址
	Catalogue    bool   `json:"is_catalogue"`                                               // 是否为目录
	Hidden       bool   `json:"is_hidden"`                                                  // 是否隐藏
	KeepAlive    bool   `json:"keep_alive"`                                                 // 是否缓存
	External     bool   `json:"is_external"`                                                // 是否外链
	ExternalLink string `gorm:"type:varchar(255)" json:"external_link"`                     // 外链地址

	Roles []*Role `json:"roles" gorm:"many2many:role_menu"`
}

type RoleResource 	struct{
	RoleId			int 	`gorm:"primaryKey;uniqueIndex:idx_role_resource" json:"-" `
	ResourceId		int 	`gorm:"primaryKey;uniqueIndex:idx_role_resource" json:"-" `
}

type UserAuthRole	struct {
	UserAuthId		int			`gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
	RoleId			int			`gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
}

type RoleMenu struct{
	RoleId			int			`gorm:"primaryKey;uniqueIndex:idx_role_menu" json:"-"`
	MenuId			int			`gorm:"primaryKey;uniqueIndex:idx_role_menu" json:"-"`
}

type RoleVO struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Label       string    `json:"label"`
	IsDisable   bool      `json:"is_disable"`
	ResourceIds []int     `json:"resource_ids" gorm:"-"`
	MenuIds     []int     `json:"menu_ids" gorm:"-"`
}

// 传入一个menu，有就更新，没有创造
func SaveOrUpdateMenu(db *gorm.DB,menu *Menu) error{
	var result *gorm.DB

	if menu.ID > 0 {
		result = db.Model(menu).
				Select("name", "path", "component", "icon", "redirect", "parent_id", "order_num", "catalogue", "hidden", "keep_alive", "external").
				Updates(menu)
	}else{
		result = db.Create(menu)
	}
	return result.Error
}

func GetMenuIdsByRoleId(db *gorm.DB,roleId int) (id []int,err error){
	result := db.Model(&RoleMenu{}).
				Where("role_id = ?",roleId).Pluck("menu_id",&id)
	
	return id,result.Error
}

func GetMenuById(db *gorm.DB,id int) (menu *Menu,err error){
	result := db.First(&menu,id)
	return menu,result.Error
}

// 查询某一个菜单能不能用  有没有对应的role有它的权限
func CheckMenuInUse(db *gorm.DB,id int)(bool, error){
	var count int64
	result := db.Model(&RoleMenu{}).Where("menu_id",id).Count(&count)
	return count > 0,result.Error
}

// 检查菜单有没有子菜单
func CheckMenuHasChidren(db *gorm.DB,id int)(bool,error){
	var count int64
	result := db.Model(&Menu{}).Where("parent_id = ?",id).Count(&count)
	return count > 0,result.Error
}

// 超级管理员一键获取所有菜单
func GetAllMenuList(db *gorm.DB) (menuList []Menu, err error){
	result := db.Find(&menuList)
	return menuList,result.Error
}

// 把user_id 和 menu通过 role 联系起来
func GetMenuListByUserId(db *gorm.DB,id int)(menus []Menu,err error){
	var userauth UserAuth
	result := db.Where(&UserAuth{Model: Model{ID: id}}).
				Preload("Roles").Preload("Roles.Menus").
				First(&userauth)
	if result.Error != nil{
		return nil,result.Error
	}


	set := make(map[int]Menu)
	for _,role := range userauth.Roles{
		for _,role_menu := range role.Menus{
			set[role_menu.ID] = role_menu
		}
	}
	
	for _,menu := range set{
		menus = append(menus, menu)
	}

	return menus,nil
}

// 查询包含关键词的munu
func GetMenuList(db *gorm.DB,keyword string) (List []Menu, err error){
	db = db.Model(&Menu{})
	if	keyword != ""{
		db = db.Where("name LIKE ?","%"+keyword+"%")
	}
	var count int64
	result := db.Count(&count).Find(&List)

	return List,result.Error
}

// Resource
func SaveOrUpdateResource(db *gorm.DB, id, pid int, name, url, method string) error{
	resource := Resource{
		Model : Model{ID : id},
		ParentId: pid,
		Name: name,
		Url: url,
		Method: method,
	}

	var result *gorm.DB
	if id > 0{
		result = db.Updates(&resource)
	}else{
		result = db.Create(&resource)
	}

	return result.Error
}

func GetResourceIdsByRoleId(db *gorm.DB,id int) (ids []int, err error){
	result := db.Model(&RoleResource{}).
				Where("role_id",id).
				Pluck("resource_id",&ids)
	return ids,result.Error
}

func GetResourceList(db *gorm.DB,keyword string) (List []Resource,err error){
	db = db.Model(&Resource{})
	if keyword != ""{
		db = db.Where("name LIKE ?","%"+keyword+"%")
	}

	result := db.Find(&List)
	return List,result.Error
}

func GetResourceListByIds(db *gorm.DB,ids []int) (List []Resource,err error){
	result := db.Where("id in ?",ids).Find(&List)
	return List,result.Error
}

// role
func SaveOrUpdateRole(db *gorm.DB, id int, name, label string, isDisable bool) error{
	role := Role{
		Model : Model{ID: id},
		Name: name,
		Label: label,
		IsDisable: isDisable,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&role)
	}else {
		result = db.Create(&role)
	}
	return result.Error
}

// 把Role中的id name都封装到list中
func GetRoleOption(db *gorm.DB) (list []OptionVO, err error){
	db = db.Model(&Role{})
	result := db.Select("id","name").Find(&list)
	if result.Error != nil {
		return list,result.Error
	}

	return list,nil
}

func GetRoleList(db *gorm.DB,num,size int ,keyWord string) (list []RoleVO, total int64,err error){
	db = db.Model(&Role{})
	
	if keyWord != ""{
		db = db.Where("name Like ?", "%"+keyWord+"%")
	}

	result := db.Count(&total).Scopes(Paginate(num,size)).
				Find(&list)
	return list,total,result.Error
}

func GetRoleIdsByUserId(db *gorm.DB, userAuthId int) (ids []int , err error){
	result := db.Model(&UserAuthRole{UserAuthId: userAuthId}).
				Pluck("role_id",&ids)

	return ids,result.Error
}

func SaveRole (db *gorm.DB, name,label string)	error {
	role := Role{
		Name : name,
		Label:label,
	}

	result := db.Create(&role)
	return result.Error
}

// 更新角色 使用事务（设计的表较多，更容易出错）
func UpdateRole(db *gorm.DB,id int, name, label string, isDisable bool, resourceIds, menuIds []int) error {
	role := Role{
		Model : Model{ID : id},
		Name  : name,
		Label:  label,
		IsDisable: isDisable,
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := db.Model(&role).Select("name","label","is_disable").Updates(&role).Error;err != nil{
			return err
		}
		// 清空相关联的资源表和菜单表
		if err := db.Delete(&RoleResource{},"role_id = ?",id).Error; err != nil {
			return err
		}
		if err := db.Delete(&RoleMenu{},"role_id = ?",id).Error; err != nil {
			return err
		}

		// 创造新的
		for _,rid := range resourceIds {
			if err := db.Create(&RoleResource{RoleId: id, ResourceId:rid}).Error; err != nil {
				return err
			}
		}

		for _,mid := range menuIds {
			if err := db.Create(&RoleMenu{RoleId: id, MenuId:mid }).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// 事务删除 role,role_resource,role_menu
func DeleteRoles(db *gorm.DB,ids []int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		result := db.Delete(&Role{},"id in ?",ids)
		if result.Error != nil {
			return result.Error
		}

		result = db.Delete(&RoleResource{},"role_id in ?",ids)
		if result.Error != nil {
			return result.Error
		}

		result = db.Delete(&RoleMenu{},"role_id in ?",ids)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}

// userauth
// 使用复合字面量在堆上创造了内存, 堆数据不会因为函数结束销毁,这个指针仍然有效
func GetUserAuthById (db *gorm.DB, id int) (*UserAuth,error){
	userauth := UserAuth{Model: Model{ID : id}}
	result := db.Model(&userauth).
				Preload("Roles").Preload("UserInfo").
				First(&userauth)
	
	return &userauth,result.Error
}

func CreateNewUser(db *gorm.DB,username, password string) (*UserAuth,*UserInfo,*UserAuthRole,error){
	// 创建userinfo
	num,err := Count(db,&UserInfo{})
	if err != nil{
		slog.Info(err.Error())
	}
	number := strconv.Itoa(num)
	userinfo := &UserInfo{
		Email : username,
		Nickname : "游客"+number,
		Avatar: "https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png",
		Intro: "我是这个程序的第"+number+"个用户",
	}
	result := db.Create(&userinfo)
	if result.Error != nil {
		return nil,nil,nil,result.Error
	}

	// 先创建userauth
	pass ,_:= utils.BcryptHash(password)
	userauth := &UserAuth{
		Username: username,
		Password: pass,
		UserInfoId: userinfo.ID,
	}
	
	result = db.Create(&userauth)
	if result.Error != nil {
		return nil,nil,nil,result.Error
	}

	// 再创建role关联表
	user_role := &UserAuthRole{
		UserAuthId: userauth.ID,
		RoleId: 2,
	}
	result = db.Create(&user_role)
	if result.Error != nil {
		return nil,nil,nil,result.Error
	}	

	return userauth,userinfo,user_role,result.Error
}