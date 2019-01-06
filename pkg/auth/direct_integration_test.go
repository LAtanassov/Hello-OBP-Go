// +build integration

package auth_test

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/LAtanassov/Hello-OBP-Go/pkg/auth"
)

func TestITDirectSerivce(t *testing.T) {

	t.Run("should login with valid credentials", func(t *testing.T) {

		directUrl := os.Getenv("DIRECT_URL")
		consumerKey := os.Getenv("CONSUMER_KEY")
		username := os.Getenv("USERNAME")
		password := os.Getenv("PASSWORD")

		u, err := url.Parse(directUrl)
		if err != nil {
			t.Errorf("url.Parse(...) error %v", err)
		}

		s := auth.NewDirectService(http.DefaultClient, u, consumerKey)
		_, err = s.Login(username, password)
		if err != nil {
			t.Errorf("s.Login(...) error %v", err)
		}
	})
}
