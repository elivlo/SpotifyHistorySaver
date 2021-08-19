package spotifySaver

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestMockedSpotifySaver_LoadToken(t *testing.T) {
	mock := MockedSpotifySaver{}

	err := mock.LoadToken("")
	assert.NoError(t, err)

	mock.LError = true
	err = mock.LoadToken("")
	assert.Error(t, err)
}

func TestMockedSpotifySaver_Authenticate(t *testing.T) {
	mock := MockedSpotifySaver{}
	mock.Authenticate("", "", "")
}

func TestMockedSpotifySaver_StartLastSongsWorker(t *testing.T) {
	mock := MockedSpotifySaver{}
	var wg sync.WaitGroup

	wg.Add(1)
	stop := make(chan bool, 1)
	mock.StartLastSongsWorker(&wg, stop)

	wg.Wait()
}