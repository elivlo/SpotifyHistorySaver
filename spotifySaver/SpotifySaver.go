package spotifySaver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elivlo/SpotifyHistoryPlaybackSaver/login"
	"github.com/gobuffalo/pop/v5"
	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"io/ioutil"
	"sync"
	"time"
)

const (
	tokenFileName = "token.json"
)

var LOG *log.Entry
var ENV string

type InterfaceSpotifySaver interface {
	LoadToken() error
	Authenticate(callbackURI, clientID, clientSecret string)
	StartLastSongsWorker(wg *sync.WaitGroup, stop chan bool)
}

// SpotifySaver will handle all the saving logic.
type SpotifySaver struct {
	dbConnection *pop.Connection
	token        *oauth2.Token
	auth         *spotifyauth.Authenticator
	client       *spotify.Client
}

// NewSpotifySaver will create a new SpotifySaver instance with database connection.
// It will throw an error when database connection fails.
func NewSpotifySaver(log *log.Entry, env string) (SpotifySaver, error) {
	LOG = log
	ENV = env
	tx, err := pop.Connect(ENV)
	if err != nil {
		return SpotifySaver{}, errors.New(fmt.Sprintf("Could not connect to database: %v", err))
	}
	return SpotifySaver{
		dbConnection: tx,
	}, nil
}

// LoadToken will load the token from file "token.json" in exec directory.
// It will throw an error when the token is expired.
func (s *SpotifySaver) LoadToken() error {
	fileBytes, err := ioutil.ReadFile(tokenFileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileBytes, &s.token)
	if err != nil {
		return err
	}

	if !s.token.Valid() && s.token.RefreshToken == "" {
		return errors.New(fmt.Sprintf("token expired at %v", s.token.Expiry))
	}
	return nil
}

// Authenticate will create a new client from token.
func (s *SpotifySaver) Authenticate(callbackURI, clientID, clientSecret string) {
	s.auth = spotifyauth.New(spotifyauth.WithRedirectURL(callbackURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadRecentlyPlayed),
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithClientSecret(clientSecret))
	s.client = spotify.New(s.auth.Client(context.Background(), s.token))
}

// StartLastSongsWorker is a worker that will send history requests every 45 minutes.
// It is not async. It accepts a wait group and will send Done when stopped. It may be stopped with stop chan value.
func (s *SpotifySaver) StartLastSongsWorker(wg *sync.WaitGroup, stop chan bool) {
	first := true
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			LOG.Info("Fetch newly listened songs")
			last, err := getLastHistoryEntry(s.dbConnection)
			if err != nil {
				LOG.Warn("Could not get last played song: ", err)
				last.PlayedAt = time.Unix(0, 0)
			}

			songs, err := s.client.PlayerRecentlyPlayedOpt(context.Background(), &spotify.RecentlyPlayedOptions{
				Limit:        50,
				AfterEpochMs: last.PlayedAt.Unix()*1000 + 1,
			})
			if err != nil {
				LOG.Error("Could not get recently played songs: ", err)
			}

			fetched := CreateFetchedSongs(s.dbConnection, songs)
			err = fetched.TransformAndInsertIntoDatabase()
			if err != nil {
				LOG.Error("Could not save recently played songs: ", err)
			}
			LOG.Info("Finished fetching newly listened songs")
			if first {
				first = false
				ticker.Reset(time.Minute * 45)
			}

			token, err := s.client.Token()
			if err != nil {
				LOG.Error("Could not get current client token: ", err)
			}
			err = login.NewLogin("").SaveToken(login.TokenFileName, token)
			if err != nil {
				LOG.Error("Could not save current client token ", err)
			}
		case <-stop:
			LOG.Info("Shutting down StartLastSongsWorker")
			ticker.Stop()
			wg.Done()
			return
		}
	}
}
