package model

import(
	"gorm.io/gorm"
)

type Category struct{
	Model
	Name	string		`gorm:"type:varchar(20);not null;unique" json:"name"`
	Article []Article	`gorm:"foreignKey:CategoryId"`
}