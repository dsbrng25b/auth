package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/dvob/auth"
	"gopkg.in/square/go-jose.v2"

	"github.com/dvob/auth/jwt"
)

func main() {
	//
	// USER
	//
	userMapAuthFunc := auth.UserMapAuth(map[string]string{"myuser1": "test123"})
	userMapAuth := auth.UserAuthenticator(auth.ExtractBasicAuth, userMapAuthFunc)
	http.Handle("/user/map", auth.AuthHandler(userMapAuth)(auth.UserInfoHandler()))

	userAllAuthFunc := auth.AuthenticateAll()
	userAllAuth := auth.UserAuthenticator(auth.ExtractFormValue, userAllAuthFunc)
	http.Handle("/user/all", auth.AuthHandler(userAllAuth)(auth.UserInfoHandler()))

	//
	// TOKEN
	//
	tokenMapAuthFunc := auth.TokenMapAuth(map[string]string{"token123456": "myuser2"})
	tokenMapAuth := auth.TokenAuthenticator(auth.ExtractBearerToken, tokenMapAuthFunc)
	http.Handle("/token/map", auth.AuthHandler(tokenMapAuth)(auth.UserInfoHandler()))

	tokenFileAuthFunc, err := auth.TokenFileAuth("token.txt")
	if err != nil {
		log.Fatal(err)
	}
	tokenFileAuth := auth.TokenAuthenticator(auth.ExtractBearerToken, tokenFileAuthFunc)
	http.Handle("/token/file", auth.AuthHandler(tokenFileAuth)(auth.UserInfoHandler()))

	//
	// JWT
	//
	// issuer/login
	jwtIssuer, err := jwt.IssueJWTHandler(jose.HS256, []byte("myhmacsecret"))
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/jwt/login", auth.AuthHandler(userAllAuth)(jwtIssuer))

	// jwt auth
	jwtTokenAuthFunc := jwt.JWTAuthenticator([]byte("myhmacsecret"))
	jwtAuth := auth.TokenAuthenticator(auth.ExtractHeader("X-Token"), jwtTokenAuthFunc)
	http.Handle("/jwt", auth.AuthHandler(jwtAuth)(auth.UserInfoHandler()))

	//
	// TLS
	//
	tlsAuth := auth.DefaultTLSAuthenticator()
	http.Handle("/tls", auth.AuthHandler(tlsAuth)(auth.UserInfoHandler()))

	//
	// SESSION
	//
	session := auth.NewDefaultSessionHandler()
	http.Handle("/session/login", auth.AuthHandler(userAllAuth)(session.Login()))
	http.Handle("/session/logout", session.Logout())
	http.Handle("/session", auth.AuthHandler(session)(auth.UserInfoHandler()))

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
