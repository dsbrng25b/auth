package main

import (
	"context"
	"net/http"
)

type contextKey int

const (
	contextKeySubject contextKey = iota
)

type Subject struct {
	Name   string
	Groups []string
}

func GetSubject(r *http.Request) *Subject {
	rawSub := r.Context().Value(contextKeySubject)
	if rawSub == nil {
		return nil
	}

	sub, ok := rawSub.(Subject)
	if !ok {
		return nil
	}
	return &sub
}

func RequestWithSubject(r *http.Request, sub *Subject) *http.Request {
	if sub == nil {
		return r
	}
	ctx := context.WithValue(r.Context(), contextKeySubject, sub)
	return r.WithContext(ctx)
}
