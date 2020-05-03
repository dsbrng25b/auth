package main

import (
	"net/http"
)

type Authenticator interface {
	Authenticate(r *http.Request) (authenticated bool, subject string, groups []string, err error)
}

func NewAuthHandler(as ...Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return authHandler(next, as...)
	}
}

func authHandler(next http.Handler, as ...Authenticator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			authenticated bool
			subject       string
			groups        []string
			err           error
		)
		for _, a := range as {
			authenticated, subject, groups, err = a.Authenticate(r)
			if err != nil {
				http.Error(w, "auth failed", 500)
				return
			}
			if authenticated {
				break
			}
		}

		if !authenticated {
			next.ServeHTTP(w, r)
			return
		}

		if subject != "" {
			r = RequestWithSubject(r, subject)
		}

		if groups != nil {
			r = RequestWithGroups(r, groups)
		}

		next.ServeHTTP(w, r)
	})
}
