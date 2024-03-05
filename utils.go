package main

import (
	"bytes"
	"net/http"

	"github.com/BrianMwangi21/anti-discover.git/templates/pages"
	"github.com/a-h/templ"
	"github.com/valyala/fasthttp"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var (
	scopes = [...]string{
		spotifyauth.ScopePlaylistReadPrivate, spotifyauth.ScopePlaylistModifyPublic, spotifyauth.ScopePlaylistModifyPrivate,
		spotifyauth.ScopePlaylistReadCollaborative, spotifyauth.ScopeUserReadEmail, spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadRecentlyPlayed, spotifyauth.ScopeUserTopRead,
	}
)

func convertRequest(req *fasthttp.Request) (*http.Request, error) {
	httpReq, err := http.NewRequest(
		string(req.Header.Method()),
		string(req.Header.RequestURI()),
		bytes.NewReader(req.Body()),
	)

	if err != nil {
		return nil, err
	}

	req.Header.VisitAll(func(key, value []byte) {
		httpReq.Header.Set(string(key), string(value))
	})

	return httpReq, nil
}

func getAuth() *spotifyauth.Authenticator {
	return spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(scopes[:]...))
}

func getMetaTags() templ.Component {
	return pages.MetaTags(
		"Anti-Discover",
		"Spotify's discover weekly rogue twin",
	)
}
