package server

import (
	"net/http"
	"web-server/cmd/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

const (
	index = "/index.html"
)

func Default(app fiber.Router) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})
}

func Web(app fiber.Router, fsys http.FileSystem, cfg config.Config) {
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         fsys,
		NotFoundFile: cfg.URIPrefix + index,
	}))
}
