package model

import "gorm.io/gorm"


type FriendLink struct{
	Model
	Name		string		`gorm:"type:varchar(255)"`
	Avatar		string		`gorm:"type:varchar(255)"`
	Link		string		`gorm:"type:varchar(255)"`
	Intro		string		`gorm:"type:varchar(255)"`
}

func GetLinkList(db *gorm.DB,page,size int,keyword string) (list []FriendLink,total int64, err error){
	db = db.Model(&FriendLink{})
	
	if keyword != ""{
		db = db.Where("name LIKE ?","%"+keyword+"%").Or("avator LIKE ? ","%"+keyword+"%").Or("intro LIKE ?","%"+keyword+"%")
	}
	
	result := db.Count(&total).Order("created_at DESC").
				Scopes(Paginate(page,size)).
				Find(&list)
	
	return list,total,result.Error
}


func SvaeOrCreateLink(db *gorm.DB,id int,name,avatar,link,intro string) (*FriendLink,error){
	friendlink := FriendLink{
		Model : Model{ID: id},
		Name: name,
		Avatar: avatar,
		Link: link,
		Intro: intro,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&friendlink)
	}else{
		result = db.Create(&friendlink)
	}

	return &friendlink,result.Error
}