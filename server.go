package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"

	gowebly "github.com/gowebly/helpers"
)

func runServer() error {
	port, err := strconv.Atoi(gowebly.Getenv("BACKEND_PORT", "7000"))
	if err != nil {
		return err
	}

	config := fiber.Config{
		Views:        html.NewFileSystem(http.Dir("./templates"), ".html"),
		ViewsLayout:  "main",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server := fiber.New(config)
	server.Use(logger.New())

	server.Static("/static", "./static")

	server.Get("/", indexViewHandler)
	server.Get("/api/connect", connectToSpotifyHandler)

	return server.Listen(fmt.Sprintf(":%d", port))
}
