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
		"gowebly, htmx example page, go with htmx",               // define meta keywords
		"Welcome to example! You're here because it worked out.", // define meta description
	)
	bodyContent := pages.BodyContent(
		"Welcome to example!",                // define h1 text
		"You're here because it worked out.", // define p text
	)

	// Define template handler.
	templateHandler := templ.Handler(
		templates.Layout(
			"Welcome to example!", // define title text
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

	return c.SendString("<p>🎉 Yes, <strong>htmx</strong> is ready to use! (<code>GET /api/hello-world</code>)</p>")
}
