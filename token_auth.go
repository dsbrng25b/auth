package main

import (
	"context"
	"net/http"
)

type TokenAuthenticator struct {
	extract tokenExtractorFunc
	auth    tokenAuthFunc
}

func (t *TokenAuthenticator) Authenticate(r *http.Request) (authenticated bool, subject string, groups []string, err error) {
	token := t.extract(r)
	ok, subject, groups, err := t.auth(r.Context(), token)
	return ok, subject, groups, err
}

type tokenAuthFunc func(ctx context.Context, token string) (authenticated bool, subject string, groups []string, err error)

func tokenAuthenticator(tokens map[string]string) tokenAuthFunc {
	return func(_ context.Context, token string) (bool, string, []string, error) {
		if subject, ok := tokens[token]; ok {
			return true, subject, nil, nil
		}
		return false, "", nil, nil
	}
}

func authenticateAllTokens() tokenAuthFunc {
	return func(_ context.Context, token string) (bool, string, []string, error) {
		return true, token, nil, nil
	}
}
