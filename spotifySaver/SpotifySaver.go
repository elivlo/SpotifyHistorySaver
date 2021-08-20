package spotifySaver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elivlo/SpotifyHistorySaver/login"
	"github.com/elivlo/SpotifyHistorySaver/models"
	"github.com/gobuffalo/pop/v5"
	"github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"io/ioutil"
	"sync"
	"time"
)

const (
	// TokenFileName is the standard file name to save the OAuth token to
	TokenFileName = "token.json"
)

// InterfaceSpotifySaver is the interface SpotifySaver implements.
// It supports loading a token and authenticating with it.
// The main purpose is to start the StartLastSongsWorker to periodically save the history.
type InterfaceSpotifySaver interface {
	LoadToken(file string) error
	Authenticate(callbackURI, clientID, clientSecret string)
	StartLastSongsWorker(wg *sync.WaitGroup, stop chan bool)
}

// SpotifySaver will handle all the saving logic.
// It supports loading a token and authenticating with it.
// The main purpose is to start the StartLastSongsWorker to periodically save the history.
type SpotifySaver struct {
	dbConnection *pop.Connection
	token        *oauth2.Token
	auth         *spotifyauth.Authenticator
	client       *spotify.Client
	log          *logrus.Entry
	env          string
}

// NewSpotifySaver will create a new SpotifySaver instance with database connection.
// It will throw an error when database connection fails.
func NewSpotifySaver(log *logrus.Entry, env string) (*SpotifySaver, error) {
	tx, err := pop.Connect(env)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to database: %v", err)
	}
	return &SpotifySaver{
		dbConnection: tx,
		log:          log,
		env:          env,
	}, nil
}

// LoadToken will load the token from file "token.json" in exec directory.
// It will throw an error when the token is expired.
func (s *SpotifySaver) LoadToken(file string) error {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileBytes, &s.token)
	if err != nil {
		return err
	}

	if !s.token.Valid() && s.token.RefreshToken == "" {
		return fmt.Errorf("token expired at %v", s.token.Expiry)
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
			s.log.Info("Fetch newly listened songs")

			last := s.getLastEntry()

			songs := s.fetchNewSongs(last)

			s.insertNewSongs(songs)

			s.log.Info("Finished fetching newly listened songs")

			s.saveNewToken(login.TokenFileName)

			if first {
				first = false
				ticker.Reset(time.Minute * 45)
			}
		case <-stop:
			s.log.Info("Shutting down StartLastSongsWorker")
			ticker.Stop()
			wg.Done()
			return
		}
	}
}

func (s *SpotifySaver) getLastEntry() models.HistoryEntry {
	last, err := getLastHistoryEntry(s.dbConnection)
	if err != nil {
		s.log.Warnf("Could not get last played song: %v", err)
		last.PlayedAt = time.Unix(0, 0)
	}
	return last
}

func (s *SpotifySaver) fetchNewSongs(last models.HistoryEntry) []spotify.RecentlyPlayedItem {
	songs, err := s.client.PlayerRecentlyPlayedOpt(context.Background(), &spotify.RecentlyPlayedOptions{
		Limit:        50,
		AfterEpochMs: last.PlayedAt.Unix()*1000 + 1000,
	})
	if err != nil {
		s.log.Error("Could not get recently played songs: ", err)
	}

	return songs
}

func (s *SpotifySaver) insertNewSongs(songs []spotify.RecentlyPlayedItem) {
	fetched := NewFetchedSongs(s.dbConnection, songs)
	err := fetched.TransformAndInsertIntoDatabase(s.log)
	if err != nil {
		s.log.Error("Could not save recently played songs: ", err)
	}
}

func (s *SpotifySaver) saveNewToken(fileName string) {
	token, err := s.client.Token()
	if err != nil {
		s.log.Error("Could not get current client token: ", err)
	}
	err = login.NewLogin("").SaveToken(fileName, token)
	if err != nil {
		s.log.Error("Could not save current client token ", err)
	}
}
