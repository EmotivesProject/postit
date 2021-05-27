package api

import (
	"context"
	"net/http"
	"postit/internal/db"
	"postit/internal/send"
	"postit/messages"
	"postit/model"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/response"
	"github.com/TomBowyerResearchProject/common/verification"
)

var (
	postParam    = "post_id"
	likeParam    = "like_id"
	commentParam = "comment_id"
)

func fetchUserFromAuth(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.ErrInvalidCheck.Error()})

		return
	}

	err := db.CheckUsername(r.Context(), username)
	if err != nil {
		logger.Error(messages.ErrInvalidUsername)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.ErrInvalidUsername.Error()})

		return
	}

	user, err := db.FindByUsername(r.Context(), username)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	response.ResultResponseJSON(w, http.StatusOK, user)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.ErrInvalidCheck.Error()})

		return
	}

	err := db.CheckUsername(r.Context(), username)
	if err != nil {
		logger.Error(messages.ErrInvalidUsername)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.ErrInvalidUsername.Error()})

		return
	}

	post, err := db.CreatePost(
		r.Context(),
		r.Body,
		username,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Successfully created post %s", username)
	response.ResultResponseJSON(w, http.StatusCreated, post)
}

func createComment(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.ErrInvalidCheck.Error()})

		return
	}

	err := db.CheckUsername(r.Context(), username)
	if err != nil {
		logger.Error(messages.ErrInvalidUsername)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.ErrInvalidUsername.Error()})

		return
	}

	postID, err := extractID(r, postParam)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	comment, err := db.CreateComment(
		r.Context(),
		r.Body,
		username,
		postID,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	post, err := db.FindPostByID(r.Context(), postID)
	if err != nil {
		logger.Error(err)
	} else if post.Username != username {
		send.SendComment(post.Username, username, post.ID)
	}

	logger.Infof("Created comment for %s", username)
	response.ResultResponseJSON(w, http.StatusCreated, comment)
}

func createLike(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.ErrInvalidCheck.Error()})

		return
	}

	err := db.CheckUsername(r.Context(), username)
	if err != nil {
		logger.Error(messages.ErrInvalidUsername)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.ErrInvalidUsername.Error()})

		return
	}

	postID, err := extractID(r, postParam)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	like, err := db.CreateLike(
		r.Context(),
		username,
		postID,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	post, err := db.FindPostByID(r.Context(), postID)
	if err != nil {
		logger.Error(err)
	} else if post.Username != username {
		send.SendLike(post.Username, username, post.ID)
	}

	logger.Infof("Created like for %s", username)
	response.ResultResponseJSON(w, http.StatusCreated, like)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	postID, err := extractID(r, postParam)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	post, err := updatePost(r.Context(), postID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Successfully deleted post %d", postID)
	response.ResultResponseJSON(w, http.StatusOK, post)
}

func updatePost(ctx context.Context, postID int) (model.Post, error) {
	post, err := db.FindPostByID(ctx, postID)
	if err != nil {
		return post, err
	}

	post.Active = false

	err = db.UpdatePost(ctx, &post)

	return post, err
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := extractID(r, commentParam)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	comment, err := updateComment(r.Context(), commentID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Successfully deleted comment %d", commentID)
	response.ResultResponseJSON(w, http.StatusOK, comment)
}

func updateComment(ctx context.Context, commentID int) (model.Comment, error) {
	comment, err := db.FindCommentByID(ctx, commentID)
	if err != nil {
		return comment, err
	}

	comment.Active = false

	err = db.UpdateComment(ctx, &comment)

	return comment, err
}

func deleteLike(w http.ResponseWriter, r *http.Request) {
	likeID, err := extractID(r, likeParam)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	like, err := updateLike(r.Context(), likeID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Successfully deleted like %d", likeID)
	response.ResultResponseJSON(w, http.StatusOK, like)
}

func updateLike(ctx context.Context, likeID int) (model.Like, error) {
	like, err := db.FindLikeByID(ctx, likeID)
	if err != nil {
		return like, err
	}

	like.Active = false

	err = db.UpdateLike(ctx, &like)

	return like, err
}

func fetchPosts(w http.ResponseWriter, r *http.Request) {
	page := findBegin(r)

	postInformations := make([]model.PostInformation, 0)

	posts, err := db.FindPosts(r.Context(), page)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	for _, post := range posts {
		comments, err := db.FindCommentsForPost(r.Context(), post.ID)
		if err != nil {
			logger.Error(err)
			response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

			return
		}

		likes, err := db.FindLikesForPost(r.Context(), post.ID)
		if err != nil {
			logger.Error(err)
			response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

			return
		}

		postInformation := model.PostInformation{
			Post:     post,
			Comments: comments,
			Likes:    likes,
		}

		postInformations = append(postInformations, postInformation)
	}

	response.ResultResponseJSON(w, http.StatusOK, postInformations)
}

func fetchIndividualPost(w http.ResponseWriter, r *http.Request) {
	postID, err := extractID(r, postParam)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	post, err := db.FindPostByID(r.Context(), postID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	comments, err := db.FindCommentsForPost(r.Context(), postID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	likes, err := db.FindLikesForPost(r.Context(), postID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	postInfo := model.PostInformation{
		Post:     post,
		Comments: comments,
		Likes:    likes,
	}

	response.ResultResponseJSON(w, http.StatusOK, postInfo)
}
