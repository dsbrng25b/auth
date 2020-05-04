package main

import (
	"context"
	"net/http"
)

func UserAuthenticator(extract UserExtractFunc, auth UserAuthFunc) Authenticator {
	return AuthenticatorFunc(func(r *http.Request) (*Subject, error) {
		user, password := extract(r)
		return auth(r.Context(), user, password)
	})
}

type UserExtractFunc func(*http.Request) (user, password string)

func BasicAuthExtractor(r *http.Request) (user, password string) {
	user, password, _ = r.BasicAuth()
	return
}

func FormValueExtractor(r *http.Request) (user, password string) {
	user = r.FormValue("username")
	password = r.FormValue("password")
	return user, password
}

type UserAuthFunc func(ctx context.Context, user, password string) (*Subject, error)

func userAuthenticator(users map[string]string) UserAuthFunc {
	return func(_ context.Context, user, password string) (*Subject, error) {
		if pw, ok := users[user]; ok && pw == password {
			return &Subject{user, nil}, nil
		}
		return nil, nil
	}
}

func authenticateAll() UserAuthFunc {
	return func(_ context.Context, user, _ string) (*Subject, error) {
		return &Subject{user, nil}, nil
	}
}
