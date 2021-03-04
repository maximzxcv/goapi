package handlers

import (
	"goapi/app/api/middle"
	"goapi/business/apiclient"
	"goapi/business/auth"
	"goapi/business/data/user"
	"goapi/ttesting"
	"log"
	"net/http"
	"strconv"
	"testing"
)

type userTests struct {
	app    http.Handler
	client apiclient.Client
}

func TestUser(t *testing.T) {
	tunit, err := ttesting.NewUnit()
	if err != nil {
		log.Fatalf("Failed to run test: %s", err)
	}

	t.Cleanup(tunit.Teardown)

	app := API(tunit.Db, middle.LoggMiddle())
	client, err := apiclient.BuildClient(app)

	if err != nil {
		log.Fatalf("Failed to run test: %s", err)
	}

	utests := userTests{
		app,
		client,
	}

	if err := utests.singupClient("UserTesterAccount", "testpassword"); err != nil {
		log.Fatalf("Failed to set authorization: %s", err)
	}

	if err := utests.client.Authorize("UserTesterAccount", "testpassword"); err != nil {
		log.Fatalf("Failed to set authorization: %s", err)
	}

	t.Log("User CRUD functionality")
	{
		usr := utests.postUser201(t)
		utests.getUser200(t, usr)
		utests.putUsers200(t, usr.ID)
		utests.deleteUsers204(t, usr.ID)
		utests.getUsersList200(t)
	}
}

func (utests *userTests) singupClient(u string, p string) error {
	signup := auth.Signup{
		Username: u,
		Password: p,
	}

	_, err := utests.client.UnauthorizedCall(http.MethodPost, "/singup", &signup, nil)

	return err
}

// create user
func (utests *userTests) postUser201(t *testing.T) user.User {
	testGoalLog := "postUser201: Should be able to create user."

	cusr := user.CreateUser{
		Name:            "HttpUserName",
		Password:        "testpassword",
		PasswordConfirm: "testpassword",
	}
	var outusr user.User

	httpCode, err := utests.client.Post("/users", &cusr, &outusr)

	if err != nil {
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	}

	expected := http.StatusCreated
	switch {
	case httpCode != expected:
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", expected, httpCode))
	case err != nil:
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	case outusr.Name != cusr.Name:
		t.Error(ttesting.FailedLog(testGoalLog, "user.Name", cusr.Name, outusr.Name))
	default:
		t.Log(ttesting.SuccessLog(testGoalLog))
	}

	return outusr
}

// get single user by id
func (utests *userTests) getUser200(t *testing.T, usr user.User) {
	testGoalLog := "getUsers200: Should be able to get user by id."

	var outusr user.User

	httpCode, err := utests.client.Get("/users/"+usr.ID, &outusr)

	if err != nil {
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	}

	expected := http.StatusOK
	switch {
	case httpCode != expected:
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", expected, httpCode))
	case err != nil:
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	case outusr.Name != usr.Name:
		t.Error(ttesting.FailedLog(testGoalLog, "user.Name", usr.Name, outusr.Name))
	default:
		t.Log(ttesting.SuccessLog(testGoalLog))
	}
}

// update user by id
func (utests *userTests) putUsers200(t *testing.T, uid string) {
	testGoalLog := "putUsers204: Should be able to update user by id."

	inusr := user.UpdateUser{
		Name:            "NewName",
		Password:        "NewPass",
		PasswordConfirm: "NewPass",
	}

	httpCode, err := utests.client.Put("/users/"+uid, &inusr, nil)
	if err != nil {
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	}

	expected := http.StatusOK
	if httpCode != expected {
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", expected, httpCode))
	}

	var outusr user.User
	httpCode, err = utests.client.Get("/users/"+uid, &outusr)

	switch {
	case err != nil:
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	case inusr.Name != outusr.Name:
		t.Error(ttesting.FailedLog(testGoalLog, "user.Name", inusr.Name, outusr.Name))
	default:
		t.Log(ttesting.SuccessLog(testGoalLog))
	}
}

// delete user by id
func (utests *userTests) deleteUsers204(t *testing.T, uid string) {
	testGoalLog := "deleteUsers204: Should be able to delete user by id."

	httpCode, err := utests.client.Delete("/users/" + uid)
	if err != nil {
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	}

	expected := http.StatusNoContent
	if httpCode != expected {
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", expected, httpCode))
	}

	httpCode, _ = utests.client.Get("/users/"+uid, nil)

	expected = http.StatusNotFound
	switch {
	case httpCode != expected:
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", expected, httpCode))
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

		_, err := utests.client.Post("/users", &cusr, nil)
		if err != nil {
			t.Error(ttesting.ErrorLog(testGoalLog, err))
		}
	}

	var outusrs []user.User
	httpCode, err := utests.client.Get("/users", &outusrs)

	expected := http.StatusOK
	switch {
	case httpCode != expected:
		t.Error(ttesting.FailedLog(testGoalLog, "httpStatus", expected, httpCode))
	case err != nil:
		t.Error(ttesting.ErrorLog(testGoalLog, err))
	case len(outusrs) != uamount+1:
		t.Error(ttesting.FailedLog(testGoalLog, "Amount of users", uamount+1, len(outusrs)))
	default:
		t.Log(ttesting.SuccessLog(testGoalLog))
	}
}
