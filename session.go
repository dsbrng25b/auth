package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
)

func tokenGenerator() string {
	b := make([]byte, 40)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func startSessionHandler(cookieName string, store Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subject := GetSubject(r)
		if subject == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		sessionToken := tokenGenerator()
		err := store.Set(r.Context(), sessionToken, []byte(subject))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		c := &http.Cookie{
			Name:   cookieName,
			Value:  sessionToken,
			MaxAge: 99999,
		}
		http.SetCookie(w, c)
		http.Redirect(w, r, "/", 302)
	})
}

func removeSessionHandler(cookieName string, store Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &http.Cookie{
			Name:   cookieName,
			MaxAge: -1,
		}
		http.SetCookie(w, c)
		http.Redirect(w, r, "/", 302)
	})
}
