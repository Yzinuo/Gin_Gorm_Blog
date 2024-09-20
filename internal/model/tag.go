package model

type Tag struct{
	Model
	Name	string		`gorm :"type:varchar(20);unique;not null"`	
	Article	[]*Article	`gorm:"many2many:article_tag;joinForeignKey:tag_id;joinReferences:article_id" json:"articles,omitempty"`
}