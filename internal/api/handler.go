package api

import (
	"fmt"
	"net/http"
	"postit/internal/db"
	"postit/messages"
	"postit/model"
	"strconv"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/notification"
	"github.com/TomBowyerResearchProject/common/response"
	"github.com/TomBowyerResearchProject/common/verification"

	"github.com/go-chi/chi"
)

var postParam = "post_id"

func createPost(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.ErrInvalidCheck.Error()})

		return
	}

	err := db.CheckUsername(username)
	if err != nil {
		logger.Error(messages.ErrInvalidUsername)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.ErrInvalidUsername.Error()})

		return
	}

	post, err := db.CreatePost(
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

	err := db.CheckUsername(username)
	if err != nil {
		logger.Error(messages.ErrInvalidUsername)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.ErrInvalidUsername.Error()})

		return
	}

	postString := chi.URLParam(r, postParam)

	postID, err := strconv.Atoi(postString)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	comment, err := db.CreateComment(
		r.Body,
		username,
		postID,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	post, err := db.FindPostByID(postID)
	if err != nil {
		logger.Error(err)
	} else if post.Username != username {
		notif := notification.Notification{
			Username: post.Username,
			Title:    "New Comment on your post",
			Message:  fmt.Sprintf("%s commented on your post", username),
			Link:     "WHATEVER",
		}

		err := notification.SendEvent("notif.localhost/notification", notif)
		if err != nil {
			logger.Error(err)
		}
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

	err := db.CheckUsername(username)
	if err != nil {
		logger.Error(messages.ErrInvalidUsername)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.ErrInvalidUsername.Error()})

		return
	}

	postString := chi.URLParam(r, postParam)

	postID, err := strconv.Atoi(postString)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	like, err := db.CreateLike(
		username,
		postID,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Created like for %s", username)
	response.ResultResponseJSON(w, http.StatusCreated, like)
}

func fetchUserFromAuth(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.ErrInvalidCheck.Error()})

		return
	}

	err := db.CheckUsername(username)
	if err != nil {
		logger.Error(messages.ErrInvalidUsername)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: messages.ErrInvalidUsername.Error()})

		return
	}

	user, err := db.FindByUsername(username)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	response.ResultResponseJSON(w, http.StatusOK, user)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	postString := chi.URLParam(r, postParam)

	postID, err := strconv.Atoi(postString)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	post, err := db.FindPostByID(postID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	post.Active = false

	err = db.UpdatePost(&post)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Successfully deleted post %d", postID)
	response.ResultResponseJSON(w, http.StatusOK, post)
}

// nolint:dupl
func deleteComment(w http.ResponseWriter, r *http.Request) {
	commentString := chi.URLParam(r, "comment_id")

	commentID, err := strconv.Atoi(commentString)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	comment, err := db.FindCommentByID(commentID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	comment.Active = false

	err = db.UpdateComment(&comment)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Successfully deleted comment %d", commentID)
	response.ResultResponseJSON(w, http.StatusOK, comment)
}

// nolint:dupl
func deleteLike(w http.ResponseWriter, r *http.Request) {
	likeString := chi.URLParam(r, "like_id")

	likeID, err := strconv.Atoi(likeString)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	like, err := db.FindLikeByID(likeID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	like.Active = false

	err = db.UpdateLike(&like)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Successfully deleted like %d", likeID)
	response.ResultResponseJSON(w, http.StatusOK, like)
}

// Next two functions don't just return individual objects, but all the references to it
// e.g. All the comments and likes.
func fetchPosts(w http.ResponseWriter, r *http.Request) {
	page := findBegin(r)

	var postInformations []model.PostInformation // nolint: prealloc

	posts, err := db.FindPosts(page)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	for _, post := range posts {
		comments, err := db.FindCommentsForPost(post.ID)
		if err != nil {
			logger.Error(err)
			response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

			return
		}

		likes, err := db.FindLikesForPost(post.ID)
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
	postString := chi.URLParam(r, postParam)

	postID, err := strconv.Atoi(postString)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	post, err := db.FindPostByID(postID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	comments, err := db.FindCommentsForPost(postID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	likes, err := db.FindLikesForPost(postID)
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
