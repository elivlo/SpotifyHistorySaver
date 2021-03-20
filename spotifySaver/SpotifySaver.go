package spotifySaver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gobuffalo/envy"
	pop "github.com/gobuffalo/pop/v5"
	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"io/ioutil"
	"sync"
)

const (
	GO_ENV = "GO_ENV"
	tokenFileName = "token.json"
)

var logger *log.Entry
var env string

// SpotifySaver will handle all the saving logic.
type SpotifySaver struct {
	dbConnection *pop.Connection
	token oauth2.Token
	auth spotify.Authenticator
	client spotify.Client
}

// NewSpotifySaver will create a new SpotifySaver instance with database connection.
// It will throw an error when database connection fails.
func NewSpotifySaver(log *log.Entry) (SpotifySaver, error) {
	logger = log
	env = envy.Get(GO_ENV, "development")
	tx, err := pop.Connect(env)
	if err != nil {
		return SpotifySaver{}, err
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
	s.auth = spotify.NewAuthenticator(callbackURI, spotify.ScopeUserReadRecentlyPlayed)
	s.auth.SetAuthInfo(clientID, clientSecret)
	s.client = s.auth.NewClient(&s.token)
}

// StartLastSongsWorker
func (s *SpotifySaver) StartLastSongsWorker(wg *sync.WaitGroup) {
	songs, err := s.client.PlayerRecentlyPlayedOpt(&spotify.RecentlyPlayedOptions{
		Limit: 1,
	})
	if err != nil {
		logger.Fatal(err)
	}
	if len(songs) > 0 {
		fmt.Println(songs[0])
	}
	wg.Done()
}