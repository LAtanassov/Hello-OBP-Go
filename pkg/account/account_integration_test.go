// +build integration

package account_test

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/LAtanassov/Hello-OBP-Go/pkg/account"
	"github.com/LAtanassov/Hello-OBP-Go/pkg/auth"
)

func TestITGet(t *testing.T) {

	t.Run("should get accounts", func(t *testing.T) {

		directUrl := os.Getenv("DIRECT_URL")
		accountUrl := os.Getenv("ACCOUNT_URL")
		consumerKey := os.Getenv("CONSUMER_KEY")
		username := os.Getenv("USERNAME")
		password := os.Getenv("PASSWORD")

		d, err := url.Parse(directUrl)
		if err != nil {
			t.Errorf("url.Parse(...) error %v", err)
		}

		token, err := auth.Login(http.DefaultClient, d, username, password, consumerKey)
		if err != nil {
			t.Errorf("auth.Login(...) error %v", err)
		}

		a, err := url.Parse(accountUrl)
		if err != nil {
			t.Errorf("url.Parse(...) error %v", err)
		}
		_, err = account.Get(http.DefaultClient, a, auth.WithToken(token))
	})
}
