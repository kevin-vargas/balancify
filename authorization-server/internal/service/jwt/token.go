package jwt

import (
	"authorization-server/internal/entity"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	iss = "balancify"
)

type Manager[T any, K any] interface {
	PrivateKey() *rsa.PrivateKey
	Kid() string
	Generate(c T) (string, error)
	Validate(string) (*K, error)
	Duration() time.Duration
}

type TokenManager = Manager[entity.Claim, TokenClaim]

func getBaseClaims(et time.Duration) jwt.RegisteredClaims {
	d := time.Now()
	expiresAt := d.Add(et)
	return jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(d),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Issuer:    iss,
	}
}

type tokenManager struct {
	kid        string
	privateKey *rsa.PrivateKey
	expireTime time.Duration
}

const (
	kid_header = "kid"
)

type TokenClaim struct {
	entity.Claim
	jwt.RegisteredClaims
}

func (j *tokenManager) Kid() string {
	return j.kid
}

func (j *tokenManager) PrivateKey() *rsa.PrivateKey {
	return j.privateKey
}

func (j *tokenManager) Generate(cp entity.Claim) (string, error) {
	claim := TokenClaim{
		Claim:            cp,
		RegisteredClaims: getBaseClaims(j.expireTime),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	token.Header[kid_header] = j.Kid()

	tokenString, err := token.SignedString(j.PrivateKey())
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *tokenManager) Validate(t string) (*TokenClaim, error) {
	claim := new(TokenClaim)
	token, err := jwt.ParseWithClaims(
		t,
		claim,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected method: %s", token.Header["alg"])
			}
			return &j.privateKey.PublicKey, nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*TokenClaim)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("expired token")
	}
	return claim, nil
}

func (j *tokenManager) Duration() time.Duration {
	return j.expireTime
}

func NewTokenManager(privateKey *rsa.PrivateKey, kid string, expireTime time.Duration) TokenManager {
	return &tokenManager{
		kid:        kid,
		privateKey: privateKey,
		expireTime: expireTime,
	}
}
