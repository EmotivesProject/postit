package test

import (
	"log"
	"net/http/httptest"
	"postit/internal/api"

	"github.com/TomBowyerResearchProject/common/logger"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
	"github.com/TomBowyerResearchProject/common/verification"
)

var TS *httptest.Server

func SetUpIntegrationTest() {
	logger.InitLogger("postit")

	err := commonPostgres.Connect(commonPostgres.Config{
		URI: "postgres://tom:tom123@localhost:5435/postit_db",
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	verification.Init(verification.VerificationConfig{
		VerificationURL: "localhost:8082/authorize",
	})

	router := api.CreateRouter()

	TS = httptest.NewServer(router)
}

func TearDownIntegrationTest() {
	con := commonPostgres.GetDatabase()
	con.Close()

	TS.Close()
}
