package model

import (
	"gorm.io/gorm"
)

const (
	TYPE_ARTICLE = iota + 1 // 文章
	TYPE_LINK				//友链
	TYPE_TALK				//说说
)

/*
如果评论类型是文章，那么 topic_id 就是文章的 id
如果评论类型是友链，不需要 topic_id
*/

type Comment struct{
	Model
	UserId			int			`json:"user_id"`
	ReplyUserId		int			`json:"reply_user_id"`
	TopicId			int			`json:"topic_id"`
	ParentId		int			`json:"parent_id"`
	Content			string		`gorm:"type:varchar(500);not null" json:"content"`
	Type       		int    		`gorm:"type:tinyint(1);not null;comment:评论类型(1.文章 2.友链 3.说说)" json:"type"` // 评论类型 1.文章 2.友链 3.说说
	IsReview   	    bool   		`json:"is_review"`

	User			*UserAuth   `gorm:"foreignKey:UserId" json:"user"`
	ReplyUser		*UserAuth	`gorm:"foreignKey:ReplyUserId" json:"reply_user"`
	Article			*Article	`gorm:"foreignKey:TopicId" json:"article"`
}


type CommentVO struct {
	Comment
	LikeCount  int         `json:"like_count" gorm:"-"`
	ReplyCount int         `json:"reply_count" gorm:"-"`
	ReplyList  []CommentVO `json:"reply_list" gorm:"-"`
}

func AddComment (db *gorm.DB,userId,typ,topicId int,content string,is_review bool) (*Comment,error){
	comment := Comment{
		UserId: userId,
		Type: typ,
		TopicId: topicId,
		Content: content,
		IsReview: is_review,
	}

	result := db.Create(&comment)
	return &comment,result.Error
}

func AddReplyComment(db *gorm.DB,userId,replyuserID,parent_id int,content string,is_review bool) (*Comment,error){
	var com Comment
	result := db.Where("id = ?",parent_id).First(&com)
	if result.Error != nil {
		return nil,result.Error
	}
	
	comment := Comment{
		UserId: userId,
		ReplyUserId: replyuserID,
		Content: content,
		IsReview: is_review,
		ParentId: parent_id,
		Type: com.Type,  // 和父文章一样
		TopicId: com.TopicId,//和父文章一样
	}

	result = db.Create(&comment)
	return &comment,result.Error
}

func GetCommentList(db *gorm.DB,is_review *bool,page,size,typ int,nickname string) (list []Comment,count int64,err error){
	var uid int
	
	if nickname != ""{
		result := db.Model(&UserInfo{}).Where("nickname LIKE ?",nickname).Pluck("id",&uid)
		if result.Error != nil {
			return nil,0,db.Error
		}

		db = db.Where("user_id = ?",uid)
	}

	if is_review != nil {
		db = db.Where("is_review = ?",*is_review)
	}

	if typ != 0 {
		db = db.Where("type = ?",typ)
	}

	result := db.Model(&Comment{}).Count(&count).
				Preload("User").Preload("User.UserInfo").
				Preload("ReplyUser").Preload("ReplyUser.UserInfo").
				Preload("Article").
				Order("id DESC").
				Scopes(Paginate(page,size)).
				Find(&list)
	if result.Error != nil {
		return nil,0,result.Error
	}

	return list,count,nil
}

// 得到文章的评论列表
func GetBlogCommentList(db *gorm.DB,page,size, typ ,topic_id int)([]CommentVO,int64,error){
	var list []Comment
	var count int64

	db = db.Model(&Comment{})
	if typ != 0{
		db = db.Where("type = ?",typ)
	}

	if topic_id != 0{
		db = db.Where("topic_id = ?",topic_id)
	}

	// 先找到所有的父评论
	result := db.Where("parent_id = 0").Count(&count).
				Preload("User").Preload("User.UserInfo").
				Order("id DESC").
				Scopes(Paginate(page,size)).
				Find(&list)
	if result.Error != nil {
		return nil,0,result.Error
	}

	var data []CommentVO
	// 在找到每一个对应的子评论
	for _,com := range list{
		replyList := make([]CommentVO, 0)
		result := db.Where("parent_id = ?",com.ID).
					Preload("User").Preload("User.UserInfo").
					Preload("ReplyUser").Preload("ReplyUser.UserInfo").
					Order("id desc").Scopes(Paginate(page,size)).Find(&replyList)
		if result.Error != nil {
			return nil,0,result.Error
		}
		
		data = append(data, CommentVO{
				Comment: com,
				ReplyList: replyList,
				ReplyCount: len(replyList),	
		})
	}

	return  data,count,nil
}

// 更具评论id获取回复列表 
func	GetCommentByid(db *gorm.DB,page,size,id int) (data []Comment, err error){
	result := db.Model(&Comment{}).Where("parent_id = ?",id).
				Preload("User").Preload("User.UserInfo").
				 Order("id DESC").Scopes(Paginate(page,size)).
				 Find(&data)

	return data,result.Error
}

//获取这个文章的评论数
func GetCommentCountOfArticle(db *gorm.DB,topicId int) ( total int64,err error){
	result := db.Model(&Comment{}).Where("topic_id = ? AND type = 1 AND is_review = 1",topicId).Count(&total)
	if result.Error != nil {
		return 0,result.Error
	}
	return total,nil
} 

