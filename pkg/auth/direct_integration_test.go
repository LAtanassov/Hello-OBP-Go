// +build integration

package auth_test

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/LAtanassov/Hello-OBP-Go/pkg/auth"
)

func TestITLogin(t *testing.T) {

	t.Run("should login with valid credentials", func(t *testing.T) {

		directUrl := os.Getenv("DIRECT_URL")
		consumerKey := os.Getenv("CONSUMER_KEY")
		username := os.Getenv("USERNAME")
		password := os.Getenv("PASSWORD")

		u, err := url.Parse(directUrl)
		if err != nil {
			t.Errorf("url.Parse(...) error %v", err)
		}

		_, err = auth.Login(http.DefaultClient, u, username, password, consumerKey)
		if err != nil {
			t.Errorf("auth.Login(...) error %v", err)
		}
	})

	t.Run("should return ErrInvalidCredentials if login with invalid credentials", func(t *testing.T) {

		directUrl := os.Getenv("DIRECT_URL")
		consumerKey := os.Getenv("CONSUMER_KEY")
		username := os.Getenv("USERNAME")
		password := "invalid-password"

		u, err := url.Parse(directUrl)
		if err != nil {
			t.Errorf("url.Parse(...) error %v", err)
		}

		_, err = auth.Login(http.DefaultClient, u, username, password, consumerKey)

		if err == nil {
			t.Errorf("auth.Login(...) want error %v", auth.ErrInvalidCredentials)
		}

		if err != auth.ErrInvalidCredentials {
			t.Errorf("auth.Login(...) want error %v but got error %v", auth.ErrInvalidCredentials, err)
		}
	})
}
