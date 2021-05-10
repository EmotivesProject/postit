// +build integration

package api_test

import (
	"fmt"
	"postit/test"
	"testing"

	commmonTest "github.com/TomBowyerResearchProject/common/test"
)

func TestRouterUserHandling(t *testing.T) {
	test.SetUpIntegrationTest()

	token := commmonTest.CreateNewUser(t, "http://0.0.0.0:8082/user")
	fmt.Println(token)

	test.TearDownIntegrationTest()
}
