package db

import (
	"context"
	"postit/model"

	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
)

func UpdatePost(ctx context.Context, post *model.Post) error {
	connection := commonPostgres.GetDatabase()

	_, err := connection.Exec(
		ctx,
		"UPDATE posts SET active = $1 WHERE id = $2",
		post.Active,
		post.ID,
	)

	return err
}

func UpdateComment(ctx context.Context, comment *model.Comment) error {
	connection := commonPostgres.GetDatabase()

	_, err := connection.Exec(
		ctx,
		"UPDATE comments SET active = $1 WHERE id = $2",
		comment.Active,
		comment.ID,
	)

	return err
}

func UpdateLike(ctx context.Context, like *model.Like) error {
	connection := commonPostgres.GetDatabase()

	_, err := connection.Exec(
		ctx,
		"UPDATE likes SET active = $1 WHERE id = $2",
		like.Active,
		like.ID,
	)

	return err
}
