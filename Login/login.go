package Login


import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/zmb3/spotify"
)

const tokenFileName = "token.json"

var (
	logger *log.Entry
	redirectURI string
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadRecentlyPlayed)
	ch    = make(chan *spotify.Client)
	state = createCodeVerifier(20)
	codeVerifier = createCodeVerifier(96)
	codeChallenge = createVerifierChallenge(codeVerifier)
)

// Login wil open a http server to log in to your account to get a newly created OAuth2 token.
func Login(clientID, clientSecret, callbackURL string) (*oauth2.Token, error) {
	redirectURI = callbackURL
	initLogger()

	// set Spotify App Client Credentials
	auth.SetAuthInfo(clientID, clientSecret)

	// start HTTP callback server
	http.HandleFunc("/callback", authHandler)
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			logger.Fatal(err)
		}
	}()

	u := auth.AuthURLWithOpts(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)
	ur, _ := url.PathUnescape(u)
	logger.Info("Please log in to Spotify by visiting the following page in your browser:", ur)

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("You are logged in as:", user.ID)

	return client.Token()
}

// SaveToken will save access and refresh token to token.json file in exec directory.
func SaveToken(token *oauth2.Token) error {
	fileString, err := json.Marshal(token)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(tokenFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}

	_, err = file.Write(fileString)
	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	logger.Info("Wrote access token to %s/%s", dir, tokenFileName)
	return nil
}

// initLogger inits a logger with "ACCOUNT LOGIN" prefix.
func initLogger() {
	l := log.New()
	l.SetLevel(log.InfoLevel)
	l.SetFormatter(&nested.Formatter{
		FieldsOrder:   []string{"component", "category"},
		HideKeys:      true,
	})
	logger = l.WithField("component", "ACCOUNT LOGIN")
}

// authHandler will handle the incoming token from Spotify.
func authHandler(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.TokenWithOpts(state, r,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		logger.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		logger.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	_, err = fmt.Fprintf(w, "Login Completed!")
	if err != nil {
		logger.Fatal(err)
	}
	ch <- &client
}

// createCodeVerifier will create a random base64 encoded verifier.
func createCodeVerifier(size int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	b := make([]byte, size)
	for i, _ := range b {
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
