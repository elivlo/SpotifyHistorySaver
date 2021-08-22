package login

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMockedAuth_Login(t *testing.T) {
	mock := MockedAuth{}

	_, err := mock.Login()
	assert.NoError(t, err)

	mock.LError = true
	_, err = mock.Login()
	assert.Error(t, err)
}

func TestMockedAuth_SaveToken(t *testing.T) {
	mock := MockedAuth{}

	err := mock.SaveToken("", nil)
	assert.NoError(t, err)

	mock.SError = true
	err = mock.SaveToken("", nil)
	assert.Error(t, err)
}
