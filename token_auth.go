package main

import (
	"context"
	"net/http"
)

type TokenAuthenticator struct {
	extract tokenExtractorFunc
	auth    tokenAuthFunc
}

func (t *TokenAuthenticator) Authenticate(r *http.Request) (*Subject, error) {
	token := t.extract(r)
	return t.auth(r.Context(), token)
}

type tokenAuthFunc func(ctx context.Context, token string) (*Subject, error)

func tokenStoreAuthenticate(store SubjectStore) tokenAuthFunc {
	return func(ctx context.Context, token string) (*Subject, error) {
		return store.Get(ctx, token)
	}
}

func tokenAuthenticator(tokens map[string]string) tokenAuthFunc {
	return func(_ context.Context, token string) (*Subject, error) {
		if subject, ok := tokens[token]; ok {
			return &Subject{subject, nil}, nil
		}
		return nil, nil
	}
}

func authenticateAllTokens() tokenAuthFunc {
	return func(_ context.Context, token string) (*Subject, error) {
		return &Subject{token, nil}, nil
	}
}
