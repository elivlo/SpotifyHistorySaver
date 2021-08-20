package spotifySaver

import (
	"errors"
	"sync"
	"time"
)

// MockedSpotifySaver implements the InterfaceSpotifySaver interface for tests.
type MockedSpotifySaver struct {
	LError bool
}

// LoadToken will load the token from file "token.json" in exec directory.
// It will throw an error when the token is expired.
func (s *MockedSpotifySaver) LoadToken(_ string) error {
	if s.LError {
		return errors.New("load Token error")
	}
	return nil
}

// Authenticate will create a new client from token.
func (s *MockedSpotifySaver) Authenticate(_, _, _ string) {}

// StartLastSongsWorker is a worker that will send history requests every 45 minutes.
// It is not async. It accepts a wait group and will send Done when stopped. It may be stopped with stop chan value.
func (s *MockedSpotifySaver) StartLastSongsWorker(wg *sync.WaitGroup, stop chan bool) {
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			stop <- true
		case <-stop:
			ticker.Stop()
			wg.Done()
			return
		}
	}
}
