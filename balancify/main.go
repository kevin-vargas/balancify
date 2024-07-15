package main

import (
	"balancify/cmd/config"
	"balancify/internal/server"
	"balancify/internal/service"
	"balancify/internal/transaction/rdbms"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	cfg := config.Make()
	tm, err := rdbms.New(cfg.DSN)
	if err != nil {
		panic(err)
	}
	svc := service.New(cfg.SMTPAddress, tm)
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(logger.New())
	server.Default(app)
	server.Api(app, svc, cfg)
	if cfg.Secure {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(cfg.ClientCertificate)
		tlsConfig := &tls.Config{
			ClientCAs:  caCertPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
		tlsConfig.BuildNameToCertificate()
		app.Server().TLSConfig = tlsConfig
	}
	fmt.Printf("Starting server... %+v\n", cfg)
	app.Listen(cfg.Port)
}
