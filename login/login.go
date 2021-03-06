package login

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

const (
	// TokenFileName is the standard file name to save the OAuth token to
	TokenFileName = "token.json"
)

// Auth is the interface Login implements. It supports log in to the Spotify account.
// You can also save the token to a file.
type Auth interface {
	Login() *oauth2.Token
	SaveToken(string, *oauth2.Token) error
}

// Login is the type to use when logging into a Spotify account.
type Login struct {
	logger        *logrus.Entry
	callbackURI   string
	ch            chan *oauth2.Token
	auth          SpotifyAuthenticatior
	state         string
	codeVerifier  string
	codeChallenge string
}

// NewLogin creates a new Login with the given callbackURL to listen on.
// It will also create a code verifier for this login.
func NewLogin(callbackURL, clientID, clientSecret string) Login {
	login := Login{
		logger:       initLogger(logrus.New()),
		callbackURI:  callbackURL,
		state:        createCodeVerifier(20),
		codeVerifier: createCodeVerifier(96),
		ch:           make(chan *oauth2.Token),
	}
	login.codeChallenge = createVerifierChallenge(login.codeVerifier)

	// creates new Authenticator
	login.auth = spotifyauth.New(spotifyauth.WithRedirectURL(login.callbackURI),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadRecentlyPlayed),
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithClientSecret(clientSecret))

	return login
}

// Login wil open a http server to log in to your account to get a newly created OAuth2 token.
func (l Login) Login() *oauth2.Token {
	// setup http server
	servMux := http.NewServeMux()
	servMux.HandleFunc("/callback", l.authHandler)
	server := http.Server{
		Addr: ":8080",
		Handler: servMux,
	}

	// run http server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.ch <- nil
			l.logger.Fatal(err)
		}
	}()

	u := l.auth.AuthURL(l.state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", l.codeChallenge),
	)
	ur, _ := url.PathUnescape(u)
	l.logger.Info("Please log in to Spotify by visiting the following page in your browser: ", ur)

	// wait for auth to complete
	token := <-l.ch

	_ = server.Close()

	l.logger.Info("You are logged in")
	return token
}

// SaveToken will save access and refresh token to token.json file in exec directory.
func (l Login) SaveToken(file string, token *oauth2.Token) error {
	fileString, err := json.Marshal(token)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}

	_, err = f.Write(fileString)
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	l.logger.Infof("Wrote access token to %s/%s", dir, file)
	return nil
}

// authHandler will handle the incoming token from Spotify.
func (l Login) authHandler(w http.ResponseWriter, r *http.Request) {
	token, err := l.auth.Token(r.Context(), l.state, r,
		oauth2.SetAuthURLParam("code_verifier", l.codeVerifier))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		l.logger.Fatal(err)
	}
	if st := r.FormValue("state"); st != l.state {
		http.NotFound(w, r)
		l.logger.Fatalf("State mismatch: %s != %s", st, l.state)
	}
	_, err = fmt.Fprintf(w, "Login Completed!")
	if err != nil {
		l.logger.Fatal(err)
	}
	l.ch <- token
}

// initLogger inits a logger with "ACCOUNT LOGIN" prefix.
func initLogger(logger *logrus.Logger) *logrus.Entry {
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&nested.Formatter{
		FieldsOrder: []string{"component", "category"},
		HideKeys:    true,
	})
	return logger.WithField("component", "LOGIN")
}

// createCodeVerifier will create a random base64 encoded verifier.
func createCodeVerifier(size int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	b := make([]byte, size)
	for i := range b {
		b[i] = byte(r.Intn(255))
	}

	return base64Escape(b)
}

// createCodeVerifier will create a sha265 base64 encoded sum for verifier.
func createVerifierChallenge(v string) string {
	c := sha256.New()
	c.Write([]byte(v))
	return base64Escape(c.Sum(nil))
}

// base64Escape escapes some runes that would throw an error when authenticating with OAuth2.
func base64Escape(b []byte) string {
	e := base64.StdEncoding.EncodeToString(b)
	e = strings.ReplaceAll(e, "+", "-")
	e = strings.ReplaceAll(e, "/", "_")
	e = strings.ReplaceAll(e, "=", "")
	return e
}
