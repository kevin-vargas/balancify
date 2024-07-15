package config

import (
	"encoding/base64"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	github_api_uri     = "github_api_uri"
	port               = "port"
	timeout            = "timeout"
	allow_origins      = "allow_origins"
	ttl_service        = "ttl_service"
	api_prefix         = "api_prefix"
	upload_uri         = "upload_uri"
	jwks_uri           = "jwks_uri"
	secure             = "secure"
	client_certificate = "client_certificate"
	client_key         = "client_key"
)

type Config struct {
	GithubApiUri      string
	Port              string
	Timeout           time.Duration
	AllowOrigins      string
	APIPrefix         string
	TTLService        time.Duration
	UploadUri         string
	JwksUri           string
	Secure            bool
	ClientCertificate []byte
	ClientKey         []byte
}

func Make() Config {
	defaults := map[string]any{
		jwks_uri:           "http://localhost:8081/certs/jwks",
		upload_uri:         "http://localhost:8082/upload",
		github_api_uri:     "https://api.github.com",
		allow_origins:      "http://localhost,http://localhost:5173",
		ttl_service:        5 * 60,
		api_prefix:         "",
		timeout:            8,
		port:               ":8080",
		secure:             "false",
		client_certificate: Test_client_certificate,
		client_key:         Test_client_key,
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
	clientKey, err := base64.RawStdEncoding.DecodeString(defaults[client_key].(string))
	if err != nil {
		clientKey = []byte{}
	}
	return Config{
		ClientKey:         clientKey,
		ClientCertificate: clientCertificate,
		Secure:            sec,
		JwksUri:           defaults[jwks_uri].(string),
		UploadUri:         defaults[upload_uri].(string),
		GithubApiUri:      defaults[github_api_uri].(string),
		Port:              defaults[port].(string),
		AllowOrigins:      defaults[allow_origins].(string),
		APIPrefix:         defaults[api_prefix].(string),
		Timeout:           time.Duration(defaults[timeout].(int) * int(time.Second)),
		TTLService:        time.Duration(defaults[ttl_service].(int) * int(time.Second)),
	}
}

var Test_client_certificate = `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURxekNDQXBPZ0F3SUJBZ0lVS3FxQnFpTlpJNGJnSnl6THJFOFc0R3J3Vk53d0RRWUpLb1pJaHZjTkFRRUwKQlFBd1pURUxNQWtHQTFVRUJoTUNRVkl4RkRBU0JnTlZCQWdNQzBKMVpXNXZjMEZwY21Wek1SSXdFQVlEVlFRSApEQWxEUVVKQklGWnBaWGN4RERBS0JnTlZCQW9NQTJKbVpqRVFNQTRHQTFVRUN3d0hZbVptSUVSRlZqRU1NQW9HCkExVUVBd3dEWW1abU1CNFhEVEkwTURjeE5UQXpNRFV5TUZvWERUTTBNRGN4TXpBek1EVXlNRm93WlRFTE1Ba0cKQTFVRUJoTUNRVkl4RkRBU0JnTlZCQWdNQzBKMVpXNXZjMEZwY21Wek1SSXdFQVlEVlFRSERBbERRVUpCSUZacApaWGN4RERBS0JnTlZCQW9NQTJKbVpqRVFNQTRHQTFVRUN3d0hZbVptSUVSRlZqRU1NQW9HQTFVRUF3d0RZbVptCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBbEdEL3Q0eXQ0OWpzZzM3VFpFMHEKcVZmNzlSbHRVb2dDY0thSnJtWWt2aThrYWl4NmxickJFeE1GWTR3Y05IbmZxME5GbW1mTmFqcTNyRDBmeUJWLwpVaFJkK3h2TEx1SGplSk84dDN6elB1UkQvTkt0YVRCUUJLQndnZno2VUFGWHVoUGJRa2k5Z0JwV1dpZnZzK0pGCkx0eVFsdUNEN3owTjJ5dUxsdGE5Wmc2d2M4VjJrNHVzdko5VFNEWUxVaERyRkJvcHNKQXJZMzRoRm94ZTh4S00KYzVZUlB6ZzY1YTZScEFRZU1zN1Q0NWZrZnlldVpoT1pBRWc2WjZIcGhWUUFMa01zMGhPdjZDUisyUnppS29JQgpYRUZWbzdPWWZFUFZJWDdTRUtvMDQ1TjR2Um5hTFVDNC9pRmRjUzdoTFRaaGhDMUgrZE9GT2tpL2g4YVlqVTg5Cm5RSURBUUFCbzFNd1VUQWRCZ05WSFE0RUZnUVVRbFpiMUs1YmtjK0FpdkkyM1RIMTBPZ0hKU1F3SHdZRFZSMGoKQkJnd0ZvQVVRbFpiMUs1YmtjK0FpdkkyM1RIMTBPZ0hKU1F3RHdZRFZSMFRBUUgvQkFVd0F3RUIvekFOQmdrcQpoa2lHOXcwQkFRc0ZBQU9DQVFFQVBORHEzVmJMUHF0WmE3SUFYc2JDL2dYRC9pMG1XcXJQdnpvZW15Rm9aTmtoCmJOeFdFN1VudktZRlBHYklEeHpRZFYzTGVpTTdqVGJ6RGNBYnh6TEtxWkxIb1U1M3RFSElaVGlTK3hiUzRramMKS09Ubm9KZDNZRXBjdms1MWxmdWpQL0s0cWdtUzQxU3l1VWhSM3B4QVNoMTNFVnlsSk5uamErbWVrTUVibVBBTApkRDdJeDRoejZ0Rmt6Z05WdFFBdE1HTFRmRWVRd3hCTmVqWGlZcnZmTy9tNXJaUitxbTBod3B1UWlZdEdSbzZ3CmJQcjQvcFBDY2V2akFKdmpNL3NrSkJjOFIrNE9NSUhHZ0pCZnBjbVhublFDWVRiMHk3QWxEdHFycXRBSEtnSmwKVGFLckVnSXRENjdVZGZxWWhlQkxXdnhHWm96VS9HeWpnZTQ3Wms2WHFBPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=`
var Test_client_key = `LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2UUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktjd2dnU2pBZ0VBQW9JQkFRQ1VZUCszakszajJPeUQKZnROa1RTcXBWL3YxR1cxU2lBSndwb211WmlTK0x5UnFMSHFWdXNFVEV3VmpqQncwZWQrclEwV2FaODFxT3JlcwpQUi9JRlg5U0ZGMzdHOHN1NGVONGs3eTNmUE0rNUVQODBxMXBNRkFFb0hDQi9QcFFBVmU2RTl0Q1NMMkFHbFphCkorK3o0a1V1M0pDVzRJUHZQUTNiSzR1VzFyMW1EckJ6eFhhVGk2eThuMU5JTmd0U0VPc1VHaW13a0N0amZpRVcKakY3ekVveHpsaEUvT0RybHJwR2tCQjR5enRQamwrUi9KNjVtRTVrQVNEcG5vZW1GVkFBdVF5elNFNi9vSkg3WgpIT0lxZ2dGY1FWV2pzNWg4UTlVaGZ0SVFxalRqazNpOUdkb3RRTGorSVYxeEx1RXRObUdFTFVmNTA0VTZTTCtICnhwaU5UejJkQWdNQkFBRUNnZ0VBRnVaVGdQMyt0bHY3d0oyWnZYQ2xaV1pWVmZIN0s5SWU0a2pwbGRkZzAxTk8KUU82bGFxZGNkZmVwRE1DS2Q5VFpYc0t1b3RKall3STE1NmkxVjNsdDRYcVFPSm1GQmJMS0d3bGVCa21MOXdoZQpyODNLQXFKNHJ3WWQ1d25tamVOdktTSTRaQ1g2elNNRGNiMlpJbitJNHQ4YWw1YUY5aTNMamloTnpsVk0vTXVLClNYYnJvTHorV3hFZzRmV0k0dEVDb3ZmZkZOQWFpSXFHYzJoQmtEWWdMQkJ2bTdEaXFYUFA3Sm1sUkVHazROc28Kb005OUgzOU1zRUlSZlNoWjNtcjI5QjRqYWhuNFJwek1XbytJL1cvdGp2cEZ0ZGNRL0JDeHA4enVERlVYNU5COAp2WmNyTUVkaEV0ZWhldWNNWkY0bmp2bkQzazBxeGx6R2Q1aU16NXJ5WndLQmdRRFFEa2wrekZlelRkWG5MWVBxCnMrY1ZYNGhzL3FWWW9aMm5FVVFmRTREcm1GL0wzbjZNMnhLdHNPV0ZqVjF3Z0kxbHlQY0l2L0FIUS9VT0hId0sKbVgxL3RiVElSdWNISlVCVlFlV1YrWjdWbmZacy94b3g4VXRpckdxMUdlMThiRWJkZDUwcm52cnFTb1V6U0xFYQoyNCtjYllMeVNHNSt2YzFXaUtzODRTckRod0tCZ1FDMmtqbzQ2UEE5NjY0dFdHM3lGcnN1QnB5eDlYUHV1TEFTCk15MVIrSTdTWlcyaCtOTzlqRXlQQ1VPZzR5QktlMHpGSnBiejI1dlRjNzNNYzJVcndaWTBYRzhyKzJDakJzVjIKU0w2YXl1QnYvcUtxeUREbU1RYTBMT2VGT1BocjdpaXdvWnNWd21qMFVkR2pFR3Vyb2FmaXFRTnVPc0NPeCtyNApHb2V3SjJqR3V3S0JnQmg0VmJTUUhCQitxeFhSaUo1bUlsdWxMTXFFK0xWLzdLYmxwUGx5dGNyLzFPU0plcURlCmMwZnlja3hPNEJxSFJCb2dsTTEycGFoMUdiRmJNRXVlMmQvWFl6ZmEvdmtjTElEYWkwSWtaY1lDR2lXZnExa00KWkMxcTBmSVM1cGVudEgzL3Y0Q041ajBBSHNKMVhqOG1hN0dlUDdSM1NHZW5zeXJtVUIyTTdoYXBBb0dCQUpDNgpzTjhHZ2RTRWJjcFNySzNhS1Q1ZVRYK2h0ZXJMakFDUmcxN1U3TnVMUG5MRlg0MkdsL1pZQUwyYzc2ODd1V3NjCm9WUGxoczBFbHJScDBnenk1TkRUYWVueTEvUEUzV3BjVm9VOVNOaGZnckppQ3FtZ3VkREJQRFBYS3MvY3QzTDIKV1l2UlZ3US9qREY1UmZHRU1DTzFtaHViQmFUcWhMRnp6cGJ0VnRrYkFvR0FIcktta2JlTmNBTkViU1VFdXNZRgpJMEVJd1p5ZU8yaHY1SWZWSVNLcmFBRWk2cTJOYndsMmE2R05naVM0RTRhU1A0MnVETTNQeC9OTmJ2TmFaWmlsCjM2LzlaN3hQbDE0aVpmdjJlbm44a0ZCVmNzZEIxYnVoL0RVbktpOHdSVG9UVGkvWUY5MzltNHhoanRsYWJWU3YKUG9DOUNYU1JwcUlUeUhoMG90YmRnM009Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K`
