package login

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
)

func TestMockedAuth_Login(t *testing.T) {
	mock := MockedAuth{}

	tok := mock.Login()
	assert.Equal(t, &oauth2.Token{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
	}, tok)

}

func TestMockedAuth_SaveToken(t *testing.T) {
	mock := MockedAuth{}

	err := mock.SaveToken("", nil)
	assert.NoError(t, err)

	mock.SError = true
	err = mock.SaveToken("", nil)
	assert.Error(t, err)
}
