package apiclient

import (
	"bytes"
	"encoding/json"
	"goapi/business/auth"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/pkg/errors"
)

// Client - Http client for goapi
type Client struct {
	api     http.Handler
	authStr string
}

// BuildClient is constructor for Client
func BuildClient(api http.Handler) (Client, error) {

	c := Client{
		api: api,
	}

	return c, nil
}

// Post ....
func (client *Client) Post(target string, input interface{}, output interface{}) (int, error) {
	body, err := json.Marshal(&input)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	r := httptest.NewRequest(http.MethodPost, target, bytes.NewBuffer(body))
	r.Header.Add("Authorization", client.authStr)
	w := httptest.NewRecorder()
	client.api.ServeHTTP(w, r)

	if output == nil {
		return w.Code, nil
	}

	// set output
	return w.Code, decode(w.Body, output)
}

func (client *Client) UnauthorizedCall(method, target string, input interface{}, output interface{}) (int, error) {
	body, err := json.Marshal(&input)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	r := httptest.NewRequest(http.MethodPost, target, bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	client.api.ServeHTTP(w, r)

	if output == nil {
		return w.Code, nil
	}

	// set output
	return w.Code, decode(w.Body, output)
}

func decode(r io.Reader, n interface{}) error {
	return json.NewDecoder(r).Decode(&n)
}

func (client *Client) Authorize(username string, password string) error {
	login := auth.Login{
		Username: username,
		Password: password,
	}
	body, err := json.Marshal(&login)
	if err != nil {
		return err
	}

	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	client.api.ServeHTTP(w, r)

	var access auth.Access
	err = json.NewDecoder(w.Body).Decode(&access)
	if err != nil {
		return err
	}

	client.authStr = "Bearer " + access.Token
	return nil
}
