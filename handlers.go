package main

import (
	"errors"
	"net/url"

	"github.com/BrianMwangi21/anti-discover.git/templates"
	"github.com/BrianMwangi21/anti-discover.git/templates/pages"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	gowebly "github.com/gowebly/helpers"
	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"

	"github.com/gofiber/fiber/v2"
)

const redirectURI = "http://localhost:7000/anti"

var (
	scopes = [...]string{
		spotifyauth.ScopePlaylistReadPrivate, spotifyauth.ScopePlaylistModifyPublic, spotifyauth.ScopePlaylistModifyPrivate,
		spotifyauth.ScopePlaylistReadCollaborative, spotifyauth.ScopeUserReadEmail, spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadRecentlyPlayed, spotifyauth.ScopeUserTopRead,
	}
	auth  = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(scopes[:]...))
	ch    = make(chan *spotify.Client)
	state = "anti-discover"
)

func getSpotifyLink() (templ.SafeURL, error) {
	spotifyID := gowebly.Getenv("SPOTIFY_ID", "")
	if spotifyID == "" {
		return "", errors.New("SPOTIFY_ID not set")
	}

	authURL := auth.AuthURL(state)
	parsedURL, err := url.Parse(authURL)

	if err != nil {
		return "", errors.New("Parsing error failed")
	}

	query := parsedURL.Query()
	query.Set("client_id", spotifyID)
	parsedURL.RawQuery = query.Encode()
	updatedURL := parsedURL.String()

	return templ.URL(updatedURL), nil
}

func indexViewHandler(c *fiber.Ctx) error {
	link, err := getSpotifyLink()

	if err != nil {
		return err
	}

	metaTags := pages.MetaTags(
		"Anti-Discover",
		"Spotify's discover weekly rogue twin",
	)
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

func connectToSpotifyHandler(c *fiber.Ctx) error {
	if c.Get("HX-Request") == "" || c.Get("HX-Request") != "true" {
		return fiber.NewError(fiber.StatusBadRequest, "non-htmx request")
	}

	url := auth.AuthURL(state)
	return c.SendString("<p>🎉 To connect to Spotify, follow the following link: " + url + "</p>")
	// return c.SendString("<p>🎉 Yes! We connected to Spotify!</p>")
}
