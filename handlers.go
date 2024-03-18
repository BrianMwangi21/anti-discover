package main

import (
	"context"
	"errors"

	"github.com/BrianMwangi21/anti-discover.git/templates"
	"github.com/BrianMwangi21/anti-discover.git/templates/pages"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/zmb3/spotify/v2"
)

const state = "anti-discover"

func indexViewHandler(c *fiber.Ctx) error {
	link, err := getSpotifyLink()

	if err != nil {
		return err
	}

	metaTags := getMetaTags()
	bodyContent := pages.BodyContent(
		"Anti-Discover",
		"You're here because you want something outside your radar. We got you!",
		link,
	)

	templateHandler := templ.Handler(
		templates.Layout("Anti-Discover", metaTags, bodyContent),
	)

	return adaptor.HTTPHandler(templateHandler)(c)
}

func antiHandler(c *fiber.Ctx) error {
	auth := getAuth()

	request, err := convertRequest(&c.Context().Request)
	if err != nil {
		return errorHandler(c, err)
	}

	token, err := retrieveToken(request)
	if err != nil {
		token, err = auth.Token(request.Context(), state, request)
		if err != nil {
			return errorHandler(c, err)
		}

		err = saveToken(request, token)
		if err != nil {
			return errorHandler(c, err)
		}
	}

	if st := request.FormValue("state"); st != state {
		return errorHandler(c, errors.New("State mismatch"))
	}

	client := spotify.New(auth.Client(c.Context(), token))
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		return errorHandler(c, err)
	}

	recommendations, playlist, err := getRecommendationAndCreatePlaylist(client, user.ID)
	if err != nil {
		return errorHandler(c, err)
	}

	metaTags := getMetaTags()
	antiContent := pages.AntiContent(user, recommendations, playlist)

	templateHandler := templ.Handler(
		templates.Layout("Anti-Discover", metaTags, antiContent),
	)

	return adaptor.HTTPHandler(templateHandler)(c)
}

func errorHandler(c *fiber.Ctx, err error) error {
	link, linkErr := getSpotifyLink()

	if linkErr != nil {
		err = linkErr
	}

	metaTags := getMetaTags()
	errorsContent := pages.ErrorsContent(err.Error(), link)

	templateHandler := templ.Handler(
		templates.Layout("Anti-Discover", metaTags, errorsContent),
	)

	return adaptor.HTTPHandler(templateHandler)(c)
}
