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
	postitCollection        = db.Connect()
	errFailedDecoding       = errors.New("Failed during decoding request")
)

func healthz(w http.ResponseWriter, r *http.Request) {
	messageResponseJSON(w, http.StatusOK, model.Message{Message: "Health ok"})
}

func createPost(w http.ResponseWriter, r *http.Request) {
	post := &model.Post{}
	err := json.NewDecoder(r.Body).Decode(post)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: errFailedDecoding.Error()})
		return
	}

	post.Created = time.Now()
	post.User = "FetchedFromAuth"
	post.ID = primitive.NewObjectID()

	result, err := postitCollection.InsertOne(context.TODO(), post)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	resultResponseJSON(w, http.StatusBadRequest, result)
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

	result, err := postitCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}

	resultResponseJSON(w, http.StatusOK, result)
}
