package auth_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/LAtanassov/Hello-OBP-Go/pkg/auth"
	"github.com/pkg/errors"
)

func TestLogin(t *testing.T) {

	t.Run("should login with valid credentials", func(t *testing.T) {

		consumerKey := "consumer_key"
		username := "username"
		password := "password"

		wantToken := "authentication_token"
		wantHeader := fmt.Sprintf("DirectLogin username=\"%v\",password=\"%v\",consumer_key=\"%v\"", username, password, consumerKey)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotHeader := r.Header.Get("Authorization")
			if gotHeader != wantHeader {
				t.Errorf("want header: %v but got header: %v", wantHeader, gotHeader)
			}

			marshalToken(w, wantToken, nil)
			http.Redirect(w, r, "/", 200)
		}))
		defer ts.Close()

		u, err := url.Parse(ts.URL)
		if err != nil {
			t.Errorf("url.Parse(...) error %v", err)
		}

		gotToken, err := auth.Login(http.DefaultClient, u, username, password, consumerKey)
		if err != nil {
			t.Errorf("auth.Login(...) error %v", err)
		}

		if wantToken != gotToken {
			t.Errorf("want token %v but got token: %v", wantToken, gotToken)
		}
	})

	t.Run("should return error if login with invalid credentials", func(t *testing.T) {

		consumerKey := "consumer_key"
		username := "username"
		password := "invalid-password"

		wantError := auth.ErrInvalidCredentials
		wantHeader := fmt.Sprintf("DirectLogin username=\"%v\",password=\"%v\",consumer_key=\"%v\"", username, password, consumerKey)

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotHeader := r.Header.Get("Authorization")
			if gotHeader != wantHeader {
				t.Errorf("want header: %v but got header: %v", wantHeader, gotHeader)
			}

			marshalToken(w, "", wantError)
			http.Redirect(w, r, "/", 401)
		}))
		defer ts.Close()

		url, err := url.Parse(ts.URL)
		if err != nil {
			t.Errorf("url.Parse(...) error %v", err)
		}

		_, err = auth.Login(http.DefaultClient, url, username, password, consumerKey)
		if err == nil {
			t.Errorf("auth.Login(...) want error %v", wantError)
		}

		if err != wantError {
			t.Errorf("want error %v but got error: %v", wantError, err)
		}
	})

}

func marshalToken(w io.Writer, token string, err error) error {
	res := make(map[string]string)
	res["token"] = token
	if err != nil {
		res["error"] = err.Error()
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		return errors.Wrap(err, "could not marshal token")
	}
	return nil
}
