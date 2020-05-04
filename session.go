package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	sessionTokenLength = 60
)

func handlerErr(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

type SessionHandler struct {
	store      SubjectStore
	cookieName string
	duration   time.Duration
}

func NewDefaultSessionHandler() *SessionHandler {
	return &SessionHandler{
		NewMemorySubjectStore(),
		"id",
		time.Second * 60 * 3,
	}
}

func (s *SessionHandler) Authenticate(r *http.Request) (*Subject, error) {
	token := ExtractCookie(s.cookieName)(r)
	if token == "" {
		return nil, nil
	}
	return s.store.Get(r.Context(), token)
}

func (s *SessionHandler) Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sub := GetSubject(r)
		if sub == nil {
			handlerErr(w, fmt.Errorf("no subject to start session"))
			return
		}

		sessionToken, err := generateRandomString(sessionTokenLength)
		if err != nil {
			handlerErr(w, err)
			return
		}

		err = s.store.Set(r.Context(), sessionToken, sub)
		if err != nil {
			handlerErr(w, err)
			return
		}

		c := &http.Cookie{
			Name:   s.cookieName,
			Value:  sessionToken,
			MaxAge: int(s.duration.Seconds()),
		}
		http.SetCookie(w, c)

		//TODO: make this more generic
		http.Redirect(w, r, "/", 302)
	})
}

func (s *SessionHandler) Logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := ExtractCookie(s.cookieName)(r)
		if token == "" {
			http.Redirect(w, r, "/", 302)
		}
		err := s.store.Delete(r.Context(), token)
		if err != nil {
			handlerErr(w, err)
		}
		c := &http.Cookie{
			Name:   s.cookieName,
			MaxAge: -1,
		}
		http.SetCookie(w, c)
		//TODO: make more generic
		http.Redirect(w, r, "/", 302)
	})
}
