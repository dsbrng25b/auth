package main

import (
	"fmt"
	"net/http"
)

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func userInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := GetSubject(r)
	groups := GetGroups(r)
	fmt.Fprintf(w, "user: %s\n", user)
	fmt.Fprintf(w, "groups: %v\n", groups)
	fmt.Fprintf(w, "-----")
}
