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

	Title  		string `gorm:"type:varchar(100);not null" json:"title"`
	Desc   		string `json:"desc"`
	Content		string `json:"content"`
	Img			string `json:"img"`
	Type		int 	`gorm:"type:tinyint;comment:类型(1-原创 2-转载 3-翻译)" json:"type"`
	Status		int  	`gorm:"type:tinyint;comment:状态(1-公开 2-私密)" json:"status"`
	IsTop		bool 	`json:"is_top"`
	IsDelete	bool 	`json:"is_delete"`
	OriginalUrl	string	`json:"original_url"`

	CategoryId int `json:"category_id"`
	UserId     int `json:"-"` // user_auth_id
	
	//指定多对多关系的关联表 ： article_tag  它的article_id标签和本model关联
	Tags		[]*Tag	`gorm:"many2many:article_tag;joinForeignKey:article_id;joinReferences:tag_id" json:"tags"`
	// 定义外键 如Category就是定义CategoryId为指向Catefory的外键
	Category	*Category `gorm:"foreignkey:CategoryId" json:"category"`
	User		*UserAuth `gorm:"foreignkey:UserId" json:"user"`
}

type ArticleTag struct{
	ArticleId		int
	TagId			int
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
	data = &Article{}

	result := db.Preload("Category").Preload("Tags").
			Where(Article{Model : Model{ID : id}}).
			First(&data)
	
	return data,result.Error
}

// 获取第一个可获取的文章（不在回收站并且状态为公开）
func GetBlogArticle(db *gorm.DB, id int) (data *Article, err error){
	data = &Article{}
	
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
			  Select("id","title","img","created_at").
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
				Select("id","title","img","created_at").
				Where("is_delete = 0 AND status = 1").
				Order("created_at desc,id desc").
				Limit(n).
				Find(&list)
	return list,result.Error
}

// 删除文章
func DeleteArticle(db *gorm.DB,ids []int) (int64,error){
	// 删除对应的 tag-article 关联
	result := db.Where("article_id IN ?",ids).Delete(&ArticleTag{})
	if result.Error != nil{
		return 0,result.Error
	}

	// 删除文章
	result = db.Where("id IN ?",ids).Delete(&Article{})
	if result.Error != nil{
		return 0,result.Error
	}

	return result.RowsAffected,nil
}

// 软删除： 改变is_delete
func UpdateArticleSoftDlete(db *gorm.DB,ids []int,isDelete bool) (int64,error){
	result := db.Model(&Article{}).
				Where("id IN ?", ids).
				Update("is_delete",isDelete)
	
	if result.Error != nil {
		return 0,result.Error
	} 
	return result.RowsAffected,nil
}

// 新增/编辑文章, 同时根据 分类名称, 标签名称 维护关联表
func SvaeOrUpdateArticle(db *gorm.DB,article *Article,categoryname string,tagnames []string) error{
	// 如果没有这个分类，Create
	category := Category{Name : categoryname}
	result := db.Model(&Category{}).Where("name",categoryname).FirstOrCreate(&category)
	if result.Error != nil{
		return result.Error
	}

// 新增或更新文章
	if article.ID == 0{
		result = db.Create(article)
	}else{
		result = db.Model(article).Where("id",article.ID).Updates(article)
	}
	if result.Error != nil{
		return result.Error
	}

	// 更新文章后 文章对应的tag需要更新 维护关联表 
	// 清空
	result = db.Delete(&ArticleTag{},"article_id",article.ID)
	if result.Error != nil{
		return result.Error
	}

	// 如果tag表中有没有这个tag，添加
	var articletags []ArticleTag
	for _,tagname := range tagnames {
		newTag := Tag{Name: tagname}
		result = db.Model(&Tag{}).Where("name",tagname).FirstOrCreate(&newTag)
		if result.Error != nil{
			return result.Error
		}

		articletags = append(articletags, ArticleTag{
			ArticleId: article.ID,
			TagId: newTag.ID,
		})
	}
	result = db.Create(&articletags)
	return result.Error
}

func UpdatearticleTop(db *gorm.DB,id int, Istop bool) error{
	result := db.Model(&Article{Model : Model{ID: id}}).Update("is_top",Istop)
	return result.Error
}


func ImportArticle(db *gorm.DB, userAuthId int, title, content, img string) error {
	article := Article{
		Title:   title,
		Content: content,
		Img:     img,
		Status:  STATUS_DRAFT,
		Type:    TYPE_ORIGINAL,
		UserId:  userAuthId,
	}

	result := db.Create(&article)
	return result.Error
}