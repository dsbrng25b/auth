package main

import (
	"context"
	"net/http"
)

type UserAuthenticator struct {
	extract userExtractorFunc
	auth    userAuthFunc
}

func (u *UserAuthenticator) Authenticate(r *http.Request) (authenticated bool, subject string, groups []string, err error) {
	user, pw := u.extract(r)
	ok, err := u.auth(r.Context(), user, pw)
	return ok, user, nil, err
}

type userAuthFunc func(ctx context.Context, user, password string) (authenticated bool, err error)

func userAuthenticator(users map[string]string) userAuthFunc {
	return func(_ context.Context, user, password string) (bool, error) {
		if pw, ok := users[user]; ok && pw == password {
			return true, nil
		}
		return false, nil
	}
}

func authenticateAll() userAuthFunc {
	return func(_ context.Context, _, _ string) (bool, error) {
		return true, nil
	}
}
