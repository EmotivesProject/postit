package main

import (
	"log"
	"net/http"
	"os"
	"postit/internal/api"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/middlewares"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
	"github.com/TomBowyerResearchProject/common/verification"
)

func main() {
	logger.InitLogger("postit")

	verification.Init(verification.VerificationConfig{
		VerificationURL: os.Getenv("VERIFICATION_URL"),
	})

	middlewares.Init(middlewares.Config{
		AllowedOrigin:  "*",
		AllowedMethods: "GET,POST,DELETE,OPTIONS",
		AllowedHeaders: "*",
	})

	router := api.CreateRouter()

	err := commonPostgres.Connect(commonPostgres.Config{
		URI: os.Getenv("DATABASE_URL"),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal(http.ListenAndServe(os.Getenv("HOST")+":"+os.Getenv("PORT"), router))
}
