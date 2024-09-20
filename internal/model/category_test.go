package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategory(t *testing.T){
	db := setup(t)

	category := Category{
		Model: Model{ID: 1},
		Name:"后端",
		Article: []Article{
			{CategoryId: 1,IsDelete: false,Status: 1},
			{CategoryId: 1,IsDelete: false,Status: 0},
		},
	}

	db.Create(&category)

	list,total,err := GetCategoryList(db,1,10,"后")
	assert.Nil(t,err)
	assert.Equal(t,len(list),1)
	assert.Equal(t,total,int64(1))

	cg,erro := GetCategoryByname(db,"后端")
	assert.Nil(t,erro)
	assert.Equal(t,category.ID,cg.ID)

	err = SaveOrUpdateCategory(db,1,"GOlang")
	assert.Nil(t,err)

}