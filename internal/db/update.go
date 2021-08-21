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

func UpdateLikeByUsernameAndPost(ctx context.Context, like model.Like) (model.Like, error) {
	connection := commonPostgres.GetDatabase()

	err := connection.QueryRow(
		ctx,
		"UPDATE likes SET active = $1 WHERE post_id = $2 and username = $3 RETURNING id",
		like.Active,
		like.PostID,
		like.Username,
	).Scan(
		&like.ID,
	)

	return like, err
}
