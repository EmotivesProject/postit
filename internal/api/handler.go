package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"postit/internal/db"
	"postit/model"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	limit             int64 = 5
	postitDatabase          = db.Connect()
	errFailedDecoding       = errors.New("Failed during decoding request")
)

func healthz(w http.ResponseWriter, r *http.Request) {
	messageResponseJSON(w, http.StatusOK, model.Message{Message: "Health ok"})
}

func createPost(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userID).(string)
	if username == "" {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: "yo"})
		return
	}

	user, err := findUser(username)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	post := &model.Post{}
	err = json.NewDecoder(r.Body).Decode(post)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: errFailedDecoding.Error()})
		return
	}

	post.Created = time.Now()
	post.User = user.ID
	post.ID = primitive.NewObjectID()
	post.Active = true

	postitCollection := postitDatabase.Collection("posts")
	result, err := postitCollection.InsertOne(context.TODO(), post)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	resultResponseJSON(w, http.StatusCreated, result)
}

func createComment(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userID).(string)
	if username == "" {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: "yo"})
		return
	}

	user, err := findUser(username)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	postID := chi.URLParam(r, "post_id")
	if postID == "" {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: "not specified"})
		return
	}
	post, err := findPostById(postID)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	comment := &model.Comment{}
	err = json.NewDecoder(r.Body).Decode(comment)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: errFailedDecoding.Error()})
		return
	}

	comment.ID = primitive.NewObjectID()
	comment.Post = post.ID
	comment.User = user.ID
	comment.Created = time.Now()
	comment.Active = true

	commentsCollection := postitDatabase.Collection("comments")
	result, err := commentsCollection.InsertOne(context.TODO(), comment)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	post.UserComments = append(post.UserComments, comment.ID)
	postitCollection := postitDatabase.Collection("posts")
	_, err = postitCollection.ReplaceOne(context.TODO(), bson.M{"_id": post.ID}, post)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	resultResponseJSON(w, http.StatusCreated, result)
}

func createLike(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userID).(string)
	if username == "" {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: "yo"})
		return
	}

	user, err := findUser(username)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	postID := chi.URLParam(r, "post_id")
	if postID == "" {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: "not specified"})
		return
	}
	post, err := findPostById(postID)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	like := &model.Like{
		ID:      primitive.NewObjectID(),
		Post:    post.ID,
		User:    user.ID,
		Created: time.Now(),
		Active:  true,
	}

	likeCollection := postitDatabase.Collection("likes")
	result, err := likeCollection.InsertOne(context.TODO(), like)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	post.Likes = post.Likes + 1

	post.UserLikes = append(post.UserLikes, like.ID)
	postitCollection := postitDatabase.Collection("posts")
	_, err = postitCollection.ReplaceOne(context.TODO(), bson.M{"_id": post.ID}, post)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	resultResponseJSON(w, http.StatusCreated, result)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: errFailedDecoding.Error()})
		return
	}
	user.ID = primitive.NewObjectID()

	userCollection := postitDatabase.Collection("users")
	result, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	fmt.Println("SAVED NEW USER")

	resultResponseJSON(w, http.StatusCreated, result)
}

func findUser(username string) (model.User, error) {
	user := model.User{}
	filter := bson.D{primitive.E{Key: "username", Value: username}}

	usersCollection := postitDatabase.Collection("users")
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&user)
	return user, err
}

func fetchUser(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	user, err := findUser(username)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	resultResponseJSON(w, http.StatusOK, user)
}

func fetchPost(w http.ResponseWriter, r *http.Request) {
	pagesParam := r.URL.Query().Get("page")
	if pagesParam == "" {
		pagesParam = "0"
	}
	skip, err := strconv.ParseInt(pagesParam, 10, 64)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	begin := limit * skip

	filter := bson.M{}
	posts := []model.Post{}
	options := options.Find()
	options.SetLimit(limit)
	options.SetSkip(begin)
	options.SetSort(map[string]int{"_id": -1})

	postitCollection := postitDatabase.Collection("posts")
	cur, err := postitCollection.Find(context.TODO(), filter, options)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	for cur.Next(context.TODO()) {
		p := model.Post{}
		err := cur.Decode(&p)
		if err != nil {
			continue
		}
		posts = append(posts, p)
	}
	cur.Close(context.TODO())

	resultResponseJSON(w, http.StatusOK, posts)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "post_id")
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
	if likeID == "" {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: "not specified"})
		return
	}
	hex, err := primitive.ObjectIDFromHex(likeID)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	filter := bson.M{"_id": hex}

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

func findPostById(postID string) (model.Post, error) {
	var post model.Post
	hex, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return post, err
	}
	filter := bson.M{"_id": hex}

	postsCollection := postitDatabase.Collection("posts")
	err = postsCollection.FindOne(context.TODO(), filter).Decode(&post)
	return post, err
}
