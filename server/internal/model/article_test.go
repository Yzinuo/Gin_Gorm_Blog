package model

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setup(t *testing.T) *gorm.DB{

	
	db,err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //
		},})
	if err != nil {
		t.Fatalf("fail to connect to the database")
	}
	
	err = MakeMigrate(db)
	if err != nil {
		t.Fatalf("fail to migrate data")
	}

	return db
}

func TestGetBlogArticleList(t *testing.T){
	db := setup(t)

	articles := []Article{
		{Title: "Article 1", Content: "Content 1", Status: STATUS_PUBLIC, IsDelete: false},
		{Title: "Article 2", Content: "Content 2", Status: STATUS_PUBLIC, IsDelete: false},
	}

	for _,article := range articles{
		db.Create(&article)
	}

	data,err,total := GetBlogArticleList(db,1,10,0,0)
	assert.NoError(t,err)
	assert.Equal(t, int64(len(articles)), total)
	assert.Equal(t, len(articles), len(data))
}