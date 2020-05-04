package main

import (
	"context"
	"net/http"
)

type UserAuthenticator struct {
	extract userExtractorFunc
	auth    userAuthFunc
}

func (u *UserAuthenticator) Authenticate(r *http.Request) (*Subject, error) {
	user, pw := u.extract(r)
	return u.auth(r.Context(), user, pw)
}

type userAuthFunc func(ctx context.Context, user, password string) (*Subject, error)

func userAuthenticator(users map[string]string) userAuthFunc {
	return func(_ context.Context, user, password string) (*Subject, error) {
		if pw, ok := users[user]; ok && pw == password {
			return &Subject{user, nil}, nil
		}
		return nil, nil
	}
}

func authenticateAll() userAuthFunc {
	return func(_ context.Context, user, _ string) (*Subject, error) {
		return &Subject{user, nil}, nil
	}
}
