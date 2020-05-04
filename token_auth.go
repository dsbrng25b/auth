package main

import (
	"context"
	"net/http"
	"strings"
)

func TokenAuthenticator(extract TokenExtractFunc, auth TokenAuthFunc) Authenticator {
	return AuthenticatorFunc(func(r *http.Request) (*Subject, error) {
		return auth(r.Context(), extract(r))
	})
}

type TokenExtractFunc func(*http.Request) (token string)

func BearerTokenExtractor(r *http.Request) (token string) {
	const prefix = "Bearer "
	token = r.Header.Get("Authorization")
	if len(token) < len(prefix) || !strings.EqualFold(token[:len(prefix)], prefix) {
		return ""
	}
	token = token[len(prefix):]
	return token
}

func HeaderTokenExtractor(header string, r *http.Request) (token string) {
	return r.Header.Get(header)
}

func CookieTokenExtractor(name string) TokenExtractFunc {
	return func(r *http.Request) string {
		c, err := r.Cookie(name)
		if err != nil {
			return ""
		}
		return c.Value
	}
}

type TokenAuthFunc func(ctx context.Context, token string) (*Subject, error)

func tokenAuthenticator(tokens map[string]string) TokenAuthFunc {
	return func(_ context.Context, token string) (*Subject, error) {
		if subject, ok := tokens[token]; ok {
			return &Subject{subject, nil}, nil
		}
		return nil, nil
	}
}

func authenticateAllTokens() TokenAuthFunc {
	return func(_ context.Context, token string) (*Subject, error) {
		return &Subject{token, nil}, nil
	}
}
