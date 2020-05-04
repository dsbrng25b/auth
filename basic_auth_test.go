package auth

import (
	"net/http"
	"net/http/httptest"
	"reflect"
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

	handler := AuthHandler(UserAuthenticator(BasicAuthExtractor, authFunc))(testHandler)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	sub := GetSubject(innerRequest)
	if sub != nil {
		t.Errorf("subject have to be nil on unauthenticated request: want=nil, got='%v'", sub)
	}

	req, err = http.NewRequest("GET", "/welcome", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("myuser", "mypw")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	sub = GetSubject(innerRequest)
	expected := Subject{"myuser", nil}
	if reflect.DeepEqual(sub, expected) {
		t.Errorf("wrong subject: want='%v', got='%s'", expected, sub)
	}
}
