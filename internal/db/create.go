package db

import (
	"encoding/json"
	"io"
	"postit/internal/postit_messages"
	"postit/model"
)

func CreateUser(username string) (*model.User, error) {
	user := model.User{
		Username: username,
	}
	connection := GetDB()
	createdUser := connection.Create(&user)

	return &user, createdUser.Error
}

func CreatePost(body io.ReadCloser, username string) (*model.Post, error) {
	post := model.Post{}
	err := json.NewDecoder(body).Decode(&post)
	if err != nil {
		return &post, postit_messages.ErrFailedDecoding
	}

	post.Username = username
	post.Active = true

	if !post.Validate() {
		return &post, postit_messages.ErrInvalid
	}

	connection := GetDB()
	createdPost := connection.Create(&post)

	return &post, createdPost.Error
}

func CreateComment(body io.ReadCloser, username string, postID int) (*model.Comment, error) {
	comment := model.Comment{}
	err := json.NewDecoder(body).Decode(&comment)
	if err != nil {
		return &comment, postit_messages.ErrFailedDecoding
	}

	comment.Username = username
	comment.PostID = postID
	comment.Active = true

	connection := GetDB()
	createdComment := connection.Create(&comment)

	return &comment, createdComment.Error
}

func CreateLike(username string, postID int) (*model.Like, error) {
	like := model.Like{
		Username: username,
		PostID:   postID,
		Active:   true,
	}

	connection := GetDB()
	createdLike := connection.Create(&like)

	return &like, createdLike.Error
}
