package main

import (
	"authorization-server/certs"
	"authorization-server/cmd/config"
	"authorization-server/internal/server"
	"authorization-server/internal/service"
	"authorization-server/internal/service/jwt"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"golang.org/x/oauth2"
)

const (
	cert_file = "cert/cert.pem"
	key_file  = "cert/key.pem"
)

func main() {
	cfg := config.Make()

	pk, err := certs.MakePrivateKey(cfg.PrivateKey)
	if err != nil {
		os.Exit(1)
	}
	var oauthConfig = oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     cfg.OauthEndpoint,
	}
	tm := jwt.NewTokenManager(pk, cfg.Kid, cfg.TokenDuration)
	so := service.NewOauth(oauthConfig, tm)
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(logger.New())
	server.Default(app)
	server.Api(app, so, cfg)

	fmt.Printf("Starting server... %+v\n", cfg)
	err = app.Listen(cfg.Port)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
