package model

import(
	"gorm.io/gorm"
	"time"
)

type Model struct{
	ID			int 		`gorm:"primary_key;auto_increament" json: "id"`
	CreatedAt 	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
}

type OptionVO struct {
	ID			string		`json:"id"`
	Name		string 		`json:"name"`
}

// 返回一个闭包函数，因为闭包函数能实现延迟调用（能契合查询时的顺序要求），动态生成函数
func Paginate(page, size int) func(db *gorm.DB) *gorm.DB{
	return func (db *gorm.DB) *gorm.DB{
		if page <=0{
			page = 1
		}
		switch {
		case  size >= 100:
			size = 100	
		case  size <= 10 : 
			size = 10
		}
		offset := (page - 1) *size
		return db.Offset(offset).Limit(size)
	}
}

func MakeMigrate(db *gorm.DB) error {
	// 设置表关联
	db.SetupJoinTable(&Role{}, "Menus", &RoleMenu{})
	db.SetupJoinTable(&Role{}, "Resources", &RoleResource{})
	db.SetupJoinTable(&Role{}, "Users", &UserAuthRole{})
	db.SetupJoinTable(&UserAuth{}, "Roles", &UserAuthRole{})

	return db.AutoMigrate(
		&Article{},      // 文章
		&Category{},     // 分类
		&Tag{},          // 标签
		&Comment{},      // 评论
		// &Message{},      // 消息
		// &FriendLink{},   // 友链
		// &Page{},         // 页面
		&Config{},       // 网站设置
		// &OperationLog{}, // 操作日志
		&UserInfo{},     // 用户信息

		&UserAuth{},     // 用户验证
		&Role{},         // 角色
		&Menu{},         // 菜单
		&Resource{},     // 资源（接口）
		&RoleMenu{},     // 角色-菜单 关联
		&RoleResource{}, // 角色-资源 关联
		&UserAuthRole{}, // 用户-角色 关联
	)
}