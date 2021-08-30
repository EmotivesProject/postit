package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"postit/internal/db"
	"postit/internal/send"
	"postit/messages"
	"postit/model"
	"time"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/redis"
	"github.com/TomBowyerResearchProject/common/response"
	"github.com/TomBowyerResearchProject/common/verification"
)

var (
	postParam  = "post_id"
	likeParam  = "like_id"
	redisCache = time.Minute * 10
)

func fetchUserFromAuth(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusUnprocessableEntity,
			response.Message{Message: messages.ErrInvalidCheck.Error()},
		)

		return
	}

	err := db.CheckUsername(r.Context(), username)
	if err != nil {
		logger.Error(messages.ErrInvalidUsername)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusBadRequest,
			response.Message{Message: messages.ErrInvalidUsername.Error()},
		)

		return
	}

	user, err := db.FindByUsername(r.Context(), username)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	response.ResultResponseJSON(w, false, http.StatusOK, user)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusUnprocessableEntity,
			response.Message{Message: messages.ErrInvalidCheck.Error()},
		)

		return
	}

	err := db.CheckUsername(r.Context(), username)
	if err != nil {
		logger.Error(messages.ErrInvalidUsername)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusInternalServerError,
			response.Message{Message: messages.ErrInvalidUsername.Error()},
		)

		return
	}

	post, err := db.CreatePost(
		r.Context(),
		r.Body,
		username,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	postInformation, err := createPostInformation(r.Context(), *post, username)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Successfully created post %s", username)

	response.ResultResponseJSON(w, false, http.StatusCreated, postInformation)
}

func createComment(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusUnprocessableEntity,
			response.Message{Message: messages.ErrInvalidCheck.Error()},
		)

		return
	}

	err := db.CheckUsername(r.Context(), username)
	if err != nil {
		logger.Error(messages.ErrInvalidUsername)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusInternalServerError,
			response.Message{Message: messages.ErrInvalidUsername.Error()},
		)

		return
	}

	postID, err := extractID(r, postParam)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	_, err = db.CreateComment(
		r.Context(),
		r.Body,
		username,
		postID,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	postInformation, err := createPostInformationWithFetchPosts(r.Context(), postID, username)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	redisKey := fmt.Sprintf("PostInfo.%d", postID)

	err = redis.SetEx(r.Context(), redisKey, *postInformation, redisCache)
	if err != nil {
		logger.Error(err)
	}

	logger.Infof("Created comment for %s", username)
	response.ResultResponseJSON(w, false, http.StatusCreated, postInformation)
}

func createLike(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusUnprocessableEntity,
			response.Message{Message: messages.ErrInvalidCheck.Error()},
		)

		return
	}

	err := db.CheckUsername(r.Context(), username)
	if err != nil {
		logger.Error(messages.ErrInvalidUsername)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusBadRequest,
			response.Message{Message: messages.ErrInvalidUsername.Error()},
		)

		return
	}

	postID, err := extractID(r, postParam)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	like, err := db.CreateLike(
		r.Context(),
		username,
		postID,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	post, err := db.FindPostByID(r.Context(), postID)
	if err != nil {
		logger.Error(err)
	} else if post.Username != username {
		send.Like(post.Username, username, post.ID)
	}

	logger.Infof("Created like for %s", username)
	response.ResultResponseJSON(w, false, http.StatusCreated, like)
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

func deletePost(w http.ResponseWriter, r *http.Request) {
	postID, err := extractID(r, postParam)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	post, err := updatePost(r.Context(), postID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Successfully deleted post %d", postID)
	response.ResultResponseJSON(w, false, http.StatusOK, post)
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

func deleteLike(w http.ResponseWriter, r *http.Request) {
	likeID, err := extractID(r, likeParam)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	like, err := updateLike(r.Context(), likeID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Successfully deleted like %d", likeID)
	response.ResultResponseJSON(w, false, http.StatusOK, like)
}

func fetchExplorePosts(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusUnprocessableEntity,
			response.Message{Message: messages.ErrInvalidCheck.Error()},
		)

		return
	}

	lat := findPosition(r, "lat")
	lng := findPosition(r, "lng")

	if lat == 0 || lng == 0 {
		response.MessageResponseJSON(
			w,
			false,
			http.StatusUnprocessableEntity,
			response.Message{Message: messages.ErrInvalidCheck.Error()},
		)
	}

	page := findBegin(r)

	postInformations := make([]model.PostInformation, 0)

	posts, err := db.FindPostsBasedOnLatAndLng(r.Context(), lat, lng, page)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	for _, post := range posts {
		redisKey := fmt.Sprintf("PostInfo.%d", post.ID)

		result, err := redis.Get(r.Context(), redisKey)
		if err == nil {
			resultModel := model.PostInformation{}
			err = json.Unmarshal([]byte(result), &resultModel)

			if err == nil {
				postInformations = append(postInformations, resultModel)

				continue
			}
		}

		postInformation, err := createPostInformation(r.Context(), post, username)
		if err != nil {
			logger.Error(err)
			response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

			return
		}

		err = redis.SetEx(r.Context(), redisKey, *postInformation, redisCache)
		if err != nil {
			logger.Error(err)
		}

		postInformations = append(postInformations, *postInformation)
	}

	response.ResultResponseJSON(w, false, http.StatusOK, postInformations)
}

func fetchPosts(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusUnprocessableEntity,
			response.Message{Message: messages.ErrInvalidCheck.Error()},
		)

		return
	}

	page := findBegin(r)

	postInformations := make([]model.PostInformation, 0)

	posts, err := db.FindPosts(r.Context(), page)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	for _, post := range posts {
		redisKey := fmt.Sprintf("PostInfo.%d", post.ID)

		result, err := redis.Get(r.Context(), redisKey)
		if err == nil {
			resultModel := model.PostInformation{}
			err = json.Unmarshal([]byte(result), &resultModel)

			if err == nil {
				postInformations = append(postInformations, resultModel)

				continue
			}
		}

		postInformation, err := createPostInformation(r.Context(), post, username)
		if err != nil {
			logger.Error(err)
			response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

			return
		}

		err = redis.SetEx(r.Context(), redisKey, *postInformation, redisCache)
		if err != nil {
			logger.Error(err)
		}

		postInformations = append(postInformations, *postInformation)
	}

	response.ResultResponseJSON(w, false, http.StatusOK, postInformations)
}

func fetchIndividualPost(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusUnprocessableEntity,
			response.Message{Message: messages.ErrInvalidCheck.Error()},
		)

		return
	}

	postID, err := extractID(r, postParam)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	redisKey := fmt.Sprintf("PostInfo.%d", postID)

	result, err := redis.Get(r.Context(), redisKey)
	if err == nil {
		resultModel := model.PostInformation{}

		err = json.Unmarshal([]byte(result), &resultModel)
		if err == nil {
			response.ResultResponseJSON(w, false, http.StatusOK, result)

			return
		}
	}

	postInfo, err := createPostInformationWithFetchPosts(r.Context(), postID, username)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	err = redis.SetEx(r.Context(), redisKey, *postInfo, redisCache)
	if err != nil {
		logger.Error(err)
	}

	response.ResultResponseJSON(w, false, http.StatusOK, postInfo)
}
