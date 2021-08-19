package login

import (
	"errors"
	"golang.org/x/oauth2"
)


type MockedAuth struct {
	LError bool
	SError bool
}

// Login wil open a http server to log in to your account to get a newly created OAuth2 token.
func (l MockedAuth) Login(_, _ string) (*oauth2.Token, error) {
	if l.LError {
		return nil, errors.New("login error")
	}
	return &oauth2.Token{}, nil
}

// SaveToken will save access and refresh token to token.json file in exec directory.
func (l MockedAuth) SaveToken(_ *oauth2.Token) error {
	if l.SError {
		return errors.New("save error")
	}
	return nil
}
