package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_http_basic_auth(t *testing.T) {
	req, err := http.NewRequest("GET", "/welcome", nil)
	if err != nil {
		t.Fatal(err)
	}

	authFunc := userAuthenticator(map[string]string{"myuser": "mypw"})

	var innerRequest *http.Request
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		innerRequest = r
	})

	handler := NewAuthHandler(&UserAuthenticator{basicAuthExtractor, authFunc})(testHandler)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	user := GetSubject(innerRequest)
	if user != "" {
		t.Errorf("status not ok: want='', got='%s'", user)
	}

	req, err = http.NewRequest("GET", "/welcome", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("myuser", "mypw")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	user = GetSubject(innerRequest)
	if user != "myuser" {
		t.Errorf("wrong subject: want='myuser', got='%s'", user)
	}
}
