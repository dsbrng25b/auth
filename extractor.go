package main

import (
	"net/http"
	"strings"
)

type userExtractorFunc func(r *http.Request) (user, password string)

func basicAuthExtractor(r *http.Request) (user, password string) {
	user, password, _ = r.BasicAuth()
	return
}

func formAuthExtractor(r *http.Request) (user, password string) {
	user = r.FormValue("username")
	password = r.FormValue("password")
	return
}

type tokenExtractorFunc func(r *http.Request) (token string)

func bearerTokenExtractor(r *http.Request) (token string) {
	const prefix = "Bearer "
	token = r.Header.Get("Authorization")
	if len(token) < len(prefix) || !strings.EqualFold(token[:len(prefix)], prefix) {
		return ""
	}
	token = token[len(prefix):]
	return token
}

func headerTokenExtractor(header string, r *http.Request) (token string) {
	return r.Header.Get(header)
}

func cookieTokenExtractor(name string) tokenExtractorFunc {
	return func(r *http.Request) string {
		c, err := r.Cookie(name)
		if err != nil {
			return ""
		}
		return c.Value
	}
}
