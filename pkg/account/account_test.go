package account_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/LAtanassov/Hello-OBP-Go/pkg/account"
	"github.com/pkg/errors"
)

func TestGet(t *testing.T) {

	t.Run("should get accounts", func(t *testing.T) {

		wantAccounts := []account.Account{account.Account{ID: "myID"}}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			marshalAccounts(w, wantAccounts, nil)
			http.Redirect(w, r, "/", 200)
		}))
		defer ts.Close()

		u, err := url.Parse(ts.URL)
		if err != nil {
			t.Errorf("url.Parse(...) error %v", err)
		}

		gotAccounts, err := account.Get(http.DefaultClient, u)
		if err != nil {
			t.Errorf("account.Get(...) error %v", err)
		}

		if !reflect.DeepEqual(wantAccounts, gotAccounts) {
			t.Errorf("want accounts %v but got accounts: %v", wantAccounts, gotAccounts)
		}
	})
}

func marshalAccounts(w io.Writer, accounts []account.Account, err error) error {
	res := make(map[string]interface{})
	res["accounts"] = accounts
	if err != nil {
		res["error"] = err.Error()
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		return errors.Wrap(err, "could not marshal accounts")
	}
	return nil
}
