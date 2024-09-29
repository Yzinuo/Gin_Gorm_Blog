package model

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct{
	Model
	Name	string		`gorm:"type:varchar(20); unique; not null"`	
	Article	[]*Article	`gorm:"many2many:article_tag;joinForeignKey:tag_id;joinReferences:article_id" json:"articles,omitempty"`
}

type TagVO struct{
	ID			int			 `json:"id"`
	Created_at	time.Time	 `json:"created_at"`
	Updated_at	time.Time	 `json:"updated_at"`

	Name		string		 `json:"name"`
	ArticleCount int		 `json:"article_count"`
}


func GetTagList(db *gorm.DB,page,size int, keyword string) (list []TagVO,total int64,err error){
	// SELECT t.id,t.name,COUNT(t.id) as articlecount,t.created_at,t.updated_at FROM tag t LEFT JOIN article_tag at ON t.id = at.tagid
	db = db.Table("tag t").
			Select("t.id","t.name","COUNT(t.id) as articleCount","t.created_at","t.updated_at").
			Joins("LEFT JOIN article_tag at ON t.id = at.tag_id")

	if keyword != ""{
		db = db.Where("name LIKE ?","%"+keyword+"%")
	}

	result := db.Group("t.id").
				Order("created_at").
				Count(&total).
				Scopes(Paginate(page,size)).
				Find(&list)

	return list,total,result.Error
}

func GetTagOption (db *gorm.DB) ([]OptionVO,error){
	option := make([]OptionVO,0)
	result := db.Model(&Tag{}).Select("id","name").Find(&option)

	return option,result.Error
}

func GetNamesByArticleId(db *gorm.DB,id int) ([]string, error){
	list := make([]string,0)
	
	result := db.Table("tag").
		Joins("LEFT JOIN article_tag ON tag.id = article_tag.tag_id").
		Where("article_id", id).
		Pluck("name", &list)

	return list,result.Error
}


func SaveOrCreateTag(db *gorm.DB,id int,name string) (*Tag,error){
	tag := Tag{
		Model: Model{ID : id},
		Name: name,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&tag)
	}else{
		result = db.Create(&tag)
	}

	return &tag,result.Error
}
