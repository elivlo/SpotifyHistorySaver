package login

import (
	"fmt"
	"github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var hook *logtest.Hook

func TestMain(m *testing.M) {
	var logger *logrus.Logger
	logger, hook = logtest.NewNullLogger()
	initLogger(logger)
	code := m.Run()
	os.Exit(code)
}

func TestCreateCodeVerifier(t *testing.T) {
	code := createCodeVerifier(10)
	assert.Equal(t, 14, len(code))
}

func TestCreateVerifierChallenge(t *testing.T) {
	code := createVerifierChallenge("12345abcde")
	fmt.Println(code)
	assert.Equal(t, "PDc_SVO4XN6liOBDbBNMgZ9XC3LB23QOs1z8lCuqK84", code)
}

func TestBase64Escape(t *testing.T) {
	var escape = []byte("any + old & data")

	escaped := base64Escape(escape)
	assert.Equal(t, "YW55ICsgb2xkICYgZGF0YQ", escaped)
}
