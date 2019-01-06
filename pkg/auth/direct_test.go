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

func TestDirectSerivce(t *testing.T) {

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

			marshalToken(w, wantToken)
			http.Redirect(w, r, "/", 200)
		}))
		defer ts.Close()

		u, err := url.Parse(ts.URL)
		if err != nil {
			t.Errorf("url.Parse(...) error %v", err)
		}

		s := auth.NewDirectService(http.DefaultClient, u, consumerKey)
		gotToken, err := s.Login(username, password)
		if err != nil {
			t.Errorf("s.Login(...) error %v", err)
		}

		if wantToken != gotToken {
			t.Errorf("want token %v but got token: %v", wantToken, gotToken)
		}
	})

}

func marshalToken(w io.Writer, token string) error {
	t := map[string]string{"token": token}
	if err := json.NewEncoder(w).Encode(t); err != nil {
		return errors.Wrap(err, "could not marshal token")
	}
	return nil
}
