package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"web-server/cmd/config"
	"web-server/internal/server"

	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/gofiber/fiber/v2"
)

//go:embed all:app/dist
var content embed.FS

const (
	buildPath = "app/dist"
)

func Get() (http.FileSystem, error) {
	fsys, err := fs.Sub(content, buildPath)
	if err != nil {
		return nil, err
	}
	return http.FS(fsys), nil
}
func main() {
	cfg := config.Make()

	httpFileSystem, err := Get()
	if err != nil {
		log.Fatal(err)
		return
	}
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(etag.New())
	app.Use(logger.New())
	server.Default(app)
	server.Web(app, httpFileSystem, cfg)
	fmt.Printf("Starting server... %+v\n", cfg)
	app.Listen(cfg.Port)
}
