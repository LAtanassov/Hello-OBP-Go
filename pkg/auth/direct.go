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

// Doer sends an HTTP request and returns an HTTP response, following policy (such as redirects, cookies, auth) as configured on the client.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewDirectService returns a Direct Login Service
func NewDirectService(doer Doer, url *url.URL, consumerKey string) *DirectService {
	return &DirectService{
		doer:        doer,
		url:         url,
		consumerKey: consumerKey,
	}
}

// DirectService used for Direct Login
type DirectService struct {
	doer        Doer
	url         *url.URL
	consumerKey string
}

// Login authenticates user and returns token
func (s *DirectService) Login(username, password string) (string, error) {

	req, err := http.NewRequest(http.MethodPost, s.url.String(), bytes.NewReader([]byte{}))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("DirectLogin username=\"%v\",password=\"%v\",consumer_key=\"%v\"", username, password, s.consumerKey))

	res, err := s.doer.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	return unmarshalToken(res.Body)
}

func unmarshalToken(r io.Reader) (string, error) {
	var t struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r).Decode(&t); err != nil {
		return "", errors.Wrap(err, "could not unmarshal response")
	}
	return t.Token, nil
}
