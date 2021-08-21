// +build integration

package api_test

import (
	"net/http"
	"postit/test"
	"testing"

	commonTest "github.com/TomBowyerResearchProject/common/test"
	"github.com/stretchr/testify/assert"
)

func TestHealthz(t *testing.T) {
	test.SetUpIntegrationTest()

	req, _ := http.NewRequest("GET", test.TS.URL+"/healthz", nil)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusOK)

	test.TearDownIntegrationTest()
}

func TestRouterFromAuth(t *testing.T) {
	test.SetUpIntegrationTest()

	_, token := commonTest.CreateNewUser(t, test.UaclUserEndpoint)

	req, _ := http.NewRequest("GET", test.TS.URL+"/user", nil)
	req.Header.Add("Authorization", token)

	r, _, _ := commonTest.CompleteTestRequest(t, req)

	assert.EqualValues(t, r.StatusCode, http.StatusOK)

	test.TearDownIntegrationTest()
}
