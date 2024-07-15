package certs

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func MakePrivateKey(pkStr string) (*rsa.PrivateKey, error) {
	p, _ := pem.Decode([]byte(pkStr))
	if p == nil {
		return nil, errors.New("error on pem decode, invalid private key")
	}
	pk, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		return nil, err
	}
	return pk, nil
}
