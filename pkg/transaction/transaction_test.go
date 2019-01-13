package transaction_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/LAtanassov/Hello-OBP-Go/pkg/transaction"

	"github.com/pkg/errors"
)

func TestGet(t *testing.T) {

	t.Run("should get transactions", func(t *testing.T) {

		wantTransactions := []transaction.Transaction{}

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			marshalTransactions(w, wantTransactions, nil)
			http.Redirect(w, r, "/", 200)
		}))
		defer ts.Close()

		u, err := url.Parse(ts.URL)
		if err != nil {
			t.Errorf("url.Parse(...) error %v", err)
		}

		gotTransactions, err := transaction.Get(http.DefaultClient, u)
		if err != nil {
			t.Errorf("transaction.Get(...) error %v", err)
		}

		if !reflect.DeepEqual(wantTransactions, gotTransactions) {
			t.Errorf("want transaction %v but got transaction: %v", wantTransactions, gotTransactions)
		}
	})
}

func marshalTransactions(w io.Writer, transactions []transaction.Transaction, err error) error {
	res := make(map[string]interface{})
	res["transactions"] = transactions
	if err != nil {
		res["error"] = err.Error()
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		return errors.Wrap(err, "could not marshal transactions")
	}
	return nil
}
