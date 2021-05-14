package test

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"postit/internal/api"
	"strings"
	"testing"
	"time"

	"github.com/TomBowyerResearchProject/common/logger"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
	commonTest "github.com/TomBowyerResearchProject/common/test"
	"github.com/TomBowyerResearchProject/common/verification"
)

var TS *httptest.Server

func SetUpIntegrationTest() {
	rand.Seed(time.Now().Unix())

	logger.InitLogger("postit")

	err := commonPostgres.Connect(commonPostgres.Config{
		URI: "postgres://tom:tom123@localhost:5435/postit_db",
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	verification.Init(verification.VerificationConfig{
		VerificationURL: "http://0.0.0.0:8082/authorize",
	})

	router := api.CreateRouter()

	TS = httptest.NewServer(router)
}

func TearDownIntegrationTest() {
	con := commonPostgres.GetDatabase()
	con.Close()

	TS.Close()
}

func CreatePost(t *testing.T, token string) string {
	requestBody := strings.NewReader(
		"{\"content\": {\"message\": \"HELLO\"} }",
	)

	req, _ := http.NewRequest("POST", TS.URL+"/post", requestBody)
	req.Header.Add("Authorization", token)

	r, resultMap, _ := commonTest.CompleteTestRequest(t, req)
	r.Body.Close()

	intID := int64(resultMap["id"].(float64))

	return fmt.Sprintf("%d", intID)
}
