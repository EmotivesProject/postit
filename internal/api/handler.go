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
	"strings"
	"time"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	var alreadyLiked model.Like
	likeCollection := postitDatabase.Collection("likes")
	err = likeCollection.FindOne(context.TODO(), bson.M{"post": post.ID, "user": user.ID}).Decode(&alreadyLiked)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			like := &model.Like{
				ID:      primitive.NewObjectID(),
				Post:    post.ID,
				User:    user.ID,
				Created: time.Now(),
				Active:  true,
			}

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
			return
		} else {
			messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
			return
		}
	}

	messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: "Already liked"})
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

func fetchUserFromAuth(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(userID)
	usernameString := fmt.Sprintf("%v", username)
	user, err := findUser(usernameString)
	if err != nil {
		messageResponseJSON(w, http.StatusBadRequest, model.Message{Message: err.Error()})
		return
	}
	resultResponseJSON(w, http.StatusOK, user)
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

	userPipeline := fmt.Sprintf(`
[
  { "$lookup": {
    "from": "likes",
    "let": { "user_likes": "$user_likes" },
    "pipeline": [
       { "$match": { "$expr": { "$in": [ "$_id", {"$ifNull": ["$$user_likes",[]]} ] } } },
       { "$lookup": {
         "from": "users",
         "let": { "user": "$user" },
         "pipeline": [
           { "$match": { "$expr": { "$eq": [ "$_id", {"$ifNull": ["$$user",[]]} ] } } }
         ],
         "as": "user"
       }}
     ],
     "as": "user_likes"
  }},
  { "$lookup": {
    "from": "users",
    "let": { "user": "$user" },
    "pipeline": [
       { "$match": { "$expr": { "$eq": [ "$_id", {"$ifNull": ["$$user",[]]} ] } } }
     ],
     "as": "user"
  }},
  { "$lookup": {
    "from": "comments",
    "let": { "user_comments": "$user_comments" },
    "pipeline": [
       { "$match": { "$expr": { "$in": [ "$_id", {"$ifNull": ["$$user_comments",[]]} ] } } },
       { "$lookup": {
         "from": "users",
         "let": { "user": "$user" },
         "pipeline": [
           { "$match": { "$expr": { "$eq": [ "$_id", {"$ifNull": ["$$user",[]]} ] } } }
         ],
         "as": "user"
       }}
     ],
     "as": "user_comments"
  }},
  { "$skip": %d},
  { "$limit": 5},
  { "$sort": { "_id": -1 }}
 ]`, begin)

	userPipelineMongo := MongoPipeline(userPipeline)

	options := options.Aggregate()

	collection := postitDatabase.Collection("posts")
	cur, err := collection.Aggregate(context.Background(), userPipelineMongo, options)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var showsWithInfo []bson.M
	if err := cur.All(context.TODO(), &showsWithInfo); err != nil {
		panic(err)
	}

	resultResponseJSON(w, http.StatusOK, showsWithInfo)
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

// Thanks https://github.com/simagix/mongo-go-examples
func MongoPipeline(str string) mongo.Pipeline {
	var pipeline = []bson.D{}
	str = strings.TrimSpace(str)
	if strings.Index(str, "[") != 0 {
		var doc bson.D
		bson.UnmarshalExtJSON([]byte(str), false, &doc)
		pipeline = append(pipeline, doc)
	} else {
		bson.UnmarshalExtJSON([]byte(str), false, &pipeline)
	}
	fmt.Println(pipeline)
	return pipeline
}
