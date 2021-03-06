package login

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

// MockedAuth implements the Auth interface for tests.
type MockedAuth struct {
	LError bool
	SError bool
}

// Login will return a token or error.
func (l MockedAuth) Login() *oauth2.Token {
	if l.LError {
		return &oauth2.Token{}
	}
	return &oauth2.Token{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
		Expiry:       time.Now().Add(time.Hour),
	}
}

// SaveToken mocks saving the token.
func (l MockedAuth) SaveToken(_ string, _ *oauth2.Token) error {
	if l.SError {
		return errors.New("save error")
	}
	return nil
}

// SpotifyAuthenticatior is an interface for spotifyauth.Authenticator
type SpotifyAuthenticatior interface {
	AuthURL(state string, opts ...oauth2.AuthCodeOption) string
	Token(ctx context.Context, state string, r *http.Request, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	Client(ctx context.Context, token *oauth2.Token) *http.Client
}

// MockedSpotifyauthAuthenticator implements interface SpotifyAuthenticatior for tests
type MockedSpotifyauthAuthenticator struct {
	FailToken bool
}

// AuthURL returns a URL to the the Spotify Accounts Service's OAuth2 endpoint.
func (a MockedSpotifyauthAuthenticator) AuthURL(_ string, _ ...oauth2.AuthCodeOption) string {
	return "http://test.test/test"
}

// Token pulls an authorization code from an HTTP request and attempts to exchange
// it for an access token.  The standard use case is to call Token from the handler
// that handles requests to your application's redirect URL.
func (a MockedSpotifyauthAuthenticator) Token(_ context.Context, _ string, _ *http.Request, _ ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	if a.FailToken {
		return nil, errors.New("No token")
	}
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
}

// MockedResponseWriter implements http.ResponseWriter for tests
type MockedResponseWriter struct {
	body       []byte
	statusCode int
	header     http.Header
}

// Header returns the HTTP header
func (w *MockedResponseWriter) Header() http.Header {
	return w.header
}

// Write writes some bytes to the body
func (w *MockedResponseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return len(b), nil
}

// WriteHeader sets the status code
func (w *MockedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}