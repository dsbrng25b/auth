package main

import (
	"crypto/x509"
	"net/http"
)

type TLSAuthenticator struct {
	extractUser   extractCertSubject
	extractGroups extractCertGroups
}

func NewDefaultTLSAuthenticator() *TLSAuthenticator {
	return &TLSAuthenticator{
		extractCN,
		extractOUs,
	}
}

func (t *TLSAuthenticator) Authenticate(r *http.Request) (authenticated bool, subject string, groups []string, err error) {
	if r.TLS == nil {
		return
	}

	if len(r.TLS.PeerCertificates) < 1 {
		return
	}

	authenticated = true
	cert := r.TLS.PeerCertificates[0]
	if t.extractUser != nil {
		subject = t.extractUser(cert)
	}

	if t.extractGroups != nil {
		groups = t.extractGroups(cert)
	}
	return
}

type extractCertSubject func(*x509.Certificate) string

type extractCertGroups func(*x509.Certificate) []string

func extractCN(c *x509.Certificate) string {
	return c.Subject.CommonName
}

func extractOUs(c *x509.Certificate) []string {
	return c.Subject.OrganizationalUnit
}
