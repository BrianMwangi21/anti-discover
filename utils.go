package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/BrianMwangi21/anti-discover.git/templates/pages"
	"github.com/a-h/templ"
	gowebly "github.com/gowebly/helpers"
	"github.com/valyala/fasthttp"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
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
	redirectURI := gowebly.Getenv("REDIRECT_URI", "")
	if redirectURI == "" {
		redirectURI = "http://localhost:8080/anti"
	}
	return spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(scopes[:]...))
}

func getMetaTags() templ.Component {
	return pages.MetaTags(
		"Anti-Discover",
		"Spotify's discover weekly rogue twin",
	)
}

func saveToken(request *http.Request, token *oauth2.Token) error {
	ctx := context.Background()
	values := request.URL.Query()
	code := values.Get("code")

	tokenData, err := json.Marshal(token)
	if err != nil {
		return errors.New("Error marshalling struct to JSON")
	}

	tokenString := string(tokenData)

	err = redisClient.Set(ctx, code, tokenString, time.Hour).Err()
	if err != nil {
		return errors.New("Error saving token to Redis")
	}

	return nil
}

func retrieveToken(request *http.Request) (*oauth2.Token, error) {
	ctx := context.Background()
	values := request.URL.Query()
	code := values.Get("code")

	val, err := redisClient.Get(ctx, code).Result()
	if err != nil {
		return nil, errors.New("Error getting token from Redis")
	}

	var token oauth2.Token
	err = json.Unmarshal([]byte(val), &token)
	if err != nil {
		return nil, errors.New("Error unmarshalling to struct")
	}

	return &token, nil
}

func getSpotifyLink() (templ.SafeURL, error) {
	auth := getAuth()

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

func getRecommendationAndCreatePlaylist(client *spotify.Client, userID string) ([]spotify.SimpleTrack, *spotify.FullPlaylist, error) {
	ctx := context.Background()
	topTracks, err := client.CurrentUsersTopTracks(ctx)

	if err != nil {
		return nil, nil, err
	}

	rand.NewSource(time.Now().UnixNano())
	randomIndex := rand.Intn(len(topTracks.Tracks))

	seeds := spotify.Seeds{
		Tracks: []spotify.ID{
			topTracks.Tracks[randomIndex].ID,
		},
	}

	var trackIDs []spotify.ID
	for _, track := range topTracks.Tracks {
		trackIDs = append(trackIDs, track.ID)
	}

	audioFeatures, err := client.GetAudioFeatures(ctx, trackIDs...)
	if err != nil {
		return nil, nil, err
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
		TargetAcousticness(superSecretTakeItToExtremeAlgo(acousticness, count)).
		TargetDanceability(superSecretTakeItToExtremeAlgo(danceability, count)).
		TargetEnergy(superSecretTakeItToExtremeAlgo(energy, count)).
		TargetInstrumentalness(superSecretTakeItToExtremeAlgo(instrumentalness, count)).
		TargetLiveness(superSecretTakeItToExtremeAlgo(liveness, count)).
		TargetValence(superSecretTakeItToExtremeAlgo(valence, count))
}

func superSecretTakeItToExtremeAlgo(feature, count float64) float64 {
	target := 1 - feature/count

	if target > 0.5 {
		return 0.9
	}

	return 0.1
}

func createPlaylist(client *spotify.Client, recommendations []spotify.SimpleTrack, userID string) (*spotify.FullPlaylist, error) {
	ctx := context.Background()

	playlists, err := client.GetPlaylistsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, playlist := range playlists.Playlists {
		if playlist.Name == "Anti-Discover" {
			err = client.UnfollowPlaylist(ctx, playlist.ID)
			if err != nil {
				return nil, err
			}
			break
		}
	}

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
