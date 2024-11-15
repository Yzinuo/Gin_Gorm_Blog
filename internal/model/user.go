package model

import (
	"errors"

	"gorm.io/gorm"
)

type UserInfo struct {
	Model
	Email    string `json:"email" gorm:"type:varchar(30)"`
	Nickname string `json:"nickname" gorm:"unique;type:varchar(30);not null"`
	Avatar   string `json:"avatar" gorm:"type:varchar(1024);not null"`
	Intro    string `json:"intro" gorm:"type:varchar(255)"`
	Website  string `json:"website" gorm:"type:varchar(255)"`
}

type UserInfoVO struct {
	UserInfo
	ArticleLikeSet			[]string 		`json:"article_like_set"`
	CommentLikeSet			[]string		`json:"comment_like_set"`
}

func GetUserInfoById(db *gorm.DB,id int) (*UserInfo,error){
	var userinfo UserInfo
	result := db.Model(&userinfo).Where("id = ?",id).First(&userinfo)

	return &userinfo,result.Error
}

func GetUserAuthInfoByName(db *gorm.DB,name string) (*UserAuth,error){
	var userauth UserAuth
	
	result := db.Model(&userauth).Where("username LIKE ?",name).First(&userauth)
	if result.Error != nil && errors.Is(result.Error,gorm.ErrRecordNotFound){
		return nil,result.Error
	}
	
	return &userauth,result.Error
}

func GetUserList(db *gorm.DB,login_type int8,page,size int,nickname string,name string) (list []UserAuth, total int64, err error){
	if login_type != 0{
		db =db.Where("login_type = ?",login_type)
	}
	if name != "" {
		db =db.Where("username LIKE ?","%"+name+"%")
	}

	result := db.Model(&UserAuth{}).
				Joins("LEFT JOIN user_info ON user_info.id = user_auth.user_info_id").
				Where("user_info.nickname LIKE ?","%"+nickname+"%").
				Count(&total).
				Preload("UserInfo").Preload("Roles").
				Scopes(Paginate(page,size)).
				Find(&list)
	
	return list,total,result.Error
}

//更新用户的昵称和角色  更新角色先清空关连表，再添加
func UpdateUserNicknameAndRole(db *gorm.DB, authId int, nickname string, roleIds []int) error{
	auth ,err:= GetUserAuthById(db,authId)
	if err != nil {
		return err
	}

	userinfo := UserInfo{
		Model: Model{ ID: auth.UserInfoId},
		Nickname: nickname,
	}
	
	result := db.Updates(&userinfo)
	if result.Error != nil {
		return result.Error
	}

	
	
	if len(roleIds) == 0 {
		return nil
	}
	// 删除关联表的相关记录
	result = db.Model(&UserAuthRole{}).Where("userauth_id = ?",authId).Delete(&UserAuthRole{})
	if result.Error != nil {
		return result.Error
	}

	var userauth_role []UserAuthRole
	for _,roleId := range roleIds {
		userauth_role = append(userauth_role, UserAuthRole{
			UserAuthId: authId,
			RoleId: roleId,
		})
	}

	result = db.Create(&userauth_role)

	return result.Error
}

func UpdateUserPassword(db *gorm.DB, id int, password string) error{
	userauth := UserAuth{
		Model: Model{ID: id},
	}

	result := db.Model(&userauth).Update("password",password)

	return result.Error
}

func UpdateUserInfo(db *gorm.DB, id int, nickname, avatar, intro, website string) error{
	userinfo := UserInfo{
		Model: Model{ID: id},
		Nickname: nickname,
		Avatar: avatar,
		Intro: intro,
		Website: website,
	}

	result := db.Model(&userinfo).Updates(&userinfo)
	return result.Error
}


func UpdateUserDisable(db *gorm.DB, id int, isDisable bool) error {
	userauth := UserAuth{
		Model: Model{ID: id},
		IsDisable: isDisable,
	}
	
	result := db.Model(&userauth).Updates(&userauth)
	return result.Error
}


func UpdateUserLoginInfo(db *gorm.DB, id int, ipAddress, ipSource string) error {
	userauth := UserAuth{
		Model: Model{ID: id},
		IpAddress: ipAddress,
		IpSource: ipSource,
	}

	result := db.Model(&userauth).Updates(&userauth)
	return result.Error
}


