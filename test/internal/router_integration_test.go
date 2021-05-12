// +build integration

package api_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"postit/test"
	"strings"
	"testing"
	"time"

	commonTest "github.com/TomBowyerResearchProject/common/test"
	"github.com/stretchr/testify/assert"
)

func TestRouterCreatePost(t *testing.T) {
	rand.Seed(time.Now().Unix())
	test.SetUpIntegrationTest()

	token := commonTest.CreateNewUser(t, "http://0.0.0.0:8082/user")

	requestBody := strings.NewReader(
		"{\"content\": {\"message\": \"HELLO\"} }",
	)

	req, _ := http.NewRequest("POST", test.TS.URL+"/post", requestBody)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusCreated)

	test.TearDownIntegrationTest()
}

func TestRouterFromAuth(t *testing.T) {
	rand.Seed(time.Now().Unix())
	test.SetUpIntegrationTest()

	token := commonTest.CreateNewUser(t, "http://0.0.0.0:8082/user")

	req, _ := http.NewRequest("GET", test.TS.URL+"/user", nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusOK)

	test.TearDownIntegrationTest()
}

func TestRouterCreateLike(t *testing.T) {
	rand.Seed(time.Now().Unix())
	test.SetUpIntegrationTest()

	token := commonTest.CreateNewUser(t, "http://0.0.0.0:8082/user")

	id := createPost(t, token)

	url := fmt.Sprintf("%s/post/%s/like", test.TS.URL, id)

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusCreated)

	test.TearDownIntegrationTest()
}

func TestRouterCreateComment(t *testing.T) {
	rand.Seed(time.Now().Unix())
	test.SetUpIntegrationTest()

	token := commonTest.CreateNewUser(t, "http://0.0.0.0:8082/user")

	id := createPost(t, token)

	url := fmt.Sprintf("%s/post/%s/comment", test.TS.URL, id)

	requestBody := strings.NewReader(
		"{\"content\": {\"message\": \"HELLO\"} }",
	)

	req, _ := http.NewRequest("POST", url, requestBody)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusCreated)

	test.TearDownIntegrationTest()
}

func createPost(t *testing.T, token string) string {
	requestBody := strings.NewReader(
		"{\"content\": {\"message\": \"HELLO\"} }",
	)

	req, _ := http.NewRequest("POST", test.TS.URL+"/post", requestBody)
	req.Header.Add("Authorization", token)

	_, resultMap, _ := commonTest.CompleteTestRequest(t, req)

	intID := int64(resultMap["id"].(float64))
	return fmt.Sprintf("%d", intID)
}
