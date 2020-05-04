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

func (t *TLSAuthenticator) Authenticate(r *http.Request) (*Subject, error) {
	if r.TLS == nil {
		return nil, nil
	}

	if len(r.TLS.PeerCertificates) < 1 {
		return nil, nil
	}

	cert := r.TLS.PeerCertificates[0]
	sub := &Subject{}
	if t.extractUser != nil {
		sub.Name = t.extractUser(cert)
	}

	if t.extractGroups != nil {
		sub.Groups = t.extractGroups(cert)
	}
	return sub, nil
}

type extractCertSubject func(*x509.Certificate) string

type extractCertGroups func(*x509.Certificate) []string

func extractCN(c *x509.Certificate) string {
	return c.Subject.CommonName
}

func extractOUs(c *x509.Certificate) []string {
	return c.Subject.OrganizationalUnit
}
