package transaction

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/LAtanassov/Hello-OBP-Go/pkg/auth"
	"github.com/pkg/errors"
)

// Transaction of a specific account
type Transaction struct {
	ID string `json:"id"`
}

// Get transactions of a specific account (encoded in url).
func Get(doer auth.Doer, u *url.URL, options ...func(*http.Request)) ([]Transaction, error) {

	req, err := http.NewRequest(http.MethodGet, u.String(), bytes.NewReader([]byte{}))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for _, option := range options {
		option(req)
	}

	res, err := doer.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return unmarshalTransactions(res.Body)
}

func unmarshalTransactions(r io.Reader) ([]Transaction, error) {
	var t struct {
		Transactions []Transaction `json:"transactions"`
		Error        string        `json:"error,omitempty"`
	}

	if err := json.NewDecoder(r).Decode(&t); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal response")
	}

	if len(t.Error) != 0 {
		switch t.Error {
		default:
			return nil, errors.New(t.Error)
		}
	}

	return t.Transactions, nil
}
