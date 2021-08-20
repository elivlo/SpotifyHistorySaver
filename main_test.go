package main

import (
	"fmt"
	"github.com/elivlo/SpotifyHistoryPlaybackSaver/login"
	"github.com/elivlo/SpotifyHistoryPlaybackSaver/spotifySaver"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v5"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var hook *logtest.Hook
var DB *pop.Connection

func TestMain(m *testing.M) {
	var (
		logger *logrus.Logger
		err    error
	)
	logger, hook = logtest.NewNullLogger()
	initLogger(logger)
	envy.Set(GoEnv, "test")
	DB, err = pop.Connect("test")
	if err != nil {
		fmt.Println("Could not connect to test database")
		os.Exit(1)
	}
	code := m.Run()
	os.Exit(code)
}

func TestInitLogger(t *testing.T) {
	log.Error("test error")
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "test error", hook.LastEntry().Message)
	hook.Reset()
}

func TestInitEnvVariables(t *testing.T) {
	envy.Set(EnvClientID, "")
	envy.Set(EnvClientSecret, "")

	id, sec, err := initEnvVariables()
	assert.Equal(t, "", id)
	assert.Equal(t, "", sec)
	assert.Equal(t, fmt.Sprintf("env key: %s not set", EnvClientID), err.Error())

	envy.Set(EnvClientID, "client_id123")
	id, sec, err = initEnvVariables()
	assert.Equal(t, "client_id123", id)
	assert.Equal(t, "", sec)
	assert.Equal(t, fmt.Sprintf("env key: %s not set", EnvClientSecret), err.Error())

	envy.Set(EnvClientSecret, "client_secret123")
	id, sec, err = initEnvVariables()
	assert.Equal(t, "client_id123", id)
	assert.Equal(t, "client_secret123", sec)
}

func TestCreateDB(t *testing.T) {
	err := pop.DropDB(DB)
	assert.NoError(t, err)

	err = createDB(DB)
	assert.NoError(t, err)
}

func TestMigrateDB(t *testing.T) {
	err := migrateDB(DB)
	assert.NoError(t, err)
}

func TestLogin(t *testing.T) {
	mock := login.MockedAuth{
		LError: false,
		SError: false,
	}

	err := loginAccount(mock)
	assert.NoError(t, err)

	mock.SError = true
	err = loginAccount(mock)
	assert.Contains(t, err.Error(), "could not save token to file:")

	mock.LError = true
	err = loginAccount(mock)
	assert.Contains(t, err.Error(), "could not get token:")
}

func TestStartApp(t *testing.T) {
	mock := spotifySaver.MockedSpotifySaver{LError: false}

	err := startApp(&mock)
	assert.NoError(t, err)

	mock.LError = true
	err = startApp(&mock)
	assert.Contains(t, err.Error(), "could not load token:")
}

func TestStartSubCommands(t *testing.T) {
	mock := login.MockedAuth{
		LError: false,
		SError: false,
	}

	err := pop.DropDB(DB)
	assert.NoError(t, err)

	ready, err := startSubCommands(nil, nil)
	assert.NoError(t, err)
	assert.True(t, ready)

	*loginFlag = true
	ready, err = startSubCommands(nil, mock)
	assert.NoError(t, err)
	assert.False(t, ready)

	*createDb = true
	ready, err = startSubCommands(DB, mock)
	assert.NoError(t, err)
	assert.False(t, ready)

	*createDb = false
	*migrate = true
	ready, err = startSubCommands(DB, mock)
	assert.NoError(t, err)
	assert.False(t, ready)
}
