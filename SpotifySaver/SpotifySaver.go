package SpotifySaver

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"io/ioutil"
	"time"
)

const tokenFileName = "token.json"

var logger *log.Entry

// SpotifySaver will handle all the saving logic.
type SpotifySaver struct {
	token oauth2.Token
	client spotify.Client
}

// NewSpotifySaver will create a new empty SpotifySaver instance.
func NewSpotifySaver(log *log.Entry) SpotifySaver {
	logger = log
	return SpotifySaver{}
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