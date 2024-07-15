package config

import (
	"os"
	"strings"
)

const (
	port          = "port"
	uri_prefix    = "uri_prefix"
	allow_origins = "allow_origins"
)

type Config struct {
	Port         string
	AllowOrigins string
	URIPrefix    string
}

func Make() Config {
	defaults := map[string]any{
		allow_origins: "http://localhost,http://localhost:5173",
		uri_prefix:    "",
		port:          ":80",
	}

	for k := range defaults {
		if v, ok := os.LookupEnv(strings.ToUpper(k)); ok {
			defaults[k] = v
		}
	}
	return Config{
		Port:         defaults[port].(string),
		AllowOrigins: defaults[allow_origins].(string),
		URIPrefix:    defaults[uri_prefix].(string),
	}
}
