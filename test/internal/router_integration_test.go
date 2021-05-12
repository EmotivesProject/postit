// +build integration

package api_test

import (
	"net/http"
	"postit/test"
	"strings"
	"testing"

	commonTest "github.com/TomBowyerResearchProject/common/test"
	"github.com/stretchr/testify/assert"
)

func TestRouterUserHandling(t *testing.T) {
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
