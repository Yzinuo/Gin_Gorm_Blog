package model

import(
	"gorm.io/gorm"
	"time"
)

type Model struct{
	ID			int 		`gorm:"primary_key;auto_increament" json: "id"`
	CreateAt 	time.Time	`json:"create_at"`
	UPdateAt	time.Time	`json:"update_at"`
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