package db

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindPostPipeline(begin int64) mongo.Pipeline {
	userPipeline := fmt.Sprintf(`
[
  { "$lookup": {
    "from": "likes",
    "let": { "user_likes": "$user_likes" },
    "pipeline": [
       { "$match": { "$and": [ { "$expr": { "$in": [ "$_id", {"$ifNull": ["$$user_likes",[]]} ] } }, { "active": true } ] } },
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
       { "$match": { "$and": [ { "$expr": { "$in": [ "$_id", {"$ifNull": ["$$user_comments",[]]} ] } }, { "active": true } ] } },
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
  { "$match":  { "active": true } },
  { "$skip": %d},
  { "$limit": 5},
  { "$sort": { "_id": -1 }}
 ]`, begin)

	return mongoPipeline(userPipeline)
}

// Thanks https://github.com/simagix/mongo-go-examples
func mongoPipeline(str string) mongo.Pipeline {
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
