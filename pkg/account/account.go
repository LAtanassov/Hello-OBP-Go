package account

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/LAtanassov/Hello-OBP-Go/pkg/auth"
	"github.com/pkg/errors"
)

// ErrUnknown is returned by OBP API
var ErrUnknown = errors.New("OBP-50000: Unknown Error.")

// Account held by the current user
type Account struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	BankID string `json:"bank_id"`
}

// Get returns all accounts held by specific user.
// For Authentication/Authorization provide either
// - a token via WithToken option or
// - use OAuth with Secure Cookie
func Get(doer auth.Doer, u *url.URL, options ...func(*http.Request)) ([]Account, error) {

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

	return unmarshalAccounts(res.Body)
}

func unmarshalAccounts(r io.Reader) ([]Account, error) {
	var t struct {
		Accounts []Account `json:"accounts"`
		Error    string    `json:"error,omitempty"`
	}

	if err := json.NewDecoder(r).Decode(&t); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal response")
	}

	if len(t.Error) != 0 {
		switch t.Error {
		case ErrUnknown.Error():
			return nil, ErrUnknown
		default:
			return nil, errors.New(t.Error)
		}
	}

	return t.Accounts, nil
}
