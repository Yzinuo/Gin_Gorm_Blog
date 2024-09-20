package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCommentList(t *testing.T) {
	db := setup(t)

	user := UserAuth{
		Username: "username",
		Password: "123456",
		UserInfo: &UserInfo{
			Nickname: "nickname",
		},
	}
	db.Create(&user)

	article := Article{Title: "title", Content: "content"}
	db.Create(&article)

	comment, _ := AddComment(db, user.ID, TYPE_ARTICLE, article.ID, "content", true)
	_, _ = AddReplyComment(db, user.ID, user.ID, comment.ID, "reply_content", true)

	data, total, err := GetCommentList(db, nil, 1,10, TYPE_ARTICLE, "")
	assert.Nil(t, err)
	assert.Equal(t, 2, int(total))
	assert.Equal(t, "reply_content", data[0].Content)
	assert.Equal(t, "content", data[1].Content)

	v1 := data[0]
	assert.Equal(t, "reply_content", v1.Content)
	assert.Equal(t, "username", v1.User.Username)               // preload userAuth
	assert.Equal(t, "nickname", v1.User.UserInfo.Nickname)      // preload userAuth.userInfo
	assert.Equal(t, "username", v1.ReplyUser.Username)          // preload replyUser
	assert.Equal(t, "nickname", v1.ReplyUser.UserInfo.Nickname) // preload replyUser.userInfo
	assert.Equal(t, "title", v1.Article.Title)                  // preload article
}