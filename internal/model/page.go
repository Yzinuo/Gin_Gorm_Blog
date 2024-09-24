// 这个model的主要目的是把页面的 name 和 url 匹配起来
package model

import "gorm.io/gorm"


type Page struct {
	Model
	Name  string `gorm:"unique;type:varchar(20)" json:"name"`
	Label string `gorm:"unique;type:varchar(30)" json:"label"`
	Cover string `gorm:"type:varchar(255)" json:"cover"`
}

// 得到所有的 page
func GetPageList(db *gorm.DB) (list []Page,total int64,err error){
	result :=  db.Model(&Page{}).Count(&total).Find(&list)
	return list,total,result.Error
}

func SaveOrCreatePage(db *gorm.DB,id int,name,label,cover string) (*Page,error){
	page := Page{
		Model: Model{ID: id},
		Name: name,
		Label: label,
		Cover: cover,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&page)
	}else{
		result = db.Create(&page)
	}

	return &page,result.Error
}

