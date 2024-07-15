package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lestrrat-go/jwx/jwk"
)

var Token = struct{}{}

func GetToken(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(Token).(string)
	return val, ok
}

type tokenClaim struct {
	Token string `json:"token,omitempty"`
	jwt.RegisteredClaims
}

func fetchJWKS(jwksURL string) (jwk.Set, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	set, err := jwk.Fetch(ctx, jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %v", err)
	}
	return set, nil
}

func NewAuth(jwksURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := c.Cookies("token")
		if t == "" {
			return c.Status(http.StatusUnauthorized).SendString("missing token")
		}
		claim := new(tokenClaim)
		set, err := fetchJWKS(jwksURL)
		if err != nil {
			return err
		}
		jwtParsed, err := new(jwt.Parser).ParseWithClaims(
			t,
			claim,
			func(t *jwt.Token) (interface{}, error) {
				keyID, ok := t.Header["kid"].(string)
				if !ok {
					return nil, fmt.Errorf("expecting JWT header to have string kid")
				}

				if keys, ok := set.LookupKeyID(keyID); ok {
					var pubkey interface{}
					if err := keys.Raw(&pubkey); err != nil {
						return nil, fmt.Errorf("failed to parse JWK: %v", err)
					}

					return pubkey, nil
				}
				return nil, fmt.Errorf("key %q not found in JWKS", keyID)

			},
		)
		if err != nil {
			return c.Status(http.StatusUnauthorized).SendString("invalid token")
		}
		claims, ok := jwtParsed.Claims.(*tokenClaim)
		if !ok {
			return c.Status(http.StatusUnauthorized).SendString("invalid token")
		}
		c.Locals(Token, claims.Token)
		c.Next()
		return nil
	}
}
