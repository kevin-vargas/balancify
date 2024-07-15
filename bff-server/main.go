package main

import (
	"bff-server/cmd/config"
	"bff-server/internal/repository/github"
	"bff-server/internal/server"
	"bff-server/internal/service"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/fiber/v2"
)

func makeClient(secure bool, cert, key []byte) (*http.Client, error) {
	if secure {
		clientCert, err := tls.X509KeyPair(cert, key)
		if err != nil {
			return nil, err
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{clientCert},
		}
		tr := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		return &http.Client{
			Transport: tr,
		}, nil
	}
	return http.DefaultClient, nil
}

func main() {
	cfg := config.Make()
	repo := github.New(cfg.GithubApiUri)
	c, err := makeClient(cfg.Secure, cfg.ClientCertificate, cfg.ClientKey)
	if err != nil {
		panic(err)
	}
	svc := service.New(repo, cfg.UploadUri, c)
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(logger.New())
	app.Use(etag.New())
	server.Default(app)
	api := app.Group(cfg.APIPrefix)
	server.Api(api, svc, cfg)
	fmt.Printf("Starting server... %+v\n", cfg)
	app.Listen(cfg.Port)
}
