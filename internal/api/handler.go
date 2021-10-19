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
)

var (
	postParam  = "post_id"
	likeParam  = "like_id"
	redisCache = time.Minute * 10
)

func fetchUserFromAuth(w http.ResponseWriter, r *http.Request) {
	user, err := getUserAndEnsureUserInDB(r)
	if err != nil {
		logger.Error(err)

		response.MessageResponseJSON(
			w,
			false,
			http.StatusBadRequest,
			response.Message{Message: messages.ErrInvalidUsername.Error()},
		)

		return
	}

	user, err = db.FindByUsername(r.Context(), user.Username)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("authorized user for request %s", user.Username)

	response.ResultResponseJSON(w, false, http.StatusOK, user)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	user, err := getUserAndEnsureUserInDB(r)
	if err != nil {
		logger.Error(err)

		response.MessageResponseJSON(
			w,
			false,
			http.StatusBadRequest,
			response.Message{Message: messages.ErrInvalidUsername.Error()},
		)

		return
	}

	post, err := db.CreatePost(
		r.Context(),
		r.Body,
		user.Username,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	postInformation, err := createPostInformation(r.Context(), *post)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	postInfo := postInformation.Post

	postInfoPayload, err := json.Marshal(postInfo)
	if err != nil {
		logger.Error(err)

		return
	}

	logger.Infof("Successfully created post %s with information %s", user.Username, string(postInfoPayload))

	response.ResultResponseJSON(w, false, http.StatusCreated, postInformation)
}

func createComment(w http.ResponseWriter, r *http.Request) {
	user, err := getUserAndEnsureUserInDB(r)
	if err != nil {
		logger.Error(err)

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

	comment, err := db.CreateComment(
		r.Context(),
		r.Body,
		user.Username,
		postID,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	post, err := db.FindPostByIDForUser(r.Context(), postID, user)
	if err != nil {
		logger.Error(err)
	} else if post.Username != user.Username {
		send.Comment(post.Username, user.Username, post.ID)
	}

	postInformation, err := updatePostInRedis(r.Context(), postID, user)
	if err != nil {
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Created comment USER-%s COMMENT-%s", user.Username, comment.Message)
	response.ResultResponseJSON(w, false, http.StatusCreated, postInformation)
}

func createLike(w http.ResponseWriter, r *http.Request) {
	user, err := getUserAndEnsureUserInDB(r)
	if err != nil {
		logger.Error(err)

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
		user.Username,
		postID,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	post, err := db.FindPostByIDForUser(r.Context(), postID, user)
	if err != nil {
		logger.Error(err)
	} else if post.Username != user.Username {
		send.Like(post.Username, user.Username, post.ID)
	}

	_, err = updatePostInRedis(r.Context(), postID, user)
	if err != nil {
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Created like for %s on post %d", user.Username, post.ID)
	response.ResultResponseJSON(w, false, http.StatusCreated, like)
}

func updatePost(ctx context.Context, postID int) (model.Post, error) {
	post, err := db.FindPostByID(ctx, postID)
	if err != nil {
		return post, err
	}

	post.Active = false

	err = db.UpdatePost(ctx, &post)

	logger.Infof("Post id %d is now inactive", post.ID)

	return post, err
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	user, err := getUserAndEnsureUserInDB(r)
	if err != nil {
		logger.Error(err)

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

	post, err := updatePost(r.Context(), postID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	go send.DeletePost(postID)

	_, err = updatePostInRedis(r.Context(), post.ID, user)
	if err != nil {
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Post id %d is now deleted", postID)

	response.ResultResponseJSON(w, false, http.StatusOK, post)
}

func updateLike(ctx context.Context, likeID int) (model.Like, error) {
	like, err := db.FindLikeByID(ctx, likeID)
	if err != nil {
		return like, err
	}

	like.Active = false

	err = db.UpdateLike(ctx, &like)

	logger.Infof("Like id %d from user %s is now inactive", like.ID, like.Username)

	return like, err
}

func deleteLike(w http.ResponseWriter, r *http.Request) {
	user, err := getUserAndEnsureUserInDB(r)
	if err != nil {
		logger.Error(err)

		response.MessageResponseJSON(
			w,
			false,
			http.StatusBadRequest,
			response.Message{Message: messages.ErrInvalidUsername.Error()},
		)

		return
	}

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

	_, err = updatePostInRedis(r.Context(), like.PostID, user)
	if err != nil {
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Like id %d from user %s is now deleted", like.ID, like.Username)

	response.ResultResponseJSON(w, false, http.StatusOK, like)
}

func fetchExplorePosts(w http.ResponseWriter, r *http.Request) {
	user, err := getUserAndEnsureUserInDB(r)
	if err != nil {
		logger.Error(err)

		response.MessageResponseJSON(
			w,
			false,
			http.StatusBadRequest,
			response.Message{Message: messages.ErrInvalidUsername.Error()},
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

	posts, err := db.FindPostsBasedOnLatAndLng(r.Context(), lat, lng, page, user)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	postInformations := fetchPostInformationsFromPosts(r.Context(), posts)

	logger.Infof("User %s requested explore posts on page %d", user.Username, page)

	response.ResultResponseJSON(w, false, http.StatusOK, postInformations)
}

func fetchPosts(w http.ResponseWriter, r *http.Request) {
	user, err := getUserAndEnsureUserInDB(r)
	if err != nil {
		logger.Error(err)

		response.MessageResponseJSON(
			w,
			false,
			http.StatusBadRequest,
			response.Message{Message: messages.ErrInvalidUsername.Error()},
		)

		return
	}

	page := findBegin(r)

	posts, err := db.FindPosts(r.Context(), page, user)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	postInformations := fetchPostInformationsFromPosts(r.Context(), posts)

	logger.Infof("User %s requested posts on page %d", user.Username, page)

	response.ResultResponseJSON(w, false, http.StatusOK, postInformations)
}

func fetchIndividualPost(w http.ResponseWriter, r *http.Request) {
	user, err := getUserAndEnsureUserInDB(r)
	if err != nil {
		logger.Error(err)

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

	redisKey := fmt.Sprintf("PostInfo.%d", postID)

	result, err := redis.Get(r.Context(), redisKey)
	if err == nil {
		resultModel := model.PostInformation{}

		err = json.Unmarshal([]byte(result), &resultModel)
		if err == nil {
			response.ResultResponseJSON(w, false, http.StatusOK, resultModel)

			return
		}
	}

	postInfo, err := createPostInformationWithFetchPosts(r.Context(), postID, user)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	err = redis.SetEx(r.Context(), redisKey, *postInfo, redisCache)
	if err != nil {
		logger.Error(err)
	}

	logger.Infof("User %s requested individual post id %d", user.Username, postID)

	response.ResultResponseJSON(w, false, http.StatusOK, postInfo)
}
