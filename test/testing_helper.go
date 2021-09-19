package test

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"postit/internal/api"
	"strings"
	"testing"
	"time"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/middlewares"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
	"github.com/TomBowyerResearchProject/common/redis"
	commonTest "github.com/TomBowyerResearchProject/common/test"
	"github.com/TomBowyerResearchProject/common/verification"
)

const (
	uaclEndpoint     = "http://0.0.0.0:8082"
	UaclUserEndpoint = uaclEndpoint + "/user"
)

var TS *httptest.Server

func SetUpIntegrationTest() {
	rand.Seed(time.Now().Unix())

	logger.InitLogger("postit", logger.EmailConfig{
		From:     os.Getenv("EMAIL_FROM"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		Level:    os.Getenv("EMAIL_LEVEL"),
	})

	middlewares.Init(middlewares.Config{
		AllowedOrigins: "*",
		AllowedMethods: "GET,POST,OPTIONS,DELETE",
		AllowedHeaders: "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token",
	})

	err := commonPostgres.Connect(commonPostgres.Config{
		URI: "postgres://tom:tom123@localhost:5435/postit_db",
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = redis.Init(redis.Config{
		Addr:   "test_redis_db:63790",
		Prefix: "POSTIT",
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	verification.Init(verification.VerificationConfig{
		VerificationURL: uaclEndpoint + "/authorize",
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
		"{\"content\": {\"reaction\": \"😊\"} }",
	)

	req, _ := http.NewRequest("POST", TS.URL+"/post", requestBody)
	req.Header.Add("Authorization", token)

	r, resultMap, _ := commonTest.CompleteTestRequest(t, req)
	r.Body.Close()

	post, _ := resultMap["post"].(map[string]interface{})

	intID := int64(post["id"].(float64))

	return fmt.Sprintf("%d", intID)
}
