package auth

import (
	"crypto/x509"
	"net/http"
)

func DefaultTLSAuthenticator() Authenticator {
	return TLSAuthenticator(extractCertSubject)
}

func TLSAuthenticator(extract ExtractCertSubjectFunc) Authenticator {
	return AuthenticatorFunc(func(r *http.Request) (*Subject, error) {
		if r.TLS == nil {
			return nil, nil
		}

		if len(r.TLS.PeerCertificates) < 1 {
			return nil, nil
		}

		cert := r.TLS.PeerCertificates[0]
		return extract(cert), nil
	})
}

type ExtractCertSubjectFunc func(*x509.Certificate) *Subject

func extractCertSubject(c *x509.Certificate) *Subject {
	return &Subject{
		Name:   c.Subject.CommonName,
		Groups: c.Subject.OrganizationalUnit,
	}
}
