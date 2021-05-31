// +build integration

package api_test

import (
	"fmt"
	"net/http"
	"postit/test"
	"strings"
	"testing"

	"github.com/TomBowyerResearchProject/common/logger"
	commonTest "github.com/TomBowyerResearchProject/common/test"
	"github.com/stretchr/testify/assert"
)

func TestPostDeleteIndividual(t *testing.T) {
	test.SetUpIntegrationTest()

	token := commonTest.CreateNewUser(t, "http://0.0.0.0:8082/user")

	id := test.CreatePost(t, token)

	url := fmt.Sprintf("%s/post/%s", test.TS.URL, id)

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusOK)

	test.TearDownIntegrationTest()
}

func TestCommentDelete(t *testing.T) {
	test.SetUpIntegrationTest()

	token := commonTest.CreateNewUser(t, "http://0.0.0.0:8082/user")

	id := test.CreatePost(t, token)

	url := fmt.Sprintf("%s/post/%s/comment", test.TS.URL, id)

	requestBody := strings.NewReader(
		"{\"message\": \"HELLO\" }",
	)

	req, _ := http.NewRequest("POST", url, requestBody)
	req.Header.Add("Authorization", token)

	r, resultMap, _ := commonTest.CompleteTestRequest(t, req)
	r.Body.Close()

	logger.Infof("%v", resultMap)

	CommentID := int64(resultMap["id"].(float64))

	url = fmt.Sprintf("%s/post/%s/comment/%d", test.TS.URL, id, CommentID)

	req, _ = http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", token)

	r, _, _ = commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusOK)

	test.TearDownIntegrationTest()
}

func TestLikeDelete(t *testing.T) {
	test.SetUpIntegrationTest()

	token := commonTest.CreateNewUser(t, "http://0.0.0.0:8082/user")

	id := test.CreatePost(t, token)

	url := fmt.Sprintf("%s/post/%s/like", test.TS.URL, id)

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", token)

	r, resultMap, _ := commonTest.CompleteTestRequest(t, req)
	r.Body.Close()

	likeID := int64(resultMap["id"].(float64))

	url = fmt.Sprintf("%s/post/%s/like/%d", test.TS.URL, id, likeID)

	req, _ = http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", token)

	r, _, _ = commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusOK)

	test.TearDownIntegrationTest()
}
