package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_welcome(t *testing.T) {
	req, err := http.NewRequest("GET", "/welcome", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(welcomeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Error("status not ok")
	}

	expected := "Hello, world!"
	if rr.Body.String() != expected {
		t.Errorf("body invalid: want: %v, got: %v", expected, rr.Body.String())
	}
}
