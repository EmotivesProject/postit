package api

import (
	"context"
	"net/http"
	"postit/internal/db"
	"postit/model"
	"strconv"

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
	ctx context.Context, postID int, username string,
) (*model.PostInformation, error) {
	post, err := db.FindPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	return createPostInformation(ctx, post, username)
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
