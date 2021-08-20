package spotifySaver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elivlo/SpotifyHistoryPlaybackSaver/models"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"os"
	"testing"
	"time"
)

var DB *pop.Connection

func getLogger() (*logtest.Hook, *logrus.Entry) {
	logger, hook := logtest.NewNullLogger()
	logger.Level = logrus.DebugLevel
	log := logger.WithField("test", "test")

	return hook, log
}

func TestMain(m *testing.M) {
	var err error

	envy.Set("GO_ENV", "test")
	DB, err = pop.Connect("test")
	if err != nil {
		fmt.Println("Could not connect to test database")
		os.Exit(1)
	}
	_ = DB.TruncateAll()

	code := m.Run()
	os.Exit(code)
}

func TestNewSpotifySaver(t *testing.T) {
	_, log := getLogger()

	saver, err := NewSpotifySaver(log, "test")
	assert.NoError(t, err)
	assert.Equal(t, "test", saver.env)

	saver, err = NewSpotifySaver(log, "invalid")
	assert.Error(t, err)
	assert.Nil(t, saver)
}

func TestSpotifySaver_LoadToken(t *testing.T) {
	_, log := getLogger()

	saver, err := NewSpotifySaver(log, "test")
	assert.NoError(t, err)

	tokenName, err := uuid.NewV4()
	assert.NoError(t, err)

	t.Run("NoFile", func(t *testing.T) {
		err = saver.LoadToken(tokenName.String())
		assert.Error(t, err)
	})

	var file *os.File
	t.Run("FileEmpty", func(t *testing.T) {
		file, err = os.Create(tokenName.String())
		assert.NoError(t, err)

		err = saver.LoadToken(tokenName.String())
		assert.Error(t, err)
	})

	t.Run("TokenInvalid", func(t *testing.T) {
		tokenBytes, err := json.Marshal(oauth2.Token{})
		assert.NoError(t, err)

		_, err = file.Write(tokenBytes)
		assert.NoError(t, err)

		err = saver.LoadToken(tokenName.String())
		assert.Error(t, err)
	})

	t.Run("Valid", func(t *testing.T) {
		tokenBytes, err := json.Marshal(&oauth2.Token{
			AccessToken:  "aaaa",
			TokenType:    "tttt",
			RefreshToken: "rrrr",
			Expiry:       time.Now().Add(time.Hour),
		})
		assert.NoError(t, err)

		_, err = file.WriteAt(tokenBytes, 0)
		assert.NoError(t, err)

		err = file.Close()
		assert.NoError(t, err)

		err = saver.LoadToken(tokenName.String())
		assert.Nil(t, err)
	})

	err = os.Remove(tokenName.String())
	assert.NoError(t, err)
}

func TestSpotifySaver_Authenticate(t *testing.T) {
	_, log := getLogger()

	saver, err := NewSpotifySaver(log, "test")
	assert.NoError(t, err)

	saver.Authenticate("url", "id", "secret")
}

func TestSpotifySaver_getLastEntry(t *testing.T) {
	_, log := getLogger()

	saver, err := NewSpotifySaver(log, "test")
	assert.NoError(t, err)

	entry := saver.getLastEntry()
	assert.Equal(t, time.Unix(0, 0), entry.PlayedAt)
}

func TestSpotifySaver_fetchNewSongs(t *testing.T) {
	_, log := getLogger()

	saver, err := NewSpotifySaver(log, "test")
	assert.NoError(t, err)

	saver.auth = spotifyauth.New()
	saver.client = spotify.New(saver.auth.Client(context.Background(), &oauth2.Token{}))

	items := saver.fetchNewSongs(models.HistoryEntry{
		PlayedAt: time.Unix(0, 0),
	})
	assert.Equal(t, 0, len(items))
}

func TestSpotifySaver_InsertNewSongs(t *testing.T) {
	_, log := getLogger()

	saver, err := NewSpotifySaver(log, "test")
	assert.NoError(t, err)

	saver.insertNewSongs([]spotify.RecentlyPlayedItem{{
		Track:           spotify.SimpleTrack{},
		PlayedAt:        time.Now(),
		PlaybackContext: spotify.PlaybackContext{},
	}})
}

func TestSpotifySaver_SaveNewToken(t *testing.T) {
	_, log := getLogger()

	saver, err := NewSpotifySaver(log, "test")
	assert.NoError(t, err)

	saver.auth = spotifyauth.New()
	saver.client = spotify.New(saver.auth.Client(context.Background(), &oauth2.Token{}))

	tokenName, err := uuid.NewV4()
	assert.NoError(t, err)

	saver.saveNewToken(tokenName.String())

	err = os.Remove(tokenName.String())
	assert.NoError(t, err)
}
