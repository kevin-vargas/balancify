package server

import (
	"authorization-server/cmd/config"
	"authorization-server/internal/service"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Default(app fiber.Router) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})
}

func Api(app fiber.Router, s service.Oauth, cfg config.Config) {
	app.Get("/certs/jwks", func(c *fiber.Ctx) error {
		cert, err := s.GetCertificateJWKS()
		if err != nil {
			log.Errorf("on get certificate jwks, %v", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		err = c.JSON(cert)
		if err != nil {
			log.Errorf("on encode json, %v", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		return nil
	})
	app.Get("/oauth/callback", func(c *fiber.Ctx) error {
		code := c.Query("code")
		if code == "" {
			return c.Status(http.StatusBadRequest).SendString("no code query")
		}
		exchangedToken, err := s.Exchange(code)
		if err != nil {
			log.Errorf("on exchange token, %v", err)
			return c.SendStatus(http.StatusBadRequest)
		}
		token, err := s.GenerateJWT(exchangedToken)
		if err != nil {
			log.Errorf("on generate jwt, %v", err)
			return c.SendStatus(http.StatusBadRequest)
		}
		origin, err := url.QueryUnescape(c.Query("state"))
		if err != nil {
			log.Errorf("on url query escape, %v", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		tokenCookie := &fiber.Cookie{
			Name:     "token",
			Value:    token,
			HTTPOnly: true,
			Path:     "/",
		}
		loggedCookie := &fiber.Cookie{
			Name:  "logged",
			Value: "true",
			Path:  "/",
		}
		if cfg.Secure == "true" {
			tokenCookie.Domain = cfg.ParentDomain
			tokenCookie.Secure = true
			loggedCookie.Domain = cfg.ParentDomain
			loggedCookie.Secure = true
		}
		c.Cookie(tokenCookie)
		c.Cookie(loggedCookie)
		return c.Redirect(origin, http.StatusSeeOther)
	})
}
