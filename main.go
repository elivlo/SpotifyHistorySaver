package main

import (
	"flag"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/elivlo/SpotifyHistoryPlaybackSaver/login"
	"github.com/elivlo/SpotifyHistoryPlaybackSaver/models"
	"github.com/elivlo/SpotifyHistoryPlaybackSaver/spotifySaver"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/pop/v5"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"sync"
)

const (
	ENV_CLIENT_ID     = "CLIENT_ID"
	ENV_CLIENT_SECRET = "CLIENT_SECRET"
	GO_ENV            = "GO_ENV"

	CallbackURI = "http://localhost:8080/callback"
)

var (
	LOG          *log.Entry
	ENV          string
	ClientId     string
	ClientSecret string
)

// init logging
func initLogger() {
	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	logger.SetFormatter(&nested.Formatter{
		FieldsOrder: []string{"component", "category"},
		HideKeys:    true,
	})
	LOG = logger.WithField("component", "SpotifyPlaybackSaver")
	LOG.Info("Setup SpotifyPlaybackSaver...")
}

// load env variables
func initEnvVariables() {
	var err error
	ClientId, err = envy.MustGet(ENV_CLIENT_ID)
	if err != nil {
		LOG.Fatalf("Env key: %s not set", ENV_CLIENT_ID)
	}
	ClientSecret, err = envy.MustGet(ENV_CLIENT_SECRET)
	if err != nil {
		LOG.Fatalf("Env key: %s not set", ENV_CLIENT_SECRET)
	}
	ENV = envy.Get(GO_ENV, "development")
}

func StartApp() {
	loginFlag := flag.Bool("login", false, "login: will get you an OAuth2 token for further usage")
	migrate := flag.Bool("migrate", false, "login: will get you an OAuth2 token for further usage")
	createDb := flag.Bool("create_db", false, "login: will create the database")
	flag.Parse()

	if *createDb {
		err := pop.CreateDB(models.DB)
		if err != nil && strings.Contains(err.Error(), "database exists") {
			LOG.Fatalf("Could not connect to database: %v", err)
		}
		return
	}

	if *migrate {
		box, err := pop.NewMigrationBox(packr.New("migrations", "./migrations"), models.DB)
		if err != nil {
			LOG.Fatalf("Could not load migrations: %s", err)
		}
		err = box.Up()
		if err != nil {
			LOG.Fatalf("Could not migrate: %s", err)
		}
		return
	}

	if *loginFlag {
		LOG.Info("Start login to your account...")
		token, err := login.Login(ClientId, ClientSecret, CallbackURI)
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
	var wg sync.WaitGroup

	s, err := spotifySaver.NewSpotifySaver(LOG, ENV)
	if err != nil {
		LOG.Fatalf("Could not connect to database: %v", err)
	}

	err = s.LoadToken()
	if err != nil {
		LOG.Fatalf("Could not load token: %v", err)
	}
	s.Authenticate(CallbackURI, ClientId, ClientSecret)

	stop := make(chan bool, 1)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			stop <- true
		}
	}()

	wg.Add(1)
	go s.StartLastSongsWorker(&wg, stop)

	wg.Wait()
	LOG.Info("Shutting down...")
}

func init() {
	initLogger()
	initEnvVariables()
}

func main() {
	StartApp()
}
