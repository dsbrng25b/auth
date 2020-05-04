package main

import (
	"crypto/tls"
	"log"
	"net/http"
)

func main() {
	//authFunc := userAuthenticator(map[string]string{"dave": "foobar"})
	authFunc := authenticateAll()
	tokenAuthFunc := tokenAuthenticator(map[string]string{"123456": "davem2m"})

	userAuth := &UserAuthenticator{formAuthExtractor, authFunc}
	tokenAuth := &TokenAuthenticator{bearerTokenExtractor, tokenAuthFunc}
	tlsAuth := NewDefaultTLSAuthenticator()

	store := NewMemorySubjectStore()
	tokenStoreAuth := tokenStoreAuthenticate(store)
	sessionAuth := &TokenAuthenticator{cookieTokenExtractor("id"), tokenStoreAuth}

	initSession := startSessionHandler("id", store)
	logoutSession := removeSessionHandler("id", store)

	http.Handle("/login", NewAuthHandler(userAuth)(initSession))
	http.Handle("/logout", logoutSession)

	http.Handle("/", NewAuthHandler(sessionAuth)(http.HandlerFunc(userInfoHandler)))

	http.Handle("/token", NewAuthHandler(tlsAuth, tokenAuth, userAuth)(http.HandlerFunc(userInfoHandler)))

	//http.Handle("/", basicAuthMiddleware(authFunc, basicAuthRequestMiddleware("My Realm", http.HandlerFunc(userInfoHandler))))

	go func() {
		log.Println("start http server")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	tlsSrv := &http.Server{
		Addr:    ":8443",
		Handler: http.DefaultServeMux,
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequestClientCert,
		},
	}

	log.Println("start https server")
	log.Fatal(tlsSrv.ListenAndServeTLS("tls.crt", "tls.key"))
}
