package login

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"net/http"
	"testing"
)

func TestMockedAuth_Login(t *testing.T) {
	mock := MockedAuth{}

	tok := mock.Login()
	assert.True(t, tok.Valid())

	mock.LError = true

	tok = mock.Login()
	assert.False(t, tok.Valid())
}

func TestMockedAuth_SaveToken(t *testing.T) {
	mock := MockedAuth{}

	err := mock.SaveToken("", nil)
	assert.NoError(t, err)

	mock.SError = true
	err = mock.SaveToken("", nil)
	assert.Error(t, err)
}

func TestMockedSpotifyauthAuthenticator_Token(t *testing.T) {
	mock := MockedSpotifyauthAuthenticator{}

	tok, err := mock.Token(nil, "", nil)
	assert.NoError(t, err)
	assert.Equal(t, &oauth2.Token{}, tok)

	mock.FailToken = true
	tok, err = mock.Token(nil, "", nil)
	assert.Error(t, err)
	assert.Nil(t, tok)
}

func TestMockedSpotifyauthAuthenticator(t *testing.T) {
	mock := MockedSpotifyauthAuthenticator{}

	tok, err := mock.Exchange(nil, "")
	assert.NoError(t, err)
	assert.Equal(t, &oauth2.Token{}, tok)

	client := mock.Client(nil, nil)
	assert.Equal(t, &http.Client{}, client)
}
