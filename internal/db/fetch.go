package db

import (
	"postit/model"

	"gorm.io/gorm"
)

func FindByUsername(username string) (model.User, error) {
	user := &model.User{}
	database := GetDB()

	if err := database.Where("username = ?", username).First(user).Error; err != nil {
		return *user, err
	}

	return *user, nil
}

func CheckUsername(username string) error {
	_, err := FindByUsername(username)
	if err == gorm.ErrRecordNotFound {
		_, err = CreateUser(username)
	}
	return err
}

func FindPostById(postID int) (model.Post, error) {
	post := &model.Post{}
	database := GetDB()

	if err := database.Where("id = ?", postID).First(post).Error; err != nil {
		return *post, err
	}

	return *post, nil
}

func FindCommentById(commentID int) (model.Comment, error) {
	comment := &model.Comment{}
	database := GetDB()

	if err := database.Where("id = ?", commentID).First(comment).Error; err != nil {
		return *comment, err
	}

	return *comment, nil
}

func FindLikeById(likeID int) (model.Like, error) {
	like := &model.Like{}
	database := GetDB()

	if err := database.Where("id = ?", likeID).First(like).Error; err != nil {
		return *like, err
	}

	return *like, nil
}

func FindCommentsForPost(postID int) ([]model.Comment, error) {
	var comments []model.Comment
	database := GetDB()

	if err := database.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		return comments, err
	}

	return comments, nil
}

func FindLikesForPost(postID int) ([]model.Like, error) {
	var likes []model.Like
	database := GetDB()

	if err := database.Where("post_id = ? AND active = true", postID).Find(&likes).Error; err != nil {
		return likes, err
	}

	return likes, nil
}
