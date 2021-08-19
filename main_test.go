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
		err error
	)
	logger, hook = logtest.NewNullLogger()
	initLogger(logger)
	envy.Set(GO_ENV, "test")
	DB, err = pop.Connect("test")
	if err != nil {
		fmt.Println("Could not connect to test database")
		os.Exit(1)
	}
	code := m.Run()
	os.Exit(code)
}

func TestInitLogger(t *testing.T) {
	LOG.Error("test error")
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "test error", hook.LastEntry().Message)
	hook.Reset()
}

func TestInitEnvVariables(t *testing.T) {
	envy.Set(ENV_CLIENT_ID, "")
	envy.Set(ENV_CLIENT_SECRET, "")

	id, sec, err := initEnvVariables()
	assert.Equal(t, "", id)
	assert.Equal(t, "", sec)
	assert.Equal(t, fmt.Sprintf("Env key: %s not set", ENV_CLIENT_ID), err.Error())


	envy.Set(ENV_CLIENT_ID, "client_id123")
	id, sec, err = initEnvVariables()
	assert.Equal(t, "client_id123", id)
	assert.Equal(t, "", sec)
	assert.Equal(t, fmt.Sprintf("Env key: %s not set", ENV_CLIENT_SECRET), err.Error())


	envy.Set(ENV_CLIENT_SECRET, "client_secret123")
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
	assert.Contains(t, err.Error(), "Could not save token to file:")

	mock.LError = true
	err = loginAccount(mock)
	assert.Contains(t, err.Error(), "Could not get token:")
}

func TestStartApp(t *testing.T) {
	mock := spotifySaver.MockedSpotifySaver{LError: false}

	err := StartApp(&mock)
	assert.NoError(t, err)

	mock.LError = true
	err = StartApp(&mock)
	assert.Contains(t, err.Error(), "Could not load token:")
}

func TestStartSubCommands(t *testing.T) {
	mock := login.MockedAuth{
		LError: false,
		SError: false,
	}

	err := pop.DropDB(DB)
	assert.NoError(t, err)

	err, ready := StartSubCommands(nil, nil)
	assert.NoError(t, err)
	assert.True(t, ready)

	*loginFlag = true
	err, ready = StartSubCommands(nil, mock)
	assert.NoError(t, err)
	assert.False(t, ready)

	*createDb = true
	err, ready = StartSubCommands(DB, mock)
	assert.NoError(t, err)
	assert.False(t, ready)

	*createDb = false
	*migrate = true
	err, ready = StartSubCommands(DB, mock)
	assert.NoError(t, err)
	assert.False(t, ready)
}