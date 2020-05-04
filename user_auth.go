package auth

import (
	"context"
	"net/http"
)

func UserAuthenticator(extract ExtractUserFunc, auth AuthUserFunc) Authenticator {
	return AuthenticatorFunc(func(r *http.Request) (*Subject, error) {
		user, password := extract(r)
		return auth(r.Context(), user, password)
	})
}

type ExtractUserFunc func(*http.Request) (user, password string)

func ExtractBasicAuth(r *http.Request) (user, password string) {
	user, password, _ = r.BasicAuth()
	return
}

func ExtractFormValue(r *http.Request) (user, password string) {
	user = r.FormValue("username")
	password = r.FormValue("password")
	return user, password
}

type AuthUserFunc func(ctx context.Context, user, password string) (*Subject, error)

func UserMapAuth(users map[string]string) AuthUserFunc {
	return func(_ context.Context, user, password string) (*Subject, error) {
		if pw, ok := users[user]; ok && pw == password {
			return &Subject{user, nil}, nil
		}
		return nil, nil
	}
}

func AuthenticateAll() AuthUserFunc {
	return func(_ context.Context, user, _ string) (*Subject, error) {
		return &Subject{user, nil}, nil
	}
}
