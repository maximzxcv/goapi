package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"goapi/business/auth"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// CleverClient - Http client for goapi
type CleverClient struct {
	baseURL  string
	authStr  string
	httpclnt *http.Client
}

// BuildClever is constructor for CleverClient
func BuildClever(baseURL string) *CleverClient {
	return &CleverClient{
		baseURL: baseURL,
		httpclnt: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// Login user to use client
func (client *CleverClient) Login(username string, password string) error {
	login := auth.Login{
		Username: username,
		Password: password,
	}

	var access auth.Access
	_, err := client.call(http.MethodPost, "/login", false, &login, &access, true)
	if err != nil {
		return err
	}

	client.authStr = "Bearer " + access.Token

	return nil
}

// Post ....
func (client *CleverClient) Post(target string, in interface{}, out interface{}) (int, error) {
	return client.call(http.MethodPost, target, true, in, out, false)
}

// Get ....
func (client *CleverClient) Get(target string, out interface{}) (int, error) {
	return client.call(http.MethodGet, target, true, nil, out, true)
}

// Put ....
func (client *CleverClient) Put(target string, in interface{}, out interface{}) (int, error) {
	return client.call(http.MethodPut, target, true, in, out, false)
}

// Delete ....
func (client *CleverClient) Delete(target string) (int, error) {
	return client.call(http.MethodDelete, target, true, nil, nil, false)
}

// UnauthorizedCall ...
func (client *CleverClient) UnauthorizedCall(method, target string, in interface{}, out interface{}) (int, error) {
	return client.call(method, target, false, in, out, false)
}

func (client *CleverClient) buildRequest(method, target string, in interface{}) (*http.Request, error) {
	if in == nil {
		return http.NewRequest(method, target, nil)
	}

	json, err := json.Marshal(&in)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return http.NewRequest(method, target, bytes.NewBuffer(json))

}
func (client *CleverClient) call(method, target string, isAuth bool, in interface{}, out interface{}, allowRetry bool) (int, error) {

	req, err := client.buildRequest(method, target, in)
	req.URL.Host = client.baseURL
	req.URL.Scheme = "http"

	if err != nil {
		return 0, errors.WithStack(err)
	}

	if isAuth {
		req.Header.Add("Authorization", client.authStr)
	}
	i := 2
	if allowRetry {
		i = 1
	}

	ctx, cancel := context.WithTimeout(req.Context(), time.Duration(i)*300*time.Millisecond)
	defer cancel()

	res, err := client.httpclnt.Do(req.WithContext(ctx))
	if err != nil {
		if allowRetry && errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Re-call [%d]: %v", i, req.URL)
			return client.call(method, target, isAuth, in, out, false)
		}
		return 0, errors.WithStack(err)
	}
	defer res.Body.Close()

	if out == nil {
		return res.StatusCode, nil
	}

	return res.StatusCode, json.NewDecoder(res.Body).Decode(&out)
}

func (client *CleverClient) getResponse(req *http.Request) (*http.Response, error) {
	res := &http.Response{}
	err := errors.New("")
	c := req.Context()
	r := req
	for i := 1; i < 3; i++ {
		ctx, cancel := context.WithTimeout(c, time.Duration(i)*30*time.Millisecond)
		defer cancel()

		res, err = client.httpclnt.Do(r.WithContext(ctx))
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			log.Printf("Re-call [%d]: %v", i, req.URL)
			continue
		case err != nil:
			break
		default:
			return res, nil
		}
	}

	return res, err
}
