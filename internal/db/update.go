package db

import (
	"context"
	"postit/model"

	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
)

func UpdatePost(post *model.Post) error {
	connection := commonPostgres.GetDatabase()

	_, err := connection.Exec(
		context.TODO(),
		"UPDATE post SET active = $1 WHERE post_id = $2",
		post.Active,
		post.ID,
	)

	return err
}

func UpdateComment(comment *model.Comment) error {
	connection := commonPostgres.GetDatabase()

	_, err := connection.Exec(
		context.TODO(),
		"UPDATE comment SET active = $1 WHERE comment_id = $2",
		comment.Active,
		comment.ID,
	)

	return err
}

func UpdateLike(like *model.Like) error {
	connection := commonPostgres.GetDatabase()

	_, err := connection.Exec(
		context.TODO(),
		"UPDATE like SET active = $1 WHERE like_id = $2",
		like.Active,
		like.ID,
	)

	return err
}
