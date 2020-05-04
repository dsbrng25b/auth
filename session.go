package main

import (
	"log"
	"net/http"
)

const (
	sessionTokenLength = 40
)

func handlerErr(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func startSessionHandler(cookieName string, store SubjectStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sub := GetSubject(r)
		if sub == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		sessionToken, err := GenerateRandomString(sessionTokenLength)
		if err != nil {
			handlerErr(w, err)
			return
		}

		err = store.Set(r.Context(), sessionToken, sub)
		if err != nil {
			handlerErr(w, err)
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

func removeSessionHandler(cookieName string, store SubjectStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &http.Cookie{
			Name:   cookieName,
			MaxAge: -1,
		}
		http.SetCookie(w, c)
		http.Redirect(w, r, "/", 302)
	})
}
