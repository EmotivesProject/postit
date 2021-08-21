// +build integration

package api_test

import (
	"fmt"
	"net/http"
	"postit/test"
	"testing"

	commonTest "github.com/TomBowyerResearchProject/common/test"
	"github.com/stretchr/testify/assert"
)

func TestPostDeleteIndividual(t *testing.T) {
	test.SetUpIntegrationTest()

	_, token := commonTest.CreateNewUser(t, test.UaclUserEndpoint)

	id := test.CreatePost(t, token)

	url := fmt.Sprintf("%s/post/%s", test.TS.URL, id)

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusOK)

	test.TearDownIntegrationTest()
}

func TestLikeDelete(t *testing.T) {
	test.SetUpIntegrationTest()

	_, token := commonTest.CreateNewUser(t, test.UaclUserEndpoint)

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
