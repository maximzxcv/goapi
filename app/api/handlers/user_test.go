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
	"strconv"
	"testing"
)

type userTests struct {
	app     http.Handler
	authStr string
}

func TestUser(t *testing.T) {
	tunit, err := ttesting.NewUnit()
	if err != nil {
		log.Fatalf("Failed to run test: %s", err)
	}

	t.Cleanup(tunit.Teardown)

	utests := userTests{
		app: API(tunit.Db, middle.LoggMiddle()),
	}

	if err := utests.setAuthStr(); err != nil {
		log.Fatalf("Failed to set authorization: %s", err)
	}

	t.Log("User CRUD functionality")
	{
		usr := utests.postUser201(t)
		utests.getUser200(t, usr)
		utests.putUsers200(t, usr)
		utests.deleteUsers204(t, usr)
		utests.getUsersList200(t)
	}
}

// create user
func (utests *userTests) postUser201(t *testing.T) user.User {
	testGoalLog := "postUser201: Should be able to create user."

	cusr := user.CreateUser{
		Name:            "HttpUserName",
		Password:        "testpassword",
		PasswordConfirm: "testpassword",
	}
	body, err := json.Marshal(&cusr)
	if err != nil {
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	}

	r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	r.Header.Add("Authorization", utests.authStr)
	w := httptest.NewRecorder()
	utests.app.ServeHTTP(w, r)

	var rusr user.User
	err = json.NewDecoder(w.Body).Decode(&rusr)

	estatus := http.StatusCreated
	switch {
	case w.Code != estatus:
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", estatus, w.Code))
	case err != nil:
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	case rusr.Name != cusr.Name:
		t.Error(ttesting.FailedLog(testGoalLog, "user.Name", cusr.Name, rusr.Name))
	default:
		t.Log(ttesting.SuccessLog(testGoalLog))
	}

	return rusr
}

// get single user by id
func (utests *userTests) getUser200(t *testing.T, usr user.User) {
	testGoalLog := "getUsers200: Should be able to get user by id."

	r := httptest.NewRequest(http.MethodGet, "/users/"+usr.ID, nil)
	r.Header.Add("Authorization", utests.authStr)
	w := httptest.NewRecorder()
	utests.app.ServeHTTP(w, r)

	var rusr user.User
	err := json.NewDecoder(w.Body).Decode(&rusr)

	estatus := http.StatusOK
	switch {
	case w.Code != estatus:
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", estatus, w.Code))
	case err != nil:
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	case rusr.Name != usr.Name:
		t.Error(ttesting.FailedLog(testGoalLog, "user.Name", usr.Name, rusr.Name))
	default:
		t.Log(ttesting.SuccessLog(testGoalLog))
	}
}

// update user by id
func (utests *userTests) putUsers200(t *testing.T, usr user.User) {
	testGoalLog := "putUsers204: Should be able to update user by id."

	uusr := user.UpdateUser{
		Name:            "NewName",
		Password:        "NewPass",
		PasswordConfirm: "NewPass",
	}
	body, err := json.Marshal(&uusr)
	if err != nil {
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	}

	r := httptest.NewRequest(http.MethodPut, "/users/"+usr.ID, bytes.NewBuffer(body))
	r.Header.Add("Authorization", utests.authStr)
	w := httptest.NewRecorder()

	utests.app.ServeHTTP(w, r)

	estatus := http.StatusOK
	if w.Code != estatus {
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", estatus, w.Code))
	}

	// get it from service to check results
	r = httptest.NewRequest(http.MethodGet, "/users/"+usr.ID, nil)
	r.Header.Add("Authorization", utests.authStr)
	w = httptest.NewRecorder()

	utests.app.ServeHTTP(w, r)
	var rusr user.User
	err = json.NewDecoder(w.Body).Decode(&rusr)

	switch {
	case err != nil:
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	case uusr.Name != rusr.Name:
		t.Error(ttesting.FailedLog(testGoalLog, "user.Name", uusr.Name, rusr.Name))
	default:
		t.Log(ttesting.SuccessLog(testGoalLog))
	}
}

// delete user by id
func (utests *userTests) deleteUsers204(t *testing.T, usr user.User) {
	testGoalLog := "putUsers204: Should be able to delete user by id."

	r := httptest.NewRequest(http.MethodDelete, "/users/"+usr.ID, nil)
	r.Header.Add("Authorization", utests.authStr)
	w := httptest.NewRecorder()

	utests.app.ServeHTTP(w, r)

	estatus := http.StatusNoContent
	if w.Code != estatus {
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", estatus, w.Code))
	}

	// get it from service to check results
	r = httptest.NewRequest(http.MethodGet, "/users/"+usr.ID, nil)
	r.Header.Add("Authorization", utests.authStr)
	w = httptest.NewRecorder()

	utests.app.ServeHTTP(w, r)
	estatus = http.StatusNotFound

	switch {
	case w.Code != estatus:
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", estatus, w.Code))
	default:
		t.Log(ttesting.SuccessLog(testGoalLog))
	}
}

// get multiple users by id
func (utests *userTests) getUsersList200(t *testing.T) {
	testGoalLog := "getUsers200: Should be able to get LIST of users by id."

	uamount := 7
	cusr := user.CreateUser{
		Name:            "HttpUserName",
		Password:        "testpassword",
		PasswordConfirm: "testpassword",
	}
	// populate more users
	for i := 0; i < uamount; i++ {
		cusr.Name = cusr.Name + strconv.Itoa(i)
		body, err := json.Marshal(&cusr)
		if err != nil {
			t.Error(ttesting.ErrorLog(testGoalLog, err))
		}

		r := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		r.Header.Add("Authorization", utests.authStr)
		w := httptest.NewRecorder()
		utests.app.ServeHTTP(w, r)
	}

	r := httptest.NewRequest(http.MethodGet, "/users", nil)
	r.Header.Add("Authorization", utests.authStr)
	w := httptest.NewRecorder()
	utests.app.ServeHTTP(w, r)

	var rusrs []user.User
	err := json.NewDecoder(w.Body).Decode(&rusrs)

	estatus := http.StatusOK
	switch {
	case w.Code != estatus:
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", estatus, w.Code))
	case err != nil:
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	case len(rusrs) != uamount+1:
		t.Error(ttesting.FailedLog(testGoalLog, "Amount of users", uamount+1, len(rusrs)))
	default:
		t.Log(ttesting.SuccessLog(testGoalLog))
	}
}

func (usets *userTests) setAuthStr() error {
	signup := auth.Signup{
		Username: "UserTester",
		Password: "testpassword",
	}
	body, err := json.Marshal(&signup)
	if err != nil {
		return err
	}

	r := httptest.NewRequest(http.MethodPost, "/singup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	usets.app.ServeHTTP(w, r)

	login := auth.Login{
		Username: signup.Username,
		Password: signup.Password,
	}
	body, err = json.Marshal(&login)
	if err != nil {
		return err
	}

	r = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	usets.app.ServeHTTP(w, r)

	var access auth.Access
	err = json.NewDecoder(w.Body).Decode(&access)
	if err != nil {
		return err
	}

	usets.authStr = "Bearer " + access.Token
	return nil
}
