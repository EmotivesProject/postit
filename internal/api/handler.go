package api

import (
	"fmt"
	"net/http"
	"postit/internal/db"
	"postit/internal/event"
	"postit/internal/logger"
	"postit/internal/postit_messages"
	"postit/model"

	"github.com/go-chi/chi"
)

var (
	postParam         = "post_id"
	limit       int64 = 5
	msgHealthOK       = "Health ok"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	messageResponseJSON(w, http.StatusOK, model.Message{Message: msgHealthOK})
}

func createPost(w http.ResponseWriter, r *http.Request) {
	postitDatabase := db.GetDatabase()
	username := r.Context().Value(userID).(string)
	user, err := db.FindUser(username, postitDatabase)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	result, post, err := db.CreatePost(
		r.Body,
		user.ID,
		postitDatabase,
	)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	logger.Infof("Successfully created post with for %s", username)
	event.SendPostEvent(username, model.StatusCreated, post)

	resultResponseJSON(w, http.StatusCreated, result)
}

func createComment(w http.ResponseWriter, r *http.Request) {
	postitDatabase := db.GetDatabase()
	username := r.Context().Value(userID).(string)
	user, err := db.FindUser(username, postitDatabase)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	postID := chi.URLParam(r, postParam)
	post, err := db.FindPostById(postID, postitDatabase)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	comment, err := db.CreateComment(
		r.Body,
		user.ID,
		postitDatabase,
	)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	post.UserComments = append(post.UserComments, comment.ID)
	err = db.UpdatePost(
		&post,
		postitDatabase,
	)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	logger.Infof("Created comment for %s", username)
	event.SendCommentEvent(username, model.StatusCreated, comment)

	resultResponseJSON(w, http.StatusCreated, comment)
}

func createLike(w http.ResponseWriter, r *http.Request) {
	postitDatabase := db.GetDatabase()
	username := r.Context().Value(userID).(string)
	user, err := db.FindUser(username, postitDatabase)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	postID := chi.URLParam(r, postParam)
	post, err := db.FindPostById(postID, postitDatabase)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	// Check if the user has already liked the post
	// Get all the likes on the post and see if the user is in there.
	// Set it active to true if it's false or error if already liked
	if post.UserLikes != nil {
		likesOnPost, err := db.FindByLikeIDS(post.UserLikes, postitDatabase)
		if err != nil {
			logger.Error(err)
			messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
			return
		}

		for _, rangedLike := range likesOnPost {
			if rangedLike.User == user.ID {
				if rangedLike.Active {
					logger.Infof("User %s has already liked %s", username, postID)
					messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: postit_messages.ErrAlreadyLiked.Error()})
					return
				}
				rangedLike.Active = true
				err = db.UpdateLike(&rangedLike, postitDatabase)
				if err != nil {
					logger.Error(err)
					messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
					return
				}

				resultResponseJSON(w, http.StatusCreated, rangedLike)
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
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	post.UserLikes = append(post.UserLikes, like.ID)
	err = db.UpdatePost(
		&post,
		postitDatabase,
	)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	logger.Infof("Created like for %s", username)
	event.SendLikeEvent(username, model.StatusCreated, like)
	resultResponseJSON(w, http.StatusCreated, like)
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
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	resultResponseJSON(w, http.StatusOK, rawData)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	postitDatabase := db.GetDatabase()
	username := r.Context().Value(userID)
	usernameString := fmt.Sprintf("%v", username)
	postID := chi.URLParam(r, postParam)
	post, err := db.FindPostById(postID, postitDatabase)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	post.Active = false
	err = db.UpdatePost(&post, postitDatabase)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	logger.Infof("Successfully deleted post %s", postID)
	event.SendPostEvent(usernameString, model.StatusDeleted, &post)
	resultResponseJSON(w, http.StatusOK, post)
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userID)
	usernameString := fmt.Sprintf("%v", username)
	postitDatabase := db.GetDatabase()
	commentID := chi.URLParam(r, "comment_id")
	comment, err := db.FindCommentById(commentID, postitDatabase)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	comment.Active = false
	err = db.UpdateComment(&comment, postitDatabase)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	logger.Infof("Successfully deleted comment %s", commentID)
	event.SendCommentEvent(usernameString, model.StatusDeleted, &comment)
	resultResponseJSON(w, http.StatusOK, comment)
}

func deleteLike(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userID)
	usernameString := fmt.Sprintf("%v", username)
	postitDatabase := db.GetDatabase()
	likeID := chi.URLParam(r, "like_id")
	like, err := db.FindLikeById(likeID, postitDatabase)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	like.Active = false
	err = db.UpdateLike(&like, postitDatabase)
	if err != nil {
		logger.Error(err)
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	logger.Infof("Successfully deleted like %s", likeID)
	event.SendLikeEvent(usernameString, model.StatusDeleted, &like)
	resultResponseJSON(w, http.StatusOK, like)
}
