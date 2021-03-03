package handlers

import (
	"bytes"
	"encoding/json"
	"goapi/app/api/middle"
	"goapi/business/auth"
	"goapi/business/data/user"
	"goapi/ttesting"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type authTests struct {
	app http.Handler
}

func TestAuth(t *testing.T) {
	tunit, err := ttesting.NewUnit()
	if err != nil {
		log.Fatalf("Failed to run test: %s", err)
	}

	t.Cleanup(tunit.Teardown)

	atests := authTests{
		app: API(tunit.Db, middle.LoggMiddle()),
	}

	t.Log("Authentication functionality")
	{
		signup := atests.postSingup201(t)
		atests.postLogin401(t)
		access := atests.postLogin200(t, signup.Username, signup.Password)

		atests.getUsers200(t, access)
		atests.getLogout200(t, access)
		atests.getUsers401(t, access)
		atests.getUsers401(t, nil)
	}

}

// Register new account
func (atests *authTests) postSingup201(t *testing.T) *auth.Signup {
	testGoalLog := "postSingup201: Should return status:204"

	signup := auth.Signup{
		Username: "NewUserName",
		Password: "testpassword",
	}
	body, err := json.Marshal(&signup)
	if err != nil {
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	}

	r := httptest.NewRequest(http.MethodPost, "/singup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	atests.app.ServeHTTP(w, r)

	estatus := http.StatusCreated
	if w.Code != estatus {
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", estatus, w.Code))
	} else {
		t.Log(ttesting.SuccessLog(testGoalLog))
	}

	return &signup
}

// use wrong credentials
func (atests *authTests) postLogin401(t *testing.T) *auth.Access {
	testGoalLog := "postLogin401: Should fail to return JWT token"

	login := auth.Login{
		Username: "randomUsername",
		Password: "password",
	}
	body, err := json.Marshal(&login)
	if err != nil {
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	}

	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	atests.app.ServeHTTP(w, r)

	var access auth.Access
	err = json.NewDecoder(w.Body).Decode(&access)

	estatus := http.StatusUnauthorized
	switch {
	case w.Code != estatus:
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", estatus, w.Code))
	default:
		t.Log(ttesting.SuccessLog(testGoalLog))
	}

	return &access
}

// get token
func (atests *authTests) postLogin200(t *testing.T, username string, password string) *auth.Access {
	testGoalLog := "postLogin200: Should return JWT token"

	login := auth.Login{
		Username: username,
		Password: password,
	}
	body, err := json.Marshal(&login)
	if err != nil {
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	}

	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	atests.app.ServeHTTP(w, r)

	var access auth.Access
	err = json.NewDecoder(w.Body).Decode(&access)

	estatus := http.StatusOK
	switch {
	case w.Code != estatus:
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", estatus, w.Code))
	case err != nil:
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	case len(access.Token) < 100:
		t.Error(ttesting.FailedLog(testGoalLog, "access.Token", "len(access.Token)<100", len(access.Token)))
	default:
		t.Log(ttesting.SuccessLog(testGoalLog))
	}

	return &access
}

// Check if token works
func (atests *authTests) getUsers200(t *testing.T, access *auth.Access) {
	testGoalLog := "getUsers200: Should return users."
	r := httptest.NewRequest(http.MethodGet, "/users", nil)

	r.Header.Add("Authorization", "Bearer "+access.Token)
	w := httptest.NewRecorder()
	atests.app.ServeHTTP(w, r)

	var rusrs []user.User
	err := json.NewDecoder(w.Body).Decode(&rusrs)

	uamount := 1
	estatus := http.StatusOK
	switch {
	case w.Code != estatus:
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", estatus, w.Code))
	case err != nil:
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	case len(rusrs) != uamount:
		t.Error(ttesting.FailedLog(testGoalLog, "Number of Users", uamount, len(rusrs)))
	default:
		t.Log(ttesting.SuccessLog(testGoalLog))
	}
}

// Cancel token
func (atests *authTests) getLogout200(t *testing.T, access *auth.Access) {
	testGoalLog := "getLogout200: Should cancel token."

	r := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()
	atests.app.ServeHTTP(w, r)

	estatus := http.StatusOK
	if w.Code != estatus {
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", estatus, w.Code))
	} else {
		t.Log(ttesting.SuccessLog(testGoalLog))
	}
}

// Check if token is canceled
func (atests *authTests) getUsers401(t *testing.T, access *auth.Access) {
	testGoalLog := "getUsers401: Should fail to get users with [Unauthorized] status."

	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	if access != nil { //TODO: use test tables
		//set access
	}
	atests.app.ServeHTTP(w, r)

	estatus := http.StatusUnauthorized
	switch {
	case w.Code != estatus:
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", estatus, w.Code))
	default:
		t.Log(ttesting.SuccessLog(testGoalLog))
	}
}
