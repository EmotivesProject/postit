// +build integration

package api_test

import (
	"fmt"
	"net/http"
	"postit/test"
	"strings"
	"testing"

	commonTest "github.com/TomBowyerResearchProject/common/test"
	"github.com/stretchr/testify/assert"
)

func TestRouterCreatePost(t *testing.T) {
	test.SetUpIntegrationTest()

	_, token := commonTest.CreateNewUser(t, test.UaclUserEndpoint)

	requestBody := strings.NewReader(
		"{\"content\": {\"reaction\": \"😊\"} }",
	)

	req, _ := http.NewRequest("POST", test.TS.URL+"/post", requestBody)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusCreated)

	test.TearDownIntegrationTest()
}

func TestRouterCreateLike(t *testing.T) {
	test.SetUpIntegrationTest()

	_, token := commonTest.CreateNewUser(t, test.UaclUserEndpoint)

	id := test.CreatePost(t, token)

	url := fmt.Sprintf("%s/post/%s/like", test.TS.URL, id)

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusCreated)

	test.TearDownIntegrationTest()
}

func TestRouterCreateComment(t *testing.T) {
	test.SetUpIntegrationTest()

	_, token := commonTest.CreateNewUser(t, test.UaclUserEndpoint)

	id := test.CreatePost(t, token)

	url := fmt.Sprintf("%s/post/%s/comment", test.TS.URL, id)

	requestBody := strings.NewReader(
		"{\"message\": \"😊\"}",
	)

	req, _ := http.NewRequest("POST", url, requestBody)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusCreated)

	test.TearDownIntegrationTest()
}
