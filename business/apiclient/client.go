package apiclient

import (
	"bytes"
	"encoding/json"
	"goapi/business/auth"
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
func (client *Client) Post(target string, in interface{}, out interface{}) (int, error) {
	return client.call(http.MethodPost, target, true, in, out)
}

// Get ....
func (client *Client) Get(target string, out interface{}) (int, error) {
	return client.call(http.MethodGet, target, true, nil, out)
}

// Put ....
func (client *Client) Put(target string, in interface{}, out interface{}) (int, error) {
	return client.call(http.MethodPut, target, true, in, out)
}

// Delete ....
func (client *Client) Delete(target string) (int, error) {
	return client.call(http.MethodDelete, target, true, nil, nil)
}

// UnauthorizedCall ...
func (client *Client) UnauthorizedCall(method, target string, in interface{}, out interface{}) (int, error) {
	return client.call(method, target, false, in, out)
}

// Login user to use client
func (client *Client) Login(username string, password string) error {
	login := auth.Login{
		Username: username,
		Password: password,
	}

	var access auth.Access
	_, err := client.call(http.MethodPost, "/login", false, &login, &access)
	if err != nil {
		return err
	}

	client.authStr = "Bearer " + access.Token

	return nil
}

func (client *Client) call(method, target string, isAuth bool, in interface{}, out interface{}) (int, error) {
	r := &http.Request{}
	if in == nil {
		r = httptest.NewRequest(method, target, nil)
	} else {
		body, err := json.Marshal(&in)
		if err != nil {
			return 0, errors.WithStack(err)
		}
		r = httptest.NewRequest(method, target, bytes.NewBuffer(body))
	}
	if isAuth {
		r.Header.Add("Authorization", client.authStr)
	}
	w := httptest.NewRecorder()
	client.api.ServeHTTP(w, r)

	if out == nil {
		return w.Code, nil
	}

	err := json.NewDecoder(w.Body).Decode(&out)
	defer r.Body.Close()

	return w.Code, err
}
