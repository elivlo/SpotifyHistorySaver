package main

import (
	"errors"
	"flag"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/elivlo/SpotifyHistoryPlaybackSaver/login"
	"github.com/elivlo/SpotifyHistoryPlaybackSaver/models"
	"github.com/elivlo/SpotifyHistoryPlaybackSaver/spotifySaver"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/pop/v5"
	"github.com/sirupsen/logrus"
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
	LOG          *logrus.Entry
	ENV          string
	ClientId     string
	ClientSecret string
	createDb = flag.Bool("create_db", false, "create_db: will create the database")
	migrate = flag.Bool("migrate", false, "migrate: will migrate the current schema into db")
	loginFlag = flag.Bool("login", false, "login: will get you an OAuth2 token for further usage")
)

// init logging
func initLogger(logger *logrus.Logger) {
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&nested.Formatter{
		FieldsOrder: []string{"component", "category"},
		HideKeys:    true,
	})
	LOG = logger.WithField("component", "SpotifyPlaybackSaver")
	LOG.Info("Setup SpotifyPlaybackSaver...")
}

// load env variables
func initEnvVariables() (string, string, error) {
	cId, err := envy.MustGet(ENV_CLIENT_ID)
	if err != nil || cId == "" {
		return cId, "", errors.New(fmt.Sprintf("Env key: %s not set", ENV_CLIENT_ID))
	}
	cSec, err := envy.MustGet(ENV_CLIENT_SECRET)
	if err != nil || cSec == "" {
		return cId, cSec, errors.New(fmt.Sprintf("Env key: %s not set", ENV_CLIENT_SECRET))
	}
	ENV = envy.Get(GO_ENV, "development")

	return cId, cSec, nil
}

func createDB(c *pop.Connection) error {
	err := pop.CreateDB(c)
	if err != nil && strings.Contains(err.Error(), "database exists") {
		return errors.New(fmt.Sprintf("Could not connect to database: %v", err))
	}
	return nil
}

func migrateDB(c *pop.Connection) error {
	box, err := pop.NewMigrationBox(packr.New("migrations", "./migrations"), c)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not load migrations: %s", err))
	}
	err = box.Up()
	if err != nil {
		return errors.New(fmt.Sprintf("Could not migrate: %s", err))
	}
	return nil
}

func loginAccount(auth login.Auth) error {
	LOG.Info("Start login to your account...")
	token, err := auth.Login(ClientId, ClientSecret)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not get token: %v", err))
	}

	err = auth.SaveToken(login.TokenFileName, token)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not save token to file: %v", err))
	}
	return nil
}

func StartSubCommands(db *pop.Connection, auth login.Auth) error {
	flag.Parse()

	if *createDb {
		return createDB(db)
	}

	if *migrate {
		return migrateDB(db)
	}

	if *loginFlag {
		return loginAccount(auth)
	}

	return nil
}

func StartApp(s spotifySaver.InterfaceSpotifySaver) error {
	LOG.Info("Start listening to your spotify history...")
	var wg sync.WaitGroup

	err := s.LoadToken(spotifySaver.TokenFileName)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not load token: %v", err))
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

	return nil
}

func init() {
	var err error

	initLogger(logrus.New())

	ClientId, ClientSecret, err = initEnvVariables()
	if err != nil {
		LOG.Fatal(err)
	}
}

func main() {
	err := StartSubCommands(models.DB, login.NewLogin(CallbackURI))
	if err != nil {
		LOG.Error(err)
	}

	s, err := spotifySaver.NewSpotifySaver(LOG, ENV)
	if err != nil {
		LOG.Fatal(err)
	}
	err = StartApp(s)
	if err != nil {
		LOG.Fatal(err)
	}
}
