package entity

type Claim struct {
	Token string `json:"token,omitempty"`
}

type Jwks struct {
	Keys []JwksKey `json:"keys"`
}
type JwksKey struct {
	N   string `json:"n"`
	E   string `json:"e"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
}
