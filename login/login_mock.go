package login

import (
	"errors"
	"golang.org/x/oauth2"
)

// MockedAuth implements the Auth interface for tests.
type MockedAuth struct {
	LError bool
	SError bool
}

// Login will return a token or error.
func (l MockedAuth) Login(_, _ string) (*oauth2.Token, error) {
	if l.LError {
		return nil, errors.New("login error")
	}
	return &oauth2.Token{}, nil
}

// SaveToken mocks saving the token.
func (l MockedAuth) SaveToken(_ string, _ *oauth2.Token) error {
	if l.SError {
		return errors.New("save error")
	}
	return nil
}

/*type SpotifyAuthenticatior interface {
	AuthURL(state string, opts ...oauth2.AuthCodeOption) string
	Token(ctx context.Context, state string, r *http.Request, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	Client(ctx context.Context, token *oauth2.Token) *http.Client
}

type MockedSpotifyauthAuthenticator struct {
	config *oauth2.Config
}

func NewMockedSpotifyauthAuthenticator(_ ...spotifyauth.AuthenticatorOption) *MockedSpotifyauthAuthenticator {
	return &MockedSpotifyauthAuthenticator{}
}

func (a MockedSpotifyauthAuthenticator) AuthURL(_ string, _ ...oauth2.AuthCodeOption) string {
	return "authUrl"
}

// Token pulls an authorization code from an HTTP request and attempts to exchange
// it for an access token.  The standard use case is to call Token from the handler
// that handles requests to your application's redirect URL.
func (a MockedSpotifyauthAuthenticator) Token(_ context.Context, _ string, _ *http.Request, _ ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return &oauth2.Token{}, nil
}

// Exchange is like Token, except it allows you to manually specify the access
// code instead of pulling it out of an HTTP request.
func (a MockedSpotifyauthAuthenticator) Exchange(_ context.Context, _ string, _ ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return &oauth2.Token{}, nil
}

// Client creates a *http.Client that will use the specified access token for its API requests.
// Combine this with spotify.HTTPClientOpt.
func (a MockedSpotifyauthAuthenticator) Client(_ context.Context, _ *oauth2.Token) *http.Client {
	return &http.Client{}
}*/