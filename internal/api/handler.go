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
	encodedID := r.Context().Value(userID).(string)
	if encodedID == "" {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: "yo"})
		return
	}

	post := &model.Post{}
	err := json.NewDecoder(r.Body).Decode(post)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: errFailedDecoding.Error()})
		return
	}

	post.Created = time.Now()
	post.User = encodedID
	post.ID = primitive.NewObjectID()

	postitCollection := postitDatabase.Collection("posts")
	result, err := postitCollection.InsertOne(context.TODO(), post)
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

func fetchUser(w http.ResponseWriter, r *http.Request) {
	user := model.User{}

	encodedID := chi.URLParam(r, "encoded_id")
	filter := bson.D{primitive.E{Key: "encoded_id", Value: encodedID}}

	usersCollection := postitDatabase.Collection("users")
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	resultResponseJSON(w, http.StatusOK, user)
}

func fetchPost(w http.ResponseWriter, r *http.Request) {
	skip, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	begin := 0 + (limit * skip)

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
	fmt.Println(postID)

	hex, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	filter := bson.M{"_id": hex}

	postitCollection := postitDatabase.Collection("posts")
	result, err := postitCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	resultResponseJSON(w, http.StatusOK, result)
}
