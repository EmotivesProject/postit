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

//nolint
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

	comments, err := db.FindCommentsForPost(r.Context(), post.ID, false)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	emojiCounts := createEmojiCountsFromComments(comments)

	likes, err := db.FindLikesForPost(r.Context(), post.ID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	postInfo := model.PostInformation{
		Post:       *post,
		Comments:   comments,
		EmojiCount: emojiCounts,
		Likes:      likes,
	}

	logger.Infof("Successfully created post %s", username)

	response.ResultResponseJSON(w, false, http.StatusCreated, postInfo)
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

	comment, err := db.CreateComment(
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

	post, err := db.FindPostByID(r.Context(), postID)
	if err != nil {
		logger.Error(err)
	} else if post.Username != username {
		send.Comment(post.Username, username, post.ID)
	}

	logger.Infof("Created comment for %s", username)
	response.ResultResponseJSON(w, false, http.StatusCreated, comment)
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
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	comment, err := updateComment(r.Context(), commentID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	logger.Infof("Successfully deleted comment %d", commentID)
	response.ResultResponseJSON(w, false, http.StatusOK, comment)
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
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	for _, post := range posts {
		comments, err := db.FindCommentsForPost(r.Context(), post.ID, false)
		if err != nil {
			logger.Error(err)
			response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

			return
		}

		likes, err := db.FindLikesForPost(r.Context(), post.ID)
		if err != nil {
			logger.Error(err)
			response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

			return
		}

		emojiCounts := createEmojiCountsFromComments(comments)

		postInformation := model.PostInformation{
			Post:       post,
			Comments:   comments,
			EmojiCount: emojiCounts,
			Likes:      likes,
		}

		postInformations = append(postInformations, postInformation)
	}

	response.ResultResponseJSON(w, false, http.StatusOK, postInformations)
}

func fetchJustPosts(w http.ResponseWriter, r *http.Request) {
	page := findBegin(r)

	posts, err := db.FindPosts(r.Context(), page)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	response.ResultResponseJSON(w, false, http.StatusOK, posts)
}

func getCommentsForPost(w http.ResponseWriter, r *http.Request) {
	postID, err := extractID(r, postParam)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	comments, err := db.FindCommentsForPost(r.Context(), postID, true)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	commentsInformation := model.CommentsInformation{
		PostID:   postID,
		Comments: comments,
	}

	response.ResultResponseJSON(w, false, http.StatusOK, commentsInformation)
}

func fetchIndividualPost(w http.ResponseWriter, r *http.Request) {
	postID, err := extractID(r, postParam)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	post, err := db.FindPostByID(r.Context(), postID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	comments, err := db.FindCommentsForPost(r.Context(), postID, false)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusInternalServerError, response.Message{Message: err.Error()})

		return
	}

	emojiCounts := createEmojiCountsFromComments(comments)

	likes, err := db.FindLikesForPost(r.Context(), postID)
	if err != nil {
		logger.Error(err)
		response.MessageResponseJSON(w, false, http.StatusBadRequest, response.Message{Message: err.Error()})

		return
	}

	postInfo := model.PostInformation{
		Post:       post,
		Comments:   comments,
		EmojiCount: emojiCounts,
		Likes:      likes,
	}

	response.ResultResponseJSON(w, false, http.StatusOK, postInfo)
}

func createEmojiCountsFromComments(comments []model.Comment) []model.EmojiCount {
	counts := make([]model.EmojiCount, 0)

	var totalString string

	for _, v := range comments {
		totalString += v.Message
	}

	for _, individualEmoji := range totalString {
		i := exists(counts, string(individualEmoji))
		if i == -1 {
			newCount := model.EmojiCount{
				Emoji: string(individualEmoji),
				Count: 1,
			}
			counts = append(counts, newCount)
		} else {
			counts[i].Count++
		}
	}

	return counts
}

func exists(current []model.EmojiCount, needle string) int {
	for i, count := range current {
		if count.Emoji == needle {
			return i
		}
	}

	return -1
}
