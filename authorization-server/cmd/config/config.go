package config

import (
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

const (
	oauth_endpoint = "oauth_endpoint"
	private_key    = "private_key"
	client_id      = "client_id"
	client_secret  = "client_secret"
	token_url      = "token_url"
	port           = "port"
	kid            = "kid"
	token_duration = "token_duration"
	secure         = "secure"
)

type Config struct {
	Kid           string
	TokenDuration time.Duration
	PrivateKey    string
	ClientID      string
	Port          string
	ClientSecret  string
	OauthEndpoint oauth2.Endpoint
	Secure        string
}

func Make() Config {
	defaults := map[string]any{
		token_duration: 60 * 3000,
		kid:            "0-0-0-1",
		private_key:    Test_rsa_key,
		client_id:      "",
		client_secret:  "",
		oauth_endpoint: endpoints.GitHub,
		secure:         "false",
		port:           ":8081",
	}

	for k := range defaults {
		if v, ok := os.LookupEnv(strings.ToUpper(k)); ok {
			defaults[k] = v
		}
	}
	return Config{
		Kid:           defaults[kid].(string),
		TokenDuration: time.Duration(defaults[token_duration].(int)) * time.Second,
		PrivateKey:    defaults[private_key].(string),
		ClientID:      defaults[client_id].(string),
		ClientSecret:  defaults[client_secret].(string),
		OauthEndpoint: defaults[oauth_endpoint].(oauth2.Endpoint),
		Secure:        defaults[secure].(string),
		Port:          defaults[port].(string),
	}
}

var Test_rsa_key = `-----BEGIN RSA PRIVATE KEY-----
MIIBPAIBAAJBAMtX4/MbUQx13uZ8upe2t6jfDdyInpuKfxff2+JQVfShhzfhgZCo
sUvpDFBYpJ/PGrXSogdXk1dnolhlOVj2bfsCAwEAAQJAftLrdlXkP/xIMMM8caFh
fS7Za2G+Ys6HpDFX6Bgo9DCow2BnYxZotin9/cp5x6SghEzKmDKZreGwzPrj/0jH
IQIhAPs4lJntWFfWRH3tmJw2P8znWcQdQ01lRSr7+5GiiBirAiEAzzYmk6eC59bD
04yj8JMbH9SQaP2UuMf4RkFAhRiPn/ECIQCUgjAmho5Q7pNytgAfaFpy8Ni5/GqK
2DD5ZhijUSePHQIhAMSBfw4SEuPYWTfbLXGtoFCMXjMjIJIoGfxOT2ipRTORAiEA
7I2XRtEsn4q60OPy0XUk05yYzXMPKLlg9DyHsyYgK1Y=
-----END RSA PRIVATE KEY-----
`
