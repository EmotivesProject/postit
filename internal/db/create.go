package db

import (
	"context"
	"encoding/json"
	"io"
	"postit/messages"
	"postit/model"
	"time"

	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
)

func CreateUser(ctx context.Context, username string) (*model.User, error) {
	user := model.User{
		Username: username,
	}

	connection := commonPostgres.GetDatabase()
	_, err := connection.Exec(
		ctx,
		"INSERT INTO users(username) VALUES ($1)",
		user.Username,
	)

	return &user, err
}

func CreatePost(ctx context.Context, body io.ReadCloser, username string) (*model.Post, error) {
	post := model.Post{}
	jsonMap := make(map[string]interface{})

	err := json.NewDecoder(body).Decode(&jsonMap)
	if err != nil {
		return &post, err
	}

	post.Username = username
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	post.Active = true

	convertedContent, ok := jsonMap["content"].(map[string]interface{})
	if !ok {
		return &post, messages.ErrInvalid
	}

	post.Content = convertedContent

	connection := commonPostgres.GetDatabase()
	err = connection.QueryRow(
		ctx,
		"INSERT INTO posts(username,content,created_at,updated_at,active) VALUES ($1,$2,$3,$4,$5) RETURNING id",
		post.Username,
		post.Content,
		post.CreatedAt,
		post.UpdatedAt,
		post.Active,
	).Scan(
		&post.ID,
	)

	return &post, err
}

func CreateComment(ctx context.Context, body io.ReadCloser, username string, postID int) (*model.Comment, error) {
	comment := model.Comment{}

	err := json.NewDecoder(body).Decode(&comment)
	if err != nil {
		return &comment, messages.ErrFailedDecoding
	}

	if comment.Message == "" {
		return nil, messages.ErrNoMessage
	}

	comment.PostID = postID
	comment.Username = username
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	comment.Active = true

	connection := commonPostgres.GetDatabase()
	err = connection.QueryRow(
		ctx,
		"INSERT INTO comments(post_id,username,message,created_at,updated_at,active) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id",
		comment.PostID,
		comment.Username,
		comment.Message,
		comment.CreatedAt,
		comment.UpdatedAt,
		comment.Active,
	).Scan(
		&comment.ID,
	)

	return &comment, err
}

func CreateLike(ctx context.Context, username string, postID int) (*model.Like, error) {
	like := model.Like{
		Username:  username,
		PostID:    postID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Active:    true,
	}

	connection := commonPostgres.GetDatabase()
	err := connection.QueryRow(
		ctx,
		"INSERT INTO likes(post_id,username,created_at,updated_at,active) VALUES ($1,$2,$3,$4,$5) RETURNING id",
		postID,
		username,
		time.Now(),
		time.Now(),
		true,
	).Scan(
		&like.ID,
	)

	return &like, err
}
