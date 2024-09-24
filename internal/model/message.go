package model

import "gorm.io/gorm"

type Message struct{
	Model
	Nickname		string 			`gorm:"type:varchar(50);comment:昵称" json:"nickname"`
	Avatar			string			`gorm:"type:varchar(255);comment:头像地址" json:"avatar"`
	Content			string			`gorm:"type:varchar(255);comment:内容" json:"content"`
	IPAddress		string			`gorm:"type:varchar(255);comment:IP地址" json:"ip_address"`
	IPSource		string			`gorm:"type:varchar(255);comment:IP源" json:"ip_source"`
	Speed			int				`gorm:"type:tinyint(1);comment:弹幕速度" json:"speed"`
	Is_review		bool 			`json:"is_review"`	
}

// 评论的增删改查
func CreatenewMessage(db *gorm.DB,nickname,avatar,content,address,source string,speed int,is_review bool) (*Message,error) {
	message := Message{
		Nickname: nickname,
		Avatar: avatar,
		Content: content,
		IPAddress: address,
		IPSource: source,
		Speed: speed,
		Is_review: is_review,
	}	

	result := db.Create(&message)
	return &message,result.Error
}

func DeleteMessage(db *gorm.DB,ids []int) (int64 ,error){
	result := db.Where("id in ?",ids).Delete(&Message{})
	
	return db.RowsAffected,result.Error
}

// 查某个nickname 的评论
func GetMessageList(db *gorm.DB,page,size int,keyword string)(List []Message,total int64,err error){
	db = db.Model(&Message{})

	if keyword != ""{
		db = db.Where("nickname LIKE ?","%"+keyword+"%")
	}

	result := db.Count(&total).
			Order("created_at DESC").
			Scopes(Paginate(page,size)).
			Find(&List)

	return List,total,result.Error
}

func UpdateMessageReview(db *gorm.DB,ids []int, is_review bool) (int64, error){
	result := db.Model(&Message{}).
				Where("id in ?",ids).
				Update("is_review",is_review)
	
	return db.RowsAffected,result.Error 
}






