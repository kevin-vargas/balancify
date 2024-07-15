package server

import (
	"bff-server/cmd/config"
	"bff-server/internal/entity"
	"bff-server/internal/middleware"
	"context"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Default(app fiber.Router) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})
}

type Service interface {
	GetUser(context.Context) (*entity.Data[entity.User], error)
	UploadFile(context.Context, io.Reader) error
}

func Api(app fiber.Router, s Service, cfg config.Config) {
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     cfg.AllowOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))
	resources := app.Group("/", middleware.NewAuth(cfg.JwksUri))
	resources.Get("/user", func(c *fiber.Ctx) error {
		u, err := s.GetUser(c.Context())
		if err != nil {
			log.Errorf("on get user, %v", err)
			return c.Status(http.StatusBadRequest).SendString(err.Error())
		}
		c.JSON(u)
		return nil
	})
	resources.Post("/upload", func(c *fiber.Ctx) error {
		form, err := c.FormFile("file")
		if err != nil {
			log.Errorf("on get file from form, %v", err)
			return c.Status(http.StatusBadRequest).SendString("no file")
		}
		f, err := form.Open()
		if err != nil {
			log.Errorf("on open file, %v", err)
			return c.Status(http.StatusBadRequest).SendString("invalid file")
		}
		defer f.Close()

		err = s.UploadFile(c.Context(), f)
		if err != nil {
			log.Errorf("on upload file, %v", err)
			return c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		}
		return nil
	})
}
