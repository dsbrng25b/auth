package main

import (
	"net/http"
)

type Authenticator interface {
	Authenticate(r *http.Request) (*Subject, error)
}

type AuthenticatorFunc func(r *http.Request) (*Subject, error)

func (a AuthenticatorFunc) Authenticate(r *http.Request) (*Subject, error) {
	return a(r)
}

func AuthHandler(as ...Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return authHandler(next, as...)
	}
}

func authHandler(next http.Handler, as ...Authenticator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			sub *Subject
			err error
		)
		for _, a := range as {
			sub, err = a.Authenticate(r)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			if sub != nil {
				break
			}
		}

		if sub == nil {
			next.ServeHTTP(w, r)
			return
		}

		r = RequestWithSubject(r, sub)
		next.ServeHTTP(w, r)
	})
}
