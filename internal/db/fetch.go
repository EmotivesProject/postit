package db

import (
	"context"
	"postit/model"

	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
)

const (
	PostLimit    = 5
	ExploreBound = 2
)

func CheckUsername(ctx context.Context, user model.User) error {
	_, err := FindByUsername(ctx, user.Username)
	if err != nil {
		_, err = CreateUser(ctx, user)
	}

	return err
}

func FindByUsername(ctx context.Context, username string) (model.User, error) {
	user := &model.User{}
	connection := commonPostgres.GetDatabase()

	err := connection.QueryRow(
		ctx,
		"SELECT * FROM users WHERE username = $1",
		username,
	).Scan(
		&user.ID, &user.Username, &user.UserGroup,
	)

	return *user, err
}

func FindPostByIDForUser(ctx context.Context, postID int, user model.User) (model.Post, error) {
	post := &model.Post{}
	connection := commonPostgres.GetDatabase()

	err := connection.QueryRow(
		ctx,
		`SELECT id, username, content, created_at, updated_at, active FROM posts WHERE id = $1 and username in
		(select username from users where user_group = $2)`,
		postID,
		user.UserGroup,
	).Scan(
		&post.ID, &post.Username, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Active,
	)

	return *post, err
}

func FindPostByID(ctx context.Context, postID int) (model.Post, error) {
	post := &model.Post{}
	connection := commonPostgres.GetDatabase()

	err := connection.QueryRow(
		ctx,
		`SELECT id, username, content, created_at, updated_at, active FROM posts WHERE id = $1`,
		postID,
	).Scan(
		&post.ID, &post.Username, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Active,
	)

	return *post, err
}

func FindPosts(ctx context.Context, offset int, user model.User) ([]model.Post, error) {
	posts := make([]model.Post, 0)

	connection := commonPostgres.GetDatabase()

	rows, err := connection.Query(
		ctx,
		`SELECT * FROM posts where username in
		(select username from users where user_group = $1)
		ORDER BY created_at desc LIMIT 5 OFFSET $2`,
		user.UserGroup,
		offset,
	)
	if err != nil {
		return posts, err
	}

	for rows.Next() {
		var post model.Post

		err := rows.Scan(
			&post.ID,
			&post.Username,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Active,
		)
		if err != nil {
			continue
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func FindPostsBasedOnLatAndLng(
	ctx context.Context, lat, lng float64, offset int, user model.User,
) ([]model.Post, error) {
	posts := make([]model.Post, 0)

	connection := commonPostgres.GetDatabase()

	minLng := lng - ExploreBound
	maxLng := lng + ExploreBound

	minLat := lat - ExploreBound
	maxLat := lat + ExploreBound

	rows, err := connection.Query(
		ctx,
		`SELECT * FROM posts
		WHERE (content->'latitude')::float < $1 AND (content->'latitude')::float > $2
		AND (content->'longitude')::float < $3 AND (content->'longitude')::float > $4
		AND username in (select username from users where user_group = $5)
		ORDER BY created_at desc LIMIT 5 OFFSET $6`,
		maxLat,
		minLat,
		maxLng,
		minLng,
		user.UserGroup,
		offset,
	)
	if err != nil {
		return posts, err
	}

	for rows.Next() {
		var post model.Post

		err := rows.Scan(
			&post.ID,
			&post.Username,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Active,
		)
		if err != nil {
			continue
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func FindLikeByID(ctx context.Context, likeID int) (model.Like, error) {
	like := &model.Like{}
	connection := commonPostgres.GetDatabase()

	err := connection.QueryRow(
		ctx,
		"SELECT id, post_id, username, created_at, updated_at, active FROM likes WHERE id = $1",
		likeID,
	).Scan(
		&like.ID, &like.PostID, &like.Username, &like.CreatedAt, &like.UpdatedAt, &like.Active,
	)

	return *like, err
}

func FindCommentsForPost(ctx context.Context, postID int, full bool) ([]model.Comment, error) {
	comments := make([]model.Comment, 0)

	connection := commonPostgres.GetDatabase()

	var query string

	if full {
		query = "SELECT * FROM comments WHERE post_id = $1 AND active = true ORDER BY created_at desc"
	} else {
		query = "SELECT * FROM comments WHERE post_id = $1 AND active = true ORDER BY created_at desc LIMIT 3"
	}

	rows, err := connection.Query(
		ctx,
		query,
		postID,
	)
	if err != nil {
		return comments, err
	}

	for rows.Next() {
		var comment model.Comment

		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Username,
			&comment.Message,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Active,
		)
		if err != nil {
			continue
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func FindLikesForPost(ctx context.Context, postID int) ([]model.Like, error) {
	likes := make([]model.Like, 0)

	connection := commonPostgres.GetDatabase()

	rows, err := connection.Query(
		ctx,
		"SELECT * FROM likes WHERE post_id = $1 AND active = true  ORDER BY created_at desc",
		postID,
	)
	if err != nil {
		return likes, err
	}

	for rows.Next() {
		var like model.Like

		err := rows.Scan(
			&like.ID,
			&like.PostID,
			&like.Username,
			&like.CreatedAt,
			&like.UpdatedAt,
			&like.Active,
		)
		if err != nil {
			continue
		}

		likes = append(likes, like)
	}

	return likes, nil
}
