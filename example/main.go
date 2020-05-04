package main

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/dsbrng25b/auth"
	"gopkg.in/square/go-jose.v2"

	"github.com/dsbrng25b/auth/jwt"
)

func main() {
	//authFunc := auth.UserMapAuth(map[string]string{"dave": "foobar"})
	authFunc := auth.AuthenticateAll()
	tokenAuthFunc := auth.TokenMapAuth(map[string]string{"123456": "davem2m"})

	userAuth := auth.UserAuthenticator(auth.ExtractFormValue, authFunc)
	tokenAuth := auth.TokenAuthenticator(auth.ExtractBearerToken, tokenAuthFunc)
	tlsAuth := auth.DefaultTLSAuthenticator()

	session := auth.NewDefaultSessionHandler()
	jwtIssuer, err := jwt.IssueJWTHandler(jose.HS256, []byte("fooobarbla"))
	if err != nil {
		log.Fatal(err)
	}
	jwtAuth := jwt.JWTAuthenticator([]byte("fooobarbla1"))
	jwtTokenAuth := auth.TokenAuthenticator(auth.ExtractHeader("X-Token"), jwtAuth)

	http.Handle("/login/session", auth.AuthHandler(userAuth)(session.Login()))
	http.Handle("/login/jwt", auth.AuthHandler(userAuth)(jwtIssuer))
	http.Handle("/logout", session.Logout())
	http.Handle("/jwt", auth.AuthHandler(jwtTokenAuth)(http.HandlerFunc(auth.UserInfoHandler)))

	http.Handle("/", auth.AuthHandler(session)(http.HandlerFunc(auth.UserInfoHandler)))

	http.Handle("/token", auth.AuthHandler(tlsAuth, tokenAuth, userAuth)(http.HandlerFunc(auth.UserInfoHandler)))

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
