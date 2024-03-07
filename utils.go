package main

import (
	"bytes"
	"context"
	"net/http"

	"github.com/BrianMwangi21/anti-discover.git/templates/pages"
	"github.com/a-h/templ"
	"github.com/valyala/fasthttp"
	"github.com/zmb3/spotify/v2"
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

func getMusicRecommendations(client *spotify.Client) ([]spotify.SimpleTrack, error) {
	ctx := context.Background()
	topTracks, err := client.CurrentUsersTopTracks(ctx)
	if err != nil {
		return nil, err
	}

	var trackIDs []spotify.ID
	for _, track := range topTracks.Tracks {
		trackIDs = append(trackIDs, track.ID)
	}

	audioFeatures, err := client.GetAudioFeatures(ctx, trackIDs...)
	if err != nil {
		return nil, err
	}

	seeds := spotify.Seeds{
		Tracks: []spotify.ID{topTracks.Tracks[0].ID},
	}
	trackAttributes := calculateAntiFeatures(audioFeatures)

	recommendations, err := client.GetRecommendations(ctx, seeds, trackAttributes)
	if err != nil {
		return nil, err
	}

	return recommendations.Tracks, nil
}

func calculateAntiFeatures(features []*spotify.AudioFeatures) *spotify.TrackAttributes {
	var (
		acousticness     float64
		danceability     float64
		energy           float64
		instrumentalness float64
		liveness         float64
		valence          float64
	)

	for _, f := range features {
		acousticness += float64(f.Acousticness)
		danceability += float64(f.Danceability)
		energy += float64(f.Energy)
		instrumentalness += float64(f.Instrumentalness)
		liveness += float64(f.Liveness)
		valence += float64(f.Valence)
	}

	count := float64(len(features))
	return spotify.NewTrackAttributes().
		TargetAcousticness(1 - acousticness/count).
		TargetDanceability(1 - danceability/count).
		TargetEnergy(1 - energy/count).
		TargetInstrumentalness(1 - instrumentalness/count).
		TargetLiveness(1 - liveness/count).
		TargetValence(1 - valence/count)
}
