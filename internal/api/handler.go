package api

import (
	"net/http"
	"postit/internal/db"
	"postit/internal/event"
	"postit/internal/postit_messages"
	"postit/model"
	"strconv"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/response"
	"github.com/TomBowyerResearchProject/common/verification"

	"github.com/go-chi/chi"
)

var (
	postParam         = "post_id"
	limit       int64 = 5
	msgHealthOK       = "Health ok"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	response.MessageResponseJSON(w, http.StatusOK, response.Message{Message: msgHealthOK})
}

func createPost(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(verification.UserID).(string)

	err := db.CheckUsername(username)
	if err != nil {
		logger.Error(postit_messages.ErrInvalidUsername)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: postit_messages.ErrInvalidUsername.Error()})
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
	event.SendPostEvent(username, model.StatusCreated, post)

	response.ResultResponseJSON(w, http.StatusCreated, post)
}

func createComment(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(verification.UserID).(string)

	err := db.CheckUsername(username)
	if err != nil {
		logger.Error(postit_messages.ErrInvalidUsername)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: postit_messages.ErrInvalidUsername.Error()})
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

	logger.Infof("Created comment for %s", username)
	event.SendCommentEvent(username, model.StatusCreated, comment)

	response.ResultResponseJSON(w, http.StatusCreated, comment)
}

func createLike(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(verification.UserID).(string)

	err := db.CheckUsername(username)
	if err != nil {
		logger.Error(postit_messages.ErrInvalidUsername)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: postit_messages.ErrInvalidUsername.Error()})
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
	event.SendLikeEvent(username, model.StatusCreated, like)
	response.ResultResponseJSON(w, http.StatusCreated, like)
}

func fetchUserFromAuth(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(verification.UserID).(string)

	err := db.CheckUsername(username)
	if err != nil {
		logger.Error(postit_messages.ErrInvalidUsername)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: postit_messages.ErrInvalidUsername.Error()})
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
	username := r.Context().Value(verification.UserID).(string)
	postString := chi.URLParam(r, postParam)
	postID, err := strconv.Atoi(postString)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	post, err := db.FindPostById(postID)
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
	event.SendPostEvent(username, model.StatusDeleted, &post)
	response.ResultResponseJSON(w, http.StatusOK, post)
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(verification.UserID).(string)
	commentString := chi.URLParam(r, "comment_id")
	commentID, err := strconv.Atoi(commentString)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	comment, err := db.FindCommentById(commentID)
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
	event.SendCommentEvent(username, model.StatusDeleted, &comment)
	response.ResultResponseJSON(w, http.StatusOK, comment)
}

func deleteLike(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(verification.UserID).(string)
	likeString := chi.URLParam(r, "like_id")
	likeID, err := strconv.Atoi(likeString)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	like, err := db.FindLikeById(likeID)
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
	event.SendLikeEvent(username, model.StatusDeleted, &like)
	response.ResultResponseJSON(w, http.StatusOK, like)
}

// Next two functions don't just return individual objects, but all the references to it e.g. All the comments and
func fetchPosts(w http.ResponseWriter, r *http.Request) {
}

func fetchIndividualPost(w http.ResponseWriter, r *http.Request) {
	postString := chi.URLParam(r, postParam)
	postID, err := strconv.Atoi(postString)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	post, err := db.FindPostById(postID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}
	logger.Info("PAST POST")

	comments, err := db.FindCommentsForPost(postID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}
	logger.Info("PAST Comment")

	likes, err := db.FindLikesForPost(postID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}
	logger.Info("PAST Likes")

	postInfo := model.PostInformation{
		Post:     post,
		Comments: comments,
		Likes:    likes,
	}
	response.ResultResponseJSON(w, http.StatusOK, postInfo)
}
