package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"postit/internal/db"
	"postit/messages"
	"postit/model"
	"strconv"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/redis"
	"github.com/TomBowyerResearchProject/common/verification"
	"github.com/go-chi/chi"
)

const bitSize = 64

func findBegin(r *http.Request) int {
	pagesParam := r.URL.Query().Get("page")
	if pagesParam == "" {
		pagesParam = "0"
	}

	skip, err := strconv.Atoi(pagesParam)
	if err != nil {
		return 0
	}

	return db.PostLimit * skip
}

func findPosition(r *http.Request, param string) float64 {
	urlParam := r.URL.Query().Get(param)
	if urlParam == "" {
		urlParam = "0"
	}

	value, err := strconv.ParseFloat(urlParam, bitSize)
	if err != nil {
		return 0
	}

	return value
}

func extractID(r *http.Request, param string) (int, error) {
	paramString := chi.URLParam(r, param)

	return strconv.Atoi(paramString)
}

func createPostInformationWithFetchPosts(
	ctx context.Context, postID int, user model.User,
) (*model.PostInformation, error) {
	post, err := db.FindPostByIDForUser(ctx, postID, user)
	if err != nil {
		return nil, err
	}

	return createPostInformation(ctx, post)
}

func createPostInformation(ctx context.Context, post model.Post) (*model.PostInformation, error) {
	comments, err := db.FindCommentsForPost(ctx, post.ID, true)
	if err != nil {
		return nil, err
	}

	emojiSummary := createEmojiSummaryFromComments(comments, "")

	likes, err := db.FindLikesForPost(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	postInfo := model.PostInformation{
		Post:         post,
		Comments:     comments,
		EmojiSummary: emojiSummary,
		Likes:        likes,
	}

	return &postInfo, nil
}

func createEmojiSummaryFromComments(comments []model.Comment, username string) string {
	var totalString string

	for _, comment := range comments {
		if username == "" {
			totalString += comment.Message
		} else if (username != "") && (comment.Username == username) {
			totalString += comment.Message
		}
	}

	return totalString
}

func getUsernameAndGroup(r *http.Request) (model.User, error) {
	user := model.User{}

	username, ok := r.Context().Value(verification.UserID).(string)
	if !ok {
		return user, messages.ErrInvalidCheck
	}

	userGroup, ok := r.Context().Value(verification.UserGroup).(string)
	if !ok {
		return user, messages.ErrInvalidCheck
	}

	user.Username = username
	user.UserGroup = userGroup

	return user, nil
}

func getUserAndEnsureUserInDB(r *http.Request) (model.User, error) {
	user, err := getUsernameAndGroup(r)
	if err != nil {
		return user, messages.ErrInvalidCheck
	}

	err = db.CheckUsername(r.Context(), user)
	if err != nil {
		return user, messages.ErrInvalidUsername
	}

	return user, nil
}

func fetchPostInformationsFromPosts(ctx context.Context, posts []model.Post) []model.PostInformation {
	postInformations := make([]model.PostInformation, 0)

	for _, post := range posts {
		redisKey := fmt.Sprintf("PostInfo.%d", post.ID)

		result, err := redis.Get(ctx, redisKey)
		if err == nil {
			resultModel := model.PostInformation{}
			err = json.Unmarshal([]byte(result), &resultModel)

			if err == nil {
				postInformations = append(postInformations, resultModel)

				continue
			}
		}

		postInformation, err := createPostInformation(ctx, post)
		if err != nil {
			continue
		}

		err = redis.SetEx(ctx, redisKey, *postInformation, redisCache)
		if err != nil {
			logger.Error(err)
		}

		postInformations = append(postInformations, *postInformation)
	}

	return postInformations
}

func updatePostInRedis(ctx context.Context, postID int, user model.User) (model.PostInformation, error) {
	postInformation, err := createPostInformationWithFetchPosts(ctx, postID, user)
	if err != nil {
		logger.Error(err)

		return model.PostInformation{}, err
	}

	redisKey := fmt.Sprintf("PostInfo.%d", postID)

	err = redis.SetEx(ctx, redisKey, *postInformation, redisCache)
	if err != nil {
		logger.Error(err)
	}

	return *postInformation, nil
}
