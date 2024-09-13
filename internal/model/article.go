package model

import (
	"time"

	"gorm.io/gorm"
)

const(
	STATUS_PUBLIC = iota+1 // 文章状态，公开
	STATUS_SERCET //私有
	STATUS_DRAFT // 草稿
)

const(
	TYPE_ORIGINAL = iota +1 // 文章属性，原创
	TYPE_REPRINT // 转载
	TYPE_TRANSLATE // 翻译
)

// 定义model 和数据库中的表对应
type Article struct{
	Model

	Title  		string `grom:"type:varchar(100);not null" json:"title"`
	Desc   		string `json:"desc"`
	Content		string `json:"content"`
	Img			string `json:"img"`
	Type		int 	`gorm:"type:tinyint;comment:类型(1-原创 2-转载 3-翻译)" json:"type"`
	Status		int  	`grom:"type:tinyint;comment:状态(1-公开 2-私密)" json:"status"`
	IsTop		bool 	`json:"is_top"`
	IsDelete	bool 	`json:"is_delete"`
	OriginalUrl	string	`json:"original_url"`

	CategoryId int `json:"category_id"`
	UserId     int `json:"-"` // user_auth_id
	
	//指定多对多关系的关联表 ： article_tag  它的article_id标签和本model关联
	Tags		[]*Tag	`gorm:"manytomany:article_tag;joinForeignKey:article_id" json:"tags"`
	// 定义外键 如Category就是定义CategoryId为指向Catefory的外键
	Category	*Category `gorm:"foreignkey:CategoryId" json:"category"`
	User		*UserAuth `gorm:"foreignkey:UserId" json:"User"`
}

type ArticleTag struct{
	ArticleId		string
	TagId			string
}

type ArticlePaginationVO struct {
	ID    int    `json:"id"`
	Img   string `json:"img"`
	Title string `json:"title"`
}

type RecommendArticleVO struct {
	ID        int       `json:"id"`
	Img       string    `json:"img"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type BlogArticleVO struct {
	Article

	CommentCount int64 `json:"comment_count"` // 评论数量
	LikeCount    int64 `json:"like_count"`    // 点赞数量
	ViewCount    int64 `json:"view_count"`    // 访问数量

	LastArticle       ArticlePaginationVO  `gorm:"-" json:"last_article"`       // 上一篇
	NextArticle       ArticlePaginationVO  `gorm:"-" json:"next_article"`       // 下一篇
	RecommendArticles []RecommendArticleVO `gorm:"-" json:"recommend_articles"` // 推荐文章
	NewestArticles    []RecommendArticleVO `gorm:"-" json:"newest_articles"`    // 最新文章
}

func GetArticle(db *gorm.DB,id int) (data *Article,err error){
	result := db.Preload("Category").Preload("Tags").
			Where(Article{Model : Model{ID : id}}).
			First(&data)
	
	return data,result.Error
}

// 获取第一个可获取的文章（不在回收站并且状态为公开）
func GetBlogArticle(db *gorm.DB, id int) (data *Article, err error){
	result := db.Preload("Category").Preload("Tags").
			Where(Article{Model: Model{ID : id}}).
			Where("is_delete = 0 AND status = 1").
			First(&data)
	return data,result.Error
}

// 首页文章列表 
func GetBlogArticleList(db *gorm.DB,page,size,CategoryId,TagId int) ( data []Article,err error,total int64){
	db = db.Model(Article{})
	db = db.Where("is_delete = 0 AND status = 1")

	if CategoryId != 0 {
		db = db.Where("catagory_id = ?",CategoryId)
	}
	if TagId != 0 {
		db = db.Where("id IN (SELECT article_id FROM article_tag WHERE tag_id = ?)",TagId)
	}

	db = db.Count(&total)
	result := db.Preload("Category").Preload("Tags").
				Order("is_top desc,id desc").
				Scopes(Paginate(page,size)).
				Find(&data)
	
	return data,result.Error,total
}

func GetArticleList(db *gorm.DB, page, size int, title string, isDelete *bool, status, typ, categoryId, tagId int) (List []Article,total int64,err error){
	db = db.Model(Article{})

	if title != ""{
		db = db.Where("title LIKE ?","%" + title + "%")
	}
	if isDelete != nil {
		db = db.Where("is_delete",isDelete)
	}
	if status != 0 {
		db = db.Where("status",status)
	}
	if typ != 0 {
		db = db.Where("type",typ)
	}
	if categoryId != 0{
		db = db.Where("category_id",categoryId)
	}

	db = db.Preload("Category").Preload("Tags").
			Joins("LEFT JOIN article_tag ON article_tag.article_id = article.id").
			Group("id")
	if tagId != 0 {
		db = db.Where("tag_id = ?",tagId)
	}

	result := db.Count(&total).
				Scopes(Paginate(page,size)).
				Order("is_top DESC,id DESC").
				Find(&List)
	
	return List,total,result.Error
}

// 根据当前的标签，推荐文章
func GetRecommandList(db *gorm.DB,id,n int)(list []RecommendArticleVO,err error){
	// sub1: 查出对应标签列表
	// SELECT tag_id FROM `article_tag` WHERE `article_id` = ?
	sub1 := db.Table("article_tag").
			Select("tag_id").
			Where("article_id",id)

	// sub2: 根据本文章的Tag，找到有对应tag的文章id
	//SELECT DISTINCT article_id FROM `sub1 t` 
	//JOIN `article_tag ON  article_tag.tag_id = t.tag_id`
	// WHERE `aticle_id != ?`
	sub2 := db.Table("(?) t",sub1).
			Select("DISTINCT article_id").
			Joins("JOIN article_tag ON article_tag.id = t.tag_id").
			Where("article_id != ?",id)

	// 更据得到的文章id 去Article数据库中找到对应信息
	result := db.Table("(?) t1",sub2).
			  Select("id,title,img,create_at").
			  Joins("JOIN article ON article.id = t1.article_id").
			  Where("is_delete",0).
			  Order("is_top desc, id desc").
			  Limit(n).
			  Find(&list)

	return list,result.Error
}

// 查询上一篇文章
func GetLastArticle(db *gorm.DB,id int) (list ArticlePaginationVO,err error){
	//Select max(id) FROM article WHERE id < id
	sub1 := db.Table("article").
			Select("max(id)").
			Where("id < ?",id)
	
	// SELECT `id,img,title` FROM `article` WHERE `id = sub1`
	result := db.Table("article").
			Select("id,img,title").
			Where("id = ? AND is_delete = 0 AND status = 1",sub1).
			Limit(1).
			Find(&list)
	
	return list,result.Error
}

// 查询下一个文章
func GetNextArticle(db *gorm.DB,id int) (list ArticlePaginationVO,err error){
	//Select min(id) FROM article WHERE id > id
	sub1 := db.Table("article").
			Select("min(id)").
			Where("id > ?",id)
	
	// SELECT `id,img,title` FROM `article` WHERE `id = sub1`
	result := db.Table("article").
			Select("id,img,title").
			Where("id = ? AND is_delete = 0 AND status = 1",sub1).
			Limit(1).
			Find(&list)
	
	return list,result.Error
}

func GetNewestList(db *gorm.DB, n int) (list []RecommendArticleVO,err error){
	result := db.Model(&Article{}).
				Select("id,title,img,create_at").
				Where("is_delete = 0 AND status = 1").
				Order("create_at desc,id desc").
				Limit(n).
				Find(&list)
	return list,result.Error
}

