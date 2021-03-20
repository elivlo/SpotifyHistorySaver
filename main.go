package main

import (
	"flag"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/elivlo/SpotifyHistoryPlaybackSaver/login"
	"github.com/elivlo/SpotifyHistoryPlaybackSaver/spotifySaver"
	"github.com/gobuffalo/envy"
	log "github.com/sirupsen/logrus"
)

const (
	ENV_CLIENT_ID = "CLIENT_ID"
	ENV_CLIENT_SECRET = "CLIENT_SECRET"
)

var (
	LOG          *log.Entry
	ClientId     string
	ClientSecret string
)

// init logging
func init() {
	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	logger.SetFormatter(&nested.Formatter{
		FieldsOrder:   []string{"component", "category"},
		HideKeys:      true,
	})
	LOG = logger.WithField("component", "SpotifyPlaybackSaver")
	LOG.Info("Setup SpotifyPlaybackSaver...")
}

// load env variables
func init() {
	var err error
	ClientId, err = envy.MustGet(ENV_CLIENT_ID)
	if err != nil {
		LOG.Fatalf("Env key: %s not set", ENV_CLIENT_ID)
	}
	ClientSecret, err = envy.MustGet(ENV_CLIENT_SECRET)
	if err != nil {
		LOG.Fatalf("Env key: %s not set", ENV_CLIENT_SECRET)
	}
}

func main() {
	loginFlag := flag.Bool("login", false, "login: will get you an OAuth2 token for further usage")

	if *loginFlag {
		LOG.Info("Start login to your account...")
		token, err := login.Login(ClientId, ClientSecret, "http://localhost:8080/callback")
		if err != nil {
			LOG.Fatalf("Could not get token: %v", err)
		}

		err = login.SaveToken(token)
		if err != nil {
			LOG.Fatalf("Could not save token to file: %v", err)
		}
		return
	}

	LOG.Info("Start listening to your spotify history...")
	s, err := spotifySaver.NewSpotifySaver(LOG)
	if err != nil {
		LOG.Fatalf("Could not connect to database: %v", err)
	}

	err = s.LoadToken()
	if err != nil {
		LOG.Fatalf("Could not load token: %v", err)
	}
	s.Authenticate()

	

}
