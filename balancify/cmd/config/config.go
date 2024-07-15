package config

import (
	"encoding/base64"
	"os"
	"strconv"
	"strings"
)

const (
	from               = "from"
	dsn                = "dsn"
	port               = "port"
	smtp_address       = "smtp_address"
	client_certificate = "client_certificate"
	secure             = "secure"
)

type Config struct {
	From              string
	DSN               string
	Port              string
	SMTPAddress       string
	Secure            bool
	ClientCertificate []byte
}

func Make() Config {
	defaults := map[string]any{
		client_certificate: Test_client_certificate,
		smtp_address:       "localhost:1025",
		from:               "test@test.com",
		dsn:                "root:password@(localhost:3306)/balancify",
		port:               ":8082",
		secure:             "false",
	}

	for k := range defaults {
		if v, ok := os.LookupEnv(strings.ToUpper(k)); ok {
			defaults[k] = v
		}
	}
	sec, err := strconv.ParseBool(defaults[secure].(string))
	if err != nil {
		sec = false
	}
	clientCertificate, err := base64.RawStdEncoding.DecodeString(defaults[client_certificate].(string))
	if err != nil {
		clientCertificate = []byte{}
	}
	return Config{
		Secure:            sec,
		ClientCertificate: clientCertificate,
		SMTPAddress:       defaults[smtp_address].(string),
		From:              defaults[from].(string),
		DSN:               defaults[dsn].(string),
		Port:              defaults[port].(string),
	}
}

var Test_client_certificate = `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURxekNDQXBPZ0F3SUJBZ0lVS3FxQnFpTlpJNGJnSnl6THJFOFc0R3J3Vk53d0RRWUpLb1pJaHZjTkFRRUwKQlFBd1pURUxNQWtHQTFVRUJoTUNRVkl4RkRBU0JnTlZCQWdNQzBKMVpXNXZjMEZwY21Wek1SSXdFQVlEVlFRSApEQWxEUVVKQklGWnBaWGN4RERBS0JnTlZCQW9NQTJKbVpqRVFNQTRHQTFVRUN3d0hZbVptSUVSRlZqRU1NQW9HCkExVUVBd3dEWW1abU1CNFhEVEkwTURjeE5UQXpNRFV5TUZvWERUTTBNRGN4TXpBek1EVXlNRm93WlRFTE1Ba0cKQTFVRUJoTUNRVkl4RkRBU0JnTlZCQWdNQzBKMVpXNXZjMEZwY21Wek1SSXdFQVlEVlFRSERBbERRVUpCSUZacApaWGN4RERBS0JnTlZCQW9NQTJKbVpqRVFNQTRHQTFVRUN3d0hZbVptSUVSRlZqRU1NQW9HQTFVRUF3d0RZbVptCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBbEdEL3Q0eXQ0OWpzZzM3VFpFMHEKcVZmNzlSbHRVb2dDY0thSnJtWWt2aThrYWl4NmxickJFeE1GWTR3Y05IbmZxME5GbW1mTmFqcTNyRDBmeUJWLwpVaFJkK3h2TEx1SGplSk84dDN6elB1UkQvTkt0YVRCUUJLQndnZno2VUFGWHVoUGJRa2k5Z0JwV1dpZnZzK0pGCkx0eVFsdUNEN3owTjJ5dUxsdGE5Wmc2d2M4VjJrNHVzdko5VFNEWUxVaERyRkJvcHNKQXJZMzRoRm94ZTh4S00KYzVZUlB6ZzY1YTZScEFRZU1zN1Q0NWZrZnlldVpoT1pBRWc2WjZIcGhWUUFMa01zMGhPdjZDUisyUnppS29JQgpYRUZWbzdPWWZFUFZJWDdTRUtvMDQ1TjR2Um5hTFVDNC9pRmRjUzdoTFRaaGhDMUgrZE9GT2tpL2g4YVlqVTg5Cm5RSURBUUFCbzFNd1VUQWRCZ05WSFE0RUZnUVVRbFpiMUs1YmtjK0FpdkkyM1RIMTBPZ0hKU1F3SHdZRFZSMGoKQkJnd0ZvQVVRbFpiMUs1YmtjK0FpdkkyM1RIMTBPZ0hKU1F3RHdZRFZSMFRBUUgvQkFVd0F3RUIvekFOQmdrcQpoa2lHOXcwQkFRc0ZBQU9DQVFFQVBORHEzVmJMUHF0WmE3SUFYc2JDL2dYRC9pMG1XcXJQdnpvZW15Rm9aTmtoCmJOeFdFN1VudktZRlBHYklEeHpRZFYzTGVpTTdqVGJ6RGNBYnh6TEtxWkxIb1U1M3RFSElaVGlTK3hiUzRramMKS09Ubm9KZDNZRXBjdms1MWxmdWpQL0s0cWdtUzQxU3l1VWhSM3B4QVNoMTNFVnlsSk5uamErbWVrTUVibVBBTApkRDdJeDRoejZ0Rmt6Z05WdFFBdE1HTFRmRWVRd3hCTmVqWGlZcnZmTy9tNXJaUitxbTBod3B1UWlZdEdSbzZ3CmJQcjQvcFBDY2V2akFKdmpNL3NrSkJjOFIrNE9NSUhHZ0pCZnBjbVhublFDWVRiMHk3QWxEdHFycXRBSEtnSmwKVGFLckVnSXRENjdVZGZxWWhlQkxXdnhHWm96VS9HeWpnZTQ3Wms2WHFBPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=`
