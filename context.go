package main

import (
	"context"
	"net/http"
)

type contextKey string

const (
	contextSubjectKey contextKey = "subject"
	contextGroupsKey  contextKey = "groups"
)

func GetSubject(r *http.Request) string {
	u := r.Context().Value(contextSubjectKey)
	if u == nil {
		return ""
	}

	us, _ := u.(string)
	return us
}

func GetGroups(r *http.Request) []string {
	u := r.Context().Value(contextGroupsKey)
	if u == nil {
		return nil
	}

	us, _ := u.([]string)
	return us
}

func RequestWithSubject(r *http.Request, subject string) *http.Request {
	ctx := context.WithValue(r.Context(), contextSubjectKey, subject)
	return r.WithContext(ctx)
}

func RequestWithGroups(r *http.Request, groups []string) *http.Request {
	ctx := context.WithValue(r.Context(), contextGroupsKey, groups)
	return r.WithContext(ctx)
}
