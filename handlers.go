package main

import (
	"github.com/BrianMwangi21/anti-discover.git/templates"
	"github.com/BrianMwangi21/anti-discover.git/templates/pages"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2/middleware/adaptor"

	"github.com/gofiber/fiber/v2"
)

// indexViewHandler handles a view for the index page.
func indexViewHandler(c *fiber.Ctx) error {

	// Define template functions.
	metaTags := pages.MetaTags(
		"Anti-Discover",
		"Spotify's discover weekly rogue twin",
	)
	bodyContent := pages.BodyContent(
		"Anti-Discover",
		"You're here because you want something outside your radar. We got you!",
	)

	// Define template handler.
	templateHandler := templ.Handler(
		templates.Layout(
			"Anti-Discover", // define title text
			metaTags, bodyContent,
		),
	)

	// Render template layout.
	return adaptor.HTTPHandler(templateHandler)(c)

}

// showContentAPIHandler handles an API endpoint to show content.
func showContentAPIHandler(c *fiber.Ctx) error {
	// Check, if the current request has a 'HX-Request' header.
	// For more information, see https://htmx.org/docs/#request-headers
	if c.Get("HX-Request") == "" || c.Get("HX-Request") != "true" {
		// If not, return HTTP 400 error.
		return fiber.NewError(fiber.StatusBadRequest, "non-htmx request")
	}

	return c.SendString("<p>🎉 Yes! We connected to Spotify!</p>")
}
