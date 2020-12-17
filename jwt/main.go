package jwt

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/dvob/auth"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

//func JWTAuthenticator(extract auth.ExtractTokenFunc, extractSubject ExtractJWTSubjectFunc) {
//}

func DiscoveryAuthenticator(ctx context.Context, issuerURL string, config *oidc.Config) (auth.AuthTokenFunc, error) {
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, err
	}
	verifier := provider.Verifier(config)
	return func(ctx context.Context, token string) (*auth.Subject, error) {
		idToken, err := verifier.Verify(ctx, token)
		if err != nil {
			// TODO: how to handle errors
			return nil, nil
		}
		return extractSubject(idToken), nil
	}, nil
}

func JWTAuthenticator(sharedKey interface{}) auth.AuthTokenFunc {
	return func(_ context.Context, rawToken string) (*auth.Subject, error) {
		token, err := jwt.ParseSigned(rawToken)
		if err != nil {
			//TODO: how to handle errors
			return nil, nil
		}
		out := jwt.Claims{}
		if err := token.Claims(sharedKey, &out); err != nil {
			//TODO: how to handle errors
			return nil, nil
		}
		return &auth.Subject{out.Subject, nil}, nil
	}
}

func handlerErr(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IssueJWTHandler(algorithm jose.SignatureAlgorithm, key interface{}) (http.Handler, error) {
	signingKey := jose.SigningKey{Algorithm: algorithm, Key: key}
	sig, err := jose.NewSigner(signingKey, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return nil, err
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sub := auth.GetSubject(r)
		if sub == nil {
			handlerErr(w, fmt.Errorf("no subject to issue jwt"))
			return
		}

		log.Println("issue token for", sub)

		cl := jwt.Claims{
			Subject:   sub.Name,
			Issuer:    "dvob",
			NotBefore: jwt.NewNumericDate(time.Now()),
			Audience:  jwt.Audience{"dvob"},
		}

		privateCl := struct {
			Groups []string `json:"groups,omitempty"`
		}{
			sub.Groups,
		}

		raw, err := jwt.Signed(sig).Claims(cl).Claims(privateCl).CompactSerialize()
		if err != nil {
			handlerErr(w, err)
		}

		fmt.Fprintf(w, raw)
	}), nil

}

type ExtractJWTSubjectFunc func(token string) *auth.Subject

func extractSubject(token *oidc.IDToken) *auth.Subject {
	return &auth.Subject{token.Subject, nil}
}
