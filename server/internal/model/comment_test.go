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

	data, total, err := GetCommentList(db, nil, 1, 10, TYPE_ARTICLE, "")
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

func TestGetBlogCommentListHidesUnreviewedComments(t *testing.T) {
	db := setup(t)

	user := UserAuth{
		Username: "review-user",
		Password: "123456",
		UserInfo: &UserInfo{
			Nickname: "reviewer",
		},
	}
	db.Create(&user)

	article := Article{Title: "title", Content: "content"}
	db.Create(&article)

	root, err := AddComment(db, user.ID, TYPE_ARTICLE, article.ID, "visible", true)
	assert.NoError(t, err)
	_, err = AddReplyComment(db, user.ID, user.ID, root.ID, "visible reply", true)
	assert.NoError(t, err)
	_, err = AddReplyComment(db, user.ID, user.ID, root.ID, "pending reply", false)
	assert.NoError(t, err)
	_, err = AddComment(db, user.ID, TYPE_ARTICLE, article.ID, "pending root", false)
	assert.NoError(t, err)

	comments, total, err := GetBlogCommentList(db, 1, 10, TYPE_ARTICLE, article.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, comments, 1)
	assert.Equal(t, "visible", comments[0].Content)
	assert.Equal(t, 1, comments[0].ReplyCount)
	assert.Len(t, comments[0].ReplyList, 1)
	assert.Equal(t, "visible reply", comments[0].ReplyList[0].Content)

	replies, err := GetCommentByid(db, 1, 10, root.ID)
	assert.NoError(t, err)
	assert.Len(t, replies, 1)
	assert.Equal(t, "visible reply", replies[0].Content)
}
