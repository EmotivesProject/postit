package db

import (
	"postit/model"
)

func UpdatePost(post *model.Post) error {
	database := GetDB()
	result := database.Save(post)
	return result.Error
}

func UpdateComment(comment *model.Comment) error {
	database := GetDB()
	result := database.Save(comment)
	return result.Error
}

func UpdateLike(like *model.Like) error {
	database := GetDB()
	result := database.Save(like)
	return result.Error
}
