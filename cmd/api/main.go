package main

import (
	"log"
	"net/http"
	"os"
	"postit/internal/api"

	commonKafka "github.com/TomBowyerResearchProject/common/kafka"
	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/middlewares"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
	"github.com/TomBowyerResearchProject/common/verification"

	"github.com/joho/godotenv"
)

func main() {
	logger.InitLogger("postit")

	verification.Init(verification.VerificationConfig{
		VerificationURL: "http://uacl/authorize",
	})

	commonKafka.InitProducer(commonKafka.ConfigProducer{
		Topic:  "EVENT",
		Server: "kafka:9092",
	})

	middlewares.Init(middlewares.Config{
		AllowedOrigin:  "*",
		AllowedMethods: "GET,POST,DELETE,OPTIONS",
		AllowedHeaders: "Accept, Content-Type, Content-Length, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header",
	})

	router := api.CreateRouter()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	commonPostgres.Connect(commonPostgres.Config{
		URI: os.Getenv("databaseURL"),
	})

	log.Fatal(http.ListenAndServe(os.Getenv("HOST")+":"+os.Getenv("PORT"), router))
}
