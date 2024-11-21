package model

import(
	"gorm.io/gorm"
)

type Category struct{
	Model
	Name	string		`gorm:"type:varchar(20);not null;unique" json:"name"`
	Article []Article	`gorm:"foreignKey:CategoryId"`
}

type CategoryVO struct{
	Category
	ArticleCount  int	`json:"article_count"`
}

func GetCategoryList(db *gorm.DB,num ,size int, keyword string)(list []CategoryVO,total int64, err error){
	// 涉及到聚合函数，连表查询都需要用table join这类   model find等等用于简单查询
	// SELECT c.id,c.name,COUNT(c.id) as article_count,c.created_at,c.updated_at FORM category c JOIN article a ON c.id = a.category_id AND a.is_delete = 0 AND a.status = 1
	db = db.Table("category c").
			Select("c.id","c.name","COUNT(c.id) as article_count","c.created_at","c.updated_at").
			Joins("LEFT JOIN article a ON c.id = a.category_id AND a.is_delete = 0 AND a.status = 1")

	if keyword != ""{
		db = db.Where("name LIKE ?","%"+"%")
	}

	result := db.Group("c.id").
			  Order("c.created_at DESC").
			  Count(&total).
			  Scopes(Paginate(num,size)).
			  Find(&list)
	
	return list,total,result.Error
}

func GetcategoryOption (db *gorm.DB) (list []OptionVO,err error){
	result := db.Model(&Category{}).Select("id","name").Find(&list)
	return list,result.Error
}

func GetCategoryByname(db *gorm.DB,name string) (cg *Category,err error){
	cg = &Category{}

	result := db.Model(&Category{}).Where("name LIKE ?",name).Find(cg)
	return cg,result.Error
}

func GetCategoryById(db gorm.DB,id int)(cg *Category,err error){
	cg = &Category{}
	
	result := db.Model(&Category{}).Where("id = ?",id).Find(cg)
	return cg,result.Error
}

func DeleteCategory(db *gorm.DB,ids []int)(count int64,err error){
	result := db.Delete(&Category{},"id IN ?",ids)
	if result.Error != nil{
		return 0,result.Error
	}
	return db.RowsAffected,nil
}

func SaveOrUpdateCategory(db *gorm.DB,id int, name string) (*Category,error){
	category := Category{
		Model: Model{ID: id},
		Name: name,
	}

	var result *gorm.DB
	if id > 0{
		result = db.Updates(&category)
	}else{
		result = db.Create(&category)
	}

	return &category,result.Error
}