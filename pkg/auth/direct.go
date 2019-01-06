package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// ErrInvalidCredentials returned by OBP API
var ErrInvalidCredentials = errors.New("OBP-20004: Invalid login credentials. Check username/password.")

// Doer sends an HTTP request and returns an HTTP response, following policy (such as redirects, cookies, auth) as configured on the client.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Login authenticates user and returns token
func Login(doer Doer, u *url.URL, username, password, consumerKey string) (string, error) {

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader([]byte{}))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("DirectLogin username=\"%v\",password=\"%v\",consumer_key=\"%v\"", username, password, consumerKey))

	res, err := doer.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	return unmarshalToken(res.Body)
}

func unmarshalToken(r io.Reader) (string, error) {
	var t struct {
		Token string `json:"token"`
		Error string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(r).Decode(&t); err != nil {
		return "", errors.Wrap(err, "could not unmarshal response")
	}

	if len(t.Error) != 0 {
		switch t.Error {
		case ErrInvalidCredentials.Error():
			return "", ErrInvalidCredentials
		default:
			return "", errors.New(t.Error)
		}
	}

	return t.Token, nil
}

// WithToken adds token as header to request
func WithToken(token string) func(*http.Request) {
	return func(req *http.Request) {
		req.Header.Set("Authorization", fmt.Sprintf("DirectLogin token=\"%v\"", token))
	}
}
