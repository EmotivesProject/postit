package api

import "testing"

func TestCreateRouter(t *testing.T) {
	r := CreateRouter()
	if r == nil {
		t.Error("Assertion error")
	}
}
