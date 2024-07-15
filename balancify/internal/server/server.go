package server

import (
	"balancify/cmd/config"
	"balancify/internal/transaction"
	"context"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Default(app fiber.Router) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})
}

type Service interface {
	ParseCSV(context.Context, io.Reader) ([]transaction.Trx, error)
	SendEmail(context.Context, string, string) error
	ProcessTrxs(context.Context, string, []transaction.Trx) error
}

func Api(app fiber.Router, s Service, cfg config.Config) {
	app.Post("/upload", func(c *fiber.Ctx) error {
		form, err := c.FormFile("file")
		if err != nil {
			log.Errorf("on get file from from, %v", err)
			return c.Status(http.StatusBadRequest).SendString("no file")
		}
		email := c.FormValue("email")
		if email == "" {
			return c.Status(http.StatusBadRequest).SendString("no email")
		}
		f, err := form.Open()
		if err != nil {
			log.Errorf("on open file, %v", err)
			return c.Status(http.StatusBadRequest).SendString("invalid file")
		}
		defer f.Close()
		trxs, err := s.ParseCSV(c.Context(), f)
		if err != nil {
			log.Errorf("on parse csv, %v", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		err = s.ProcessTrxs(c.Context(), email, trxs)
		if err != nil {
			log.Errorf("on process trxs, %v", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		err = s.SendEmail(c.Context(), cfg.From, email)
		if err != nil {
			log.Errorf("on send email, %v", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		return nil
	})
}
