package api

import (
	"log"
	"net/http"
	"postit/internal/db"
	"postit/model"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	mongo, ctx = db.Connect()
)

func healthz(w http.ResponseWriter, r *http.Request) {
	err := mongo.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	messageResponseJSON(w, http.StatusOK, model.Message{Message: "Health ok"})
}
