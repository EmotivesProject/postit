package api

import (
	"context"
	"net/http"
	"postit/internal/db"
	"postit/messages"
	"postit/model"
	"strconv"

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

	return createPostInformation(ctx, post, user.Username)
}

func createPostInformation(ctx context.Context, post model.Post, username string) (*model.PostInformation, error) {
	comments, err := db.FindCommentsForPost(ctx, post.ID, true)
	if err != nil {
		return nil, err
	}

	emojiCounts := createEmojiCountsFromComments(comments, "")

	selfEmojiCounts := createEmojiCountsFromComments(comments, username)

	likes, err := db.FindLikesForPost(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	postInfo := model.PostInformation{
		Post:           post,
		Comments:       comments,
		EmojiCount:     emojiCounts,
		SelfEmojiCount: selfEmojiCounts,
		Likes:          likes,
	}

	return &postInfo, nil
}

func createEmojiCountsFromComments(comments []model.Comment, username string) string {
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
