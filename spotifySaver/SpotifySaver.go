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
	"time"
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
func (s SpotifySaver) LoadToken() error {
	fileBytes, err := ioutil.ReadFile(tokenFileName)
	if err != nil {
		return err
	}

	var token oauth2.Token
	err = json.Unmarshal(fileBytes, token)
	if err != nil {
		return err
	}
	s.token = token

	if s.token.Expiry.Unix() <= time.Now().Unix() {
		return errors.New(fmt.Sprintf("token expired at %v", s.token.Expiry))
	}
	return nil
}

// Authenticate will create a new client from token.
func (s SpotifySaver) Authenticate() {
	s.client = s.auth.NewClient(&s.token)
}

// StartLastSongsWorker
func (s SpotifySaver) StartLastSongsWorker() error {
	return nil
}