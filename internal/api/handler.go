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
)

var (
	postParam = "post_id"
	likeParam = "like_id"
)

func fetchUserFromAuth(w http.ResponseWriter, r *http.Request) {
	user, err := getUsernameAndGroup(r)
	if err != nil {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusUnprocessableEntity,
			response.Message{Message: messages.ErrInvalidCheck.Error()},
		)

		return
	}

	err = db.CheckUsername(r.Context(), user)
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

	user, err = db.FindByUsername(r.Context(), user.Username)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	response.ResultResponseJSON(w, false, http.StatusOK, user)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	user, err := getUsernameAndGroup(r)
	if err != nil {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusUnprocessableEntity,
			response.Message{Message: messages.ErrInvalidCheck.Error()},
		)

		return
	}

	err = db.CheckUsername(r.Context(), user)
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
		user.Username,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	postInformation, err := createPostInformation(r.Context(), *post, user.Username)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Successfully created post %s", user.Username)

	response.ResultResponseJSON(w, false, http.StatusCreated, postInformation)
}

// nolint
func createComment(w http.ResponseWriter, r *http.Request) {
	user, err := getUsernameAndGroup(r)
	if err != nil {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusUnprocessableEntity,
			response.Message{Message: messages.ErrInvalidCheck.Error()},
		)

		return
	}

	err = db.CheckUsername(r.Context(), user)
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
		user.Username,
		postID,
	)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	postInformation, err := createPostInformationWithFetchPosts(r.Context(), postID, user)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	go send.Comment(postInformation.Post.Username, user.Username, postInformation.Post.ID)

	logger.Infof("Created comment for %s", user.Username)
	response.ResultResponseJSON(w, false, http.StatusCreated, postInformation)
}

func createLike(w http.ResponseWriter, r *http.Request) {
	user, err := getUsernameAndGroup(r)
	if err != nil {
		logger.Error(messages.ErrInvalidCheck)
		response.MessageResponseJSON(
			w,
			false,
			http.StatusUnprocessableEntity,
			response.Message{Message: messages.ErrInvalidCheck.Error()},
		)

		return
	}

	err = db.CheckUsername(r.Context(), user)
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

	logger.Infof("Created like for %s", user.Username)
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

	go send.DeletePost(postID)

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

//nolint
func fetchExplorePosts(w http.ResponseWriter, r *http.Request) {
	user, err := getUsernameAndGroup(r)
	if err != nil {
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

	posts, err := db.FindPostsBasedOnLatAndLng(r.Context(), lat, lng, page, user)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	for _, post := range posts {
		postInformation, err := createPostInformation(r.Context(), post, user.Username)
		if err != nil {
			logger.Error(err)
			response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

			return
		}

		postInformations = append(postInformations, *postInformation)
	}

	response.ResultResponseJSON(w, false, http.StatusOK, postInformations)
}

func fetchPosts(w http.ResponseWriter, r *http.Request) {
	user, err := getUsernameAndGroup(r)
	if err != nil {
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

	posts, err := db.FindPosts(r.Context(), page, user)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	for _, post := range posts {
		postInformation, err := createPostInformation(r.Context(), post, user.Username)
		if err != nil {
			logger.Error(err)
			response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

			return
		}

		postInformations = append(postInformations, *postInformation)
	}

	response.ResultResponseJSON(w, false, http.StatusOK, postInformations)
}

func fetchIndividualPost(w http.ResponseWriter, r *http.Request) {
	user, err := getUsernameAndGroup(r)
	if err != nil {
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

	postInfo, err := createPostInformationWithFetchPosts(r.Context(), postID, user)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	response.ResultResponseJSON(w, false, http.StatusOK, postInfo)
}
