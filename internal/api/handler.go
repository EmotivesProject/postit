package api

import (
	"fmt"
	"net/http"
	"postit/internal/db"
	"postit/internal/event"
	"postit/internal/postit_messages"
	"postit/model"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/response"

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
	postitDatabase := db.GetDatabase()
	username := r.Context().Value(userID).(string)
	user, err := db.FindUser(username, postitDatabase)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	result, post, err := db.CreatePost(
		r.Body,
		user.ID,
		postitDatabase,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	logger.Infof("Successfully created post %s", username)
	event.SendPostEvent(username, model.StatusCreated, post)

	response.ResultResponseJSON(w, http.StatusCreated, result)
}

func createComment(w http.ResponseWriter, r *http.Request) {
	postitDatabase := db.GetDatabase()
	username := r.Context().Value(userID).(string)
	user, err := db.FindUser(username, postitDatabase)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	postID := chi.URLParam(r, postParam)
	post, err := db.FindPostById(postID, postitDatabase)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	comment, err := db.CreateComment(
		r.Body,
		user.ID,
		postitDatabase,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	post.UserComments = append(post.UserComments, comment.ID)
	err = db.UpdatePost(
		&post,
		postitDatabase,
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
	postitDatabase := db.GetDatabase()
	username := r.Context().Value(userID).(string)
	user, err := db.FindUser(username, postitDatabase)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	postID := chi.URLParam(r, postParam)
	post, err := db.FindPostById(postID, postitDatabase)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	// Check if the user has already liked the post
	// Get all the likes on the post and see if the user is in there.
	// Set it active to true if it's false or error if already liked
	if post.UserLikes != nil {
		likesOnPost, err := db.FindByLikeIDS(post.UserLikes, postitDatabase)
		if err != nil {
			logger.Error(err)
			response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
			return
		}

		for _, rangedLike := range likesOnPost {
			if rangedLike.User == user.ID {
				if rangedLike.Active {
					logger.Infof("User %s has already liked %s", username, postID)
					response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: postit_messages.ErrAlreadyLiked.Error()})
					return
				}
				rangedLike.Active = true
				err = db.UpdateLike(&rangedLike, postitDatabase)
				if err != nil {
					logger.Error(err)
					response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
					return
				}

				response.ResultResponseJSON(w, http.StatusCreated, rangedLike)
				return
			}
		}
	}

	// Continue with creating a new like
	like, err := db.CreateLike(
		user.ID,
		postitDatabase,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	post.UserLikes = append(post.UserLikes, like.ID)
	err = db.UpdatePost(
		&post,
		postitDatabase,
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
	postitDatabase := db.GetDatabase()
	username := r.Context().Value(userID)
	usernameString := fmt.Sprintf("%v", username)
	user, err := db.FindUser(usernameString, postitDatabase)

	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	response.ResultResponseJSON(w, http.StatusOK, user)
}

func fetchPost(w http.ResponseWriter, r *http.Request) {
	postitDatabase := db.GetDatabase()
	begin := findBegin(r)

	postPipelineMongo := db.FindPostPipeline(begin)

	rawData, err := db.GetRawResponseFromAggregate(
		db.PostCollection,
		postPipelineMongo,
		postitDatabase,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	response.ResultResponseJSON(w, http.StatusOK, rawData)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	postitDatabase := db.GetDatabase()
	username := r.Context().Value(userID)
	usernameString := fmt.Sprintf("%v", username)
	postID := chi.URLParam(r, postParam)
	post, err := db.FindPostById(postID, postitDatabase)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	post.Active = false
	err = db.UpdatePost(&post, postitDatabase)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}
	logger.Infof("Successfully deleted post %s", postID)
	event.SendPostEvent(usernameString, model.StatusDeleted, &post)
	response.ResultResponseJSON(w, http.StatusOK, post)
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userID)
	usernameString := fmt.Sprintf("%v", username)
	postitDatabase := db.GetDatabase()
	commentID := chi.URLParam(r, "comment_id")
	comment, err := db.FindCommentById(commentID, postitDatabase)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	comment.Active = false
	err = db.UpdateComment(&comment, postitDatabase)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}
	logger.Infof("Successfully deleted comment %s", commentID)
	event.SendCommentEvent(usernameString, model.StatusDeleted, &comment)
	response.ResultResponseJSON(w, http.StatusOK, comment)
}

func deleteLike(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userID)
	usernameString := fmt.Sprintf("%v", username)
	postitDatabase := db.GetDatabase()
	likeID := chi.URLParam(r, "like_id")
	like, err := db.FindLikeById(likeID, postitDatabase)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}

	like.Active = false
	err = db.UpdateLike(&like, postitDatabase)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, http.StatusBadRequest, response.Message{Message: err.Error()})
		return
	}
	logger.Infof("Successfully deleted like %s", likeID)
	event.SendLikeEvent(usernameString, model.StatusDeleted, &like)
	response.ResultResponseJSON(w, http.StatusOK, like)
}
