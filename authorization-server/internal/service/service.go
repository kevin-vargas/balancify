package service

import (
	"authorization-server/internal/entity"
	"authorization-server/internal/service/jwt"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

const (
	Alg = "RS256"
	Kty = "RSA"
	Use = "sig"
	E   = "AQAB"
)

type Oauth interface {
	GetCertificateJWKS() (*entity.Jwks, error)
	GenerateJWT(*oauth2.Token) (string, error)
	Exchange(string) (*oauth2.Token, error)
}

type oauth struct {
	c  *http.Client
	tm jwt.TokenManager
	oauth2.Config
}

func (o *oauth) Exchange(code string) (*oauth2.Token, error) {
	v := url.Values{
		"client_id":     {o.ClientID},
		"client_secret": {o.ClientSecret},
		"code":          {code},
	}
	req, err := http.NewRequest(
		http.MethodPost,
		o.Endpoint.TokenURL,
		strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := o.c.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"request for device code authorisation returned status %v (%v)",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var e oauth2.Token
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (o *oauth) GenerateJWT(t *oauth2.Token) (string, error) {
	c := entity.Claim{
		Token: t.AccessToken,
	}
	return o.tm.Generate(c)
}

func (o *oauth) GetCertificateJWKS() (*entity.Jwks, error) {
	pk := o.tm.PrivateKey().PublicKey
	n := base64.StdEncoding.EncodeToString(pk.N.Bytes())
	return &entity.Jwks{
		Keys: []entity.JwksKey{
			{
				Kid: o.tm.Kid(),
				Alg: Alg,
				Kty: Kty,
				Use: Use,
				N:   n,
				E:   E,
			},
		},
	}, nil
}

func NewOauth(c oauth2.Config, tm jwt.TokenManager) Oauth {
	return &oauth{
		c:      http.DefaultClient,
		tm:     tm,
		Config: c,
	}
}
