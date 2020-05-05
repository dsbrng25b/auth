package auth

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func TokenAuthenticator(extract ExtractTokenFunc, auth AuthTokenFunc) Authenticator {
	return AuthenticatorFunc(func(r *http.Request) (*Subject, error) {
		return auth(r.Context(), extract(r))
	})
}

type ExtractTokenFunc func(*http.Request) (token string)

func ExtractBearerToken(r *http.Request) (token string) {
	const prefix = "Bearer "
	token = r.Header.Get("Authorization")
	if len(token) < len(prefix) || !strings.EqualFold(token[:len(prefix)], prefix) {
		return ""
	}
	token = token[len(prefix):]
	return token
}

func ExtractHeader(header string) ExtractTokenFunc {
	return func(r *http.Request) (token string) {
		return r.Header.Get(header)
	}
}

func ExtractCookie(name string) ExtractTokenFunc {
	return func(r *http.Request) string {
		c, err := r.Cookie(name)
		if err != nil {
			return ""
		}
		return c.Value
	}
}

type AuthTokenFunc func(ctx context.Context, token string) (*Subject, error)

func TokenFileAuth(file string) (AuthTokenFunc, error) {
	tokens := map[string]Subject{}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	l := 0
	for scanner.Scan() {
		l++
		if scanner.Text() == "" {
			continue
		}
		parts := strings.Fields(scanner.Text())
		if len(parts) < 2 {
			return nil, fmt.Errorf("failed to read token file on line: %d", l)
		}
		sub := Subject{
			Name: parts[1],
		}
		if len(parts) > 2 {
			sub.Groups = strings.Split(parts[2], ",")
		}
		tokens[parts[0]] = sub
	}
	if scanner.Err() != nil {
		return nil, err
	}
	return func(_ context.Context, token string) (*Subject, error) {
		if sub, ok := tokens[token]; ok {
			return &sub, nil
		}
		return nil, nil
	}, nil
}

func TokenMapAuth(tokens map[string]string) AuthTokenFunc {
	return func(_ context.Context, token string) (*Subject, error) {
		if subject, ok := tokens[token]; ok {
			return &Subject{subject, nil}, nil
		}
		return nil, nil
	}
}

func AllTokenAuth() AuthTokenFunc {
	return func(_ context.Context, token string) (*Subject, error) {
		return &Subject{token, nil}, nil
	}
}
