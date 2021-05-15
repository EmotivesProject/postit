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

func TestPostGet(t *testing.T) {
	test.SetUpIntegrationTest()

	token := commonTest.CreateNewUser(t, "http://0.0.0.0:8082/user")

	test.CreatePost(t, token)

	url := fmt.Sprintf("%s/post", test.TS.URL)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusOK)

	test.TearDownIntegrationTest()
}

func TestPostGetIndividualPost(t *testing.T) {
	test.SetUpIntegrationTest()

	token := commonTest.CreateNewUser(t, "http://0.0.0.0:8082/user")

	id := test.CreatePost(t, token)

	url := fmt.Sprintf("%s/post/%s/", test.TS.URL, id)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusOK)

	test.TearDownIntegrationTest()
}