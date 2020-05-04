package auth

import (
	"context"
	"net/http"
	"strings"
)

func TokenAuthenticator(extract ExtractTokenFunc, auth AuthTokenFunc) Authenticator {
	return AuthenticatorFunc(func(r *http.Request) (*Subject, error) {
		return auth(r.Context(), extract(r))
	})
}

type ExtractTokenFunc func(*http.Request) (token string)

func ExtractBearerToken(r *http.Request) (token string) {
	const prefix = "Bearer "
	token = r.Header.Get("Authorization")
	if len(token) < len(prefix) || !strings.EqualFold(token[:len(prefix)], prefix) {
		return ""
	}
	token = token[len(prefix):]
	return token
}

func ExtractHeader(header string) ExtractTokenFunc {
	return func(r *http.Request) (token string) {
		return r.Header.Get(header)
	}
}

func ExtractCookie(name string) ExtractTokenFunc {
	return func(r *http.Request) string {
		c, err := r.Cookie(name)
		if err != nil {
			return ""
		}
		return c.Value
	}
}

type AuthTokenFunc func(ctx context.Context, token string) (*Subject, error)

func TokenMapAuth(tokens map[string]string) AuthTokenFunc {
	return func(_ context.Context, token string) (*Subject, error) {
		if subject, ok := tokens[token]; ok {
			return &Subject{subject, nil}, nil
		}
		return nil, nil
	}
}

func AllTokenAuth() AuthTokenFunc {
	return func(_ context.Context, token string) (*Subject, error) {
		return &Subject{token, nil}, nil
	}
}
