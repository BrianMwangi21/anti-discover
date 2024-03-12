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

func getRecommendationAndCreatePlaylist(client *spotify.Client, userID string) ([]spotify.SimpleTrack, *spotify.FullPlaylist, error) {
	ctx := context.Background()
	topTracks, err := client.CurrentUsersTopTracks(ctx, spotify.Timerange(spotify.LongTermRange), spotify.Limit(50))

	if err != nil {
		return nil, nil, err
	}

	var trackIDs []spotify.ID
	for _, track := range topTracks.Tracks {
		trackIDs = append(trackIDs, track.ID)
	}

	audioFeatures, err := client.GetAudioFeatures(ctx, trackIDs...)
	if err != nil {
		return nil, nil, err
	}

	seedTrackIds := []spotify.ID{
		topTracks.Tracks[9].ID,
		topTracks.Tracks[19].ID,
		topTracks.Tracks[29].ID,
		topTracks.Tracks[39].ID,
		topTracks.Tracks[49].ID,
	}

	seeds := spotify.Seeds{
		Tracks: seedTrackIds,
	}
	trackAttributes := calculateAntiFeatures(audioFeatures)

	recommendations, err := client.GetRecommendations(ctx, seeds, trackAttributes)
	if err != nil {
		return nil, nil, err
	}

	playlist, err := createPlaylist(client, recommendations.Tracks, userID)
	if err != nil {
		return nil, nil, err
	}

	return recommendations.Tracks, playlist, nil
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

func createPlaylist(client *spotify.Client, recommendations []spotify.SimpleTrack, userID string) (*spotify.FullPlaylist, error) {
	ctx := context.Background()

	playlist, err := client.CreatePlaylistForUser(ctx, userID, "Anti-Discover", "Playlist created from Spotify's evil twin - Anti-Discover", true, false)
	if err != nil {
		return nil, err
	}

	var trackIDs []spotify.ID
	for _, recommendation := range recommendations {
		trackIDs = append(trackIDs, recommendation.ID)
	}

	_, err = client.AddTracksToPlaylist(ctx, playlist.ID, trackIDs...)
	if err != nil {
		return nil, err
	}

	return playlist, nil
}
