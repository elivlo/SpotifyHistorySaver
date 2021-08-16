package main

import (
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"testing"
)

var hook *logtest.Hook

func Test_Setup(t *testing.T) {
	var logger *logrus.Logger
	logger, hook = logtest.NewNullLogger()
	LOG = logger.WithField("component", "SpotifyPlaybackSaver")
}
