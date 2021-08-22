package main

import (
	"flag"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/elivlo/SpotifyHistorySaver/login"
	"github.com/elivlo/SpotifyHistorySaver/models"
	"github.com/elivlo/SpotifyHistorySaver/spotifySaver"
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
	// EnvClientID is the env variable name for spotify client id
	EnvClientID = "CLIENT_ID"
	// EnvClientSecret is the env variable name for spotify client secret
	EnvClientSecret = "CLIENT_SECRET"
	// GoEnv is the env variable that defines in which stage the app is running
	// (development/production/test)
	GoEnv = "GO_ENV"

	// CallbackURI is the URL used to log in to the spotify account
	CallbackURI = "http://localhost:8080/callback"
)

var (
	log          *logrus.Entry
	env          string
	clientID     string
	clientSecret string
	createDb     = flag.Bool("create_db", false, "create_db: will create the database")
	migrate      = flag.Bool("migrate", false, "migrate: will migrate the current schema into db")
	loginFlag    = flag.Bool("login", false, "login: will get you an OAuth2 token for further usage")
)

// init logging
func initLogger(logger *logrus.Logger) {
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&nested.Formatter{
		FieldsOrder: []string{"component", "category"},
		HideKeys:    true,
	})
	log = logger.WithField("component", "SpotifyPlaybackSaver")
	log.Info("Setup SpotifyPlaybackSaver...")
}

// load env variables
func initEnvVariables() (string, string, error) {
	cID, err := envy.MustGet(EnvClientID)
	if err != nil || cID == "" {
		return cID, "", fmt.Errorf("env key: %s not set", EnvClientID)
	}
	cSec, err := envy.MustGet(EnvClientSecret)
	if err != nil || cSec == "" {
		return cID, cSec, fmt.Errorf("env key: %s not set", EnvClientSecret)
	}
	env = envy.Get(GoEnv, "development")

	return cID, cSec, nil
}

func createDB(c *pop.Connection) error {
	err := pop.CreateDB(c)
	if err != nil && strings.Contains(err.Error(), "database exists") {
		return fmt.Errorf("could not connect to database: %v", err)
	}
	return nil
}

func migrateDB(c *pop.Connection) error {
	box, err := pop.NewMigrationBox(packr.New("migrations", "./migrations"), c)
	if err != nil {
		return fmt.Errorf("could not load migrations: %s", err)
	}
	err = box.Up()
	if err != nil {
		return fmt.Errorf("could not migrate: %s", err)
	}
	return nil
}

func loginAccount(auth login.Auth) error {
	log.Info("Start login to your account...")
	token := auth.Login()

	if !token.Valid() {
		return fmt.Errorf("could not get valid token: %v", token)
	}

	err := auth.SaveToken(login.TokenFileName, token)
	if err != nil {
		return fmt.Errorf("could not save token to file: %v", err)
	}
	return nil
}

func startSubCommands(db *pop.Connection, auth login.Auth) (bool, error) {
	flag.Parse()

	if *createDb {
		return false, createDB(db)
	}

	if *migrate {
		return false, migrateDB(db)
	}

	if *loginFlag {
		return false, loginAccount(auth)
	}

	return true, nil
}

func startApp(s spotifySaver.InterfaceSpotifySaver) error {
	log.Info("Start listening to your spotify history...")
	var wg sync.WaitGroup

	err := s.LoadToken(spotifySaver.TokenFileName)
	if err != nil {
		return fmt.Errorf("could not load token: %v", err)
	}
	s.Authenticate(CallbackURI, clientID, clientSecret)

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
	log.Info("Shutting down...")

	return nil
}

func init() {
	initLogger(logrus.New())
}

func main() {
	var err error

	clientID, clientSecret, err = initEnvVariables()
	if err != nil {
		log.Fatal(err)
	}

	ready, err := startSubCommands(models.DB, login.NewLogin(CallbackURI, clientID, clientSecret))
	if err != nil {
		log.Error(err)
	}

	if !ready {
		return
	}

	s, err := spotifySaver.NewSpotifySaver(log, env)
	if err != nil {
		log.Fatal(err)
	}
	err = startApp(s)
	if err != nil {
		log.Fatal(err)
	}
}
