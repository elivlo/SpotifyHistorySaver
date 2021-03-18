package main

import (
	"flag"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/elivlo/SpotifyHistoryPlaybackSaver/Login"
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
	login := flag.Bool("login", false, "login: will get you an OAuth2 token for further usage")

	if *login {
		LOG.Info("Start login to your account...")
		token, err := Login.Login(ClientId, ClientSecret, "http://localhost:8080/callback")
		if err != nil {
			LOG.Fatal(err)
		}

		fmt.Println(token.AccessToken)
		fmt.Println(token.Expiry)
		fmt.Println(token.RefreshToken)
		fmt.Println(token.TokenType)

		return
	}

	LOG.Info("Start listening to your spotify history...")
	//SpotifySaver.NewSpotifySaver(LOG)

	//c.Authenticate(token.AccessToken, token.TokenType, token.RefreshToken, token.Expiry)

}
