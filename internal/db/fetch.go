package db

import (
	"context"
	"postit/model"

	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
	"gorm.io/gorm"
)

const (
	PostLimit = 5
)

func CheckUsername(username string) error {
	_, err := FindByUsername(username)
	if err == gorm.ErrRecordNotFound {
		_, err = CreateUser(username)
	}
	return err
}

func FindByUsername(username string) (model.User, error) {
	user := &model.User{}
	connection := commonPostgres.GetDatabase()

	err := connection.QueryRow(
		context.Background(),
		"SELECT id, username FROM users WHERE username = $1",
		username,
	).Scan(
		&user.ID, &user.Username,
	)

	return *user, err
}

func FindPostById(postID int) (model.Post, error) {
	post := &model.Post{}
	connection := commonPostgres.GetDatabase()

	err := connection.QueryRow(
		context.Background(),
		"SELECT id, username, content, created_at, updated_at, active FROM posts WHERE post_id = $1",
		postID,
	).Scan(
		&post.ID, &post.Username, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.Active,
	)

	return *post, err
}

func FindPosts(offset int) ([]model.Post, error) {
	var posts []model.Post
	connection := commonPostgres.GetDatabase()

	rows, err := connection.Query(
		context.Background(),
		"SELECT * FROM posts ORDER BY created_at desc LIMIT 5 OFFSET $1",
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

func FindCommentById(commentID int) (model.Comment, error) {
	comment := &model.Comment{}
	connection := commonPostgres.GetDatabase()

	err := connection.QueryRow(
		context.Background(),
		"SELECT id, post_id, username, message, created_at, updated_at, active FROM posts WHERE post_id = $1",
		commentID,
	).Scan(
		&comment.ID, &comment.PostID, &comment.Username, &comment.Message, &comment.CreatedAt, &comment.UpdatedAt, &comment.Active,
	)

	return *comment, err
}

func FindLikeById(likeID int) (model.Like, error) {
	like := &model.Like{}
	connection := commonPostgres.GetDatabase()

	err := connection.QueryRow(
		context.Background(),
		"SELECT id, post_id, username, created_at, updated_at, active FROM posts WHERE post_id = $1",
		likeID,
	).Scan(
		&like.ID, &like.PostID, &like.Username, &like.CreatedAt, &like.UpdatedAt, &like.Active,
	)

	return *like, err
}

func FindCommentsForPost(postID int) ([]model.Comment, error) {
	var comments []model.Comment
	connection := commonPostgres.GetDatabase()

	rows, err := connection.Query(
		context.Background(),
		"SELECT * FROM comments WHERE post_id = $1 ORDER BY created_at desc",
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

func FindLikesForPost(postID int) ([]model.Like, error) {
	var likes []model.Like
	connection := commonPostgres.GetDatabase()

	rows, err := connection.Query(
		context.Background(),
		"SELECT * FROM likes WHERE post_id = $1 ORDER BY created_at desc",
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
