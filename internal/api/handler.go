package api

import (
	"context"
	"fmt"
	"net/http"
	"postit/internal/db"
	"postit/model"
	"postit/pkg/postit_messages"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	postParam            = "post_id"
	limit          int64 = 5
	postitDatabase       = db.Connect()
	msgHealthOK          = "Health ok"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	messageResponseJSON(w, http.StatusOK, model.Message{Message: msgHealthOK})
}

func createPost(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userID).(string)
	user, err := db.FindUser(username, postitDatabase)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	result, err := db.CreatePost(
		r.Body,
		user.ID,
		postitDatabase,
	)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	resultResponseJSON(w, http.StatusCreated, result)
}

func createComment(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userID).(string)
	user, err := db.FindUser(username, postitDatabase)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	postID := chi.URLParam(r, postParam)
	post, err := db.FindPostById(postID, postitDatabase)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	comment, err := db.CreateComment(
		r.Body,
		user.ID,
		postitDatabase,
	)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	post.UserComments = append(post.UserComments, comment.ID)
	err = db.UpdatePost(
		&post,
		postitDatabase,
	)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	resultResponseJSON(w, http.StatusCreated, comment)
}

func createLike(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userID).(string)
	user, err := db.FindUser(username, postitDatabase)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	postID := chi.URLParam(r, postParam)
	post, err := db.FindPostById(postID, postitDatabase)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	// Check if the user has already liked the post
	for _, v := range post.UserLikes {
		if v == user.ID {
			messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: postit_messages.ErrAlreadyLiked.Error()})
			return
		}
	}

	like, err := db.CreateLike(
		user.ID,
		postitDatabase,
	)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	post.Likes = post.Likes + 1
	post.UserLikes = append(post.UserLikes, like.ID)
	err = db.UpdatePost(
		&post,
		postitDatabase,
	)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	resultResponseJSON(w, http.StatusCreated, like)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user, err := db.CreateUser(
		r.Body,
		postitDatabase,
	)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	resultResponseJSON(w, http.StatusCreated, user)
}

func fetchUserFromAuth(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userID)
	usernameString := fmt.Sprintf("%v", username)
	user, err := db.FindUser(usernameString, postitDatabase)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	resultResponseJSON(w, http.StatusOK, user)
}

func fetchPost(w http.ResponseWriter, r *http.Request) {
	begin := findBegin(r)

	postPipelineMongo := db.FindPostPipeline(begin)

	rawData, err := db.GetRawResponseFromAggregate(
		db.PostCollection,
		postPipelineMongo,
		postitDatabase,
	)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	resultResponseJSON(w, http.StatusOK, rawData)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, postParam)
	if postID == "" {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: "not specified"})
		return
	}

	hex, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	filter := bson.M{"_id": hex}

	var post model.Post
	postsCollection := postitDatabase.Collection("posts")
	err = postsCollection.FindOne(context.TODO(), filter).Decode(&post)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	post.Active = false
	_, err = postsCollection.ReplaceOne(context.TODO(), bson.M{"_id": post.ID}, post)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	resultResponseJSON(w, http.StatusOK, post)
}

func deleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "comment_id")
	if commentID == "" {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: "not specified"})
		return
	}
	hex, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	filter := bson.M{"_id": hex}

	var comment model.Comment
	commentsCollection := postitDatabase.Collection("comments")
	err = commentsCollection.FindOne(context.TODO(), filter).Decode(&comment)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	comment.Active = false
	_, err = commentsCollection.ReplaceOne(context.TODO(), bson.M{"_id": comment.ID}, comment)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	resultResponseJSON(w, http.StatusOK, comment)
}

func deleteLike(w http.ResponseWriter, r *http.Request) {
	likeID := chi.URLParam(r, "like_id")
	likeHexID, err := primitive.ObjectIDFromHex(likeID)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	filter := bson.M{"_id": likeHexID}

	var like model.Like
	likesCollection := postitDatabase.Collection("likes")
	err = likesCollection.FindOne(context.TODO(), filter).Decode(&like)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	like.Active = false
	_, err = likesCollection.ReplaceOne(context.TODO(), bson.M{"_id": like.ID}, like)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	resultResponseJSON(w, http.StatusOK, like)
}
