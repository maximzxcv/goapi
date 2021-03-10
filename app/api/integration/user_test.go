package integration

import (
	"context"
	"fmt"
	"goapi/business/apiclient"
	"goapi/business/auth"
	"goapi/business/data/user"
	testEnv "goapi/testing"
	"log"
	"net/http"
	"sync"
	"testing"
)

type userTests struct {
	client *apiclient.CleverClient
}

func TestUser(t *testing.T) {
	ctx := context.Background()
	tunit, err := testEnv.NewUnit()
	if err != nil {
		log.Fatalf("Failed to run test: %s", err)
	}

	t.Cleanup(tunit.Teardown)

	tunit.RunApi(ctx)

	client := apiclient.BuildClever(tunit.ServerAddress)

	utests := userTests{
		client,
	}

	if err := utests.singupClient("userTesterAccount", "testpassword"); err != nil {
		log.Fatalf("Failed to singup: %s", err)
	}

	if err := utests.client.Login("userTesterAccount", "testpassword"); err != nil {
		log.Fatalf("Failed to set authorize: %s", err)
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
		t.Error(testEnv.ErrorLog(testGoalLog, err))
	}

	expected := http.StatusCreated
	switch {
	case httpCode != expected:
		t.Error(testEnv.FailedLog(testGoalLog, "httpStatus", expected, httpCode))
	case err != nil:
		t.Error(testEnv.ErrorLog(testGoalLog, err))
	case outusr.Name != cusr.Name:
		t.Error(testEnv.FailedLog(testGoalLog, "user.Name", cusr.Name, outusr.Name))
	default:
		t.Log(testEnv.SuccessLog(testGoalLog))
	}

	return outusr
}

// get single user by id
func (utests *userTests) getUser200(t *testing.T, usr user.User) {
	testGoalLog := "getUsers200: Should be able to get user by id."

	var outusr user.User

	httpCode, err := utests.client.Get("/users/"+usr.ID, &outusr)

	if err != nil {
		t.Error(testEnv.ErrorLog(testGoalLog, err))
	}

	expected := http.StatusOK
	switch {
	case httpCode != expected:
		t.Error(testEnv.FailedLog(testGoalLog, "httpStatus", expected, httpCode))
	case err != nil:
		t.Error(testEnv.ErrorLog(testGoalLog, err))
	case outusr.Name != usr.Name:
		t.Error(testEnv.FailedLog(testGoalLog, "user.Name", usr.Name, outusr.Name))
	default:
		t.Log(testEnv.SuccessLog(testGoalLog))
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
		t.Error(testEnv.ErrorLog(testGoalLog, err))
	}

	expected := http.StatusOK
	if httpCode != expected {
		t.Error(testEnv.FailedLog(testGoalLog, "httpStatus", expected, httpCode))
	}

	var outusr user.User
	httpCode, err = utests.client.Get("/users/"+uid, &outusr)

	switch {
	case err != nil:
		t.Error(testEnv.ErrorLog(testGoalLog, err))
	case inusr.Name != outusr.Name:
		t.Error(testEnv.FailedLog(testGoalLog, "user.Name", inusr.Name, outusr.Name))
	default:
		t.Log(testEnv.SuccessLog(testGoalLog))
	}
}

// delete user by id
func (utests *userTests) deleteUsers204(t *testing.T, uid string) {
	testGoalLog := "deleteUsers204: Should be able to delete user by id."

	httpCode, err := utests.client.Delete("/users/" + uid)
	if err != nil {
		t.Error(testEnv.ErrorLog(testGoalLog, err))
	}

	expected := http.StatusNoContent
	if httpCode != expected {
		t.Error(testEnv.FailedLog(testGoalLog, "httpStatus", expected, httpCode))
	}

	httpCode, _ = utests.client.Get("/users/"+uid, nil)

	expected = http.StatusNotFound
	switch {
	case httpCode != expected:
		t.Error(testEnv.FailedLog(testGoalLog, "httpStatus", expected, httpCode))
	default:
		t.Log(testEnv.SuccessLog(testGoalLog))
	}
}

// get multiple users by id
func (utests *userTests) getUsersList200(t *testing.T) {
	testGoalLog := "getUsers200: Should be able to get LIST of users by id."

	uamount := 22

	var wg sync.WaitGroup
	c := make(chan string, 1)

	// populate more users
	for i := 0; i < uamount; i++ {
		wg.Add(1)
		c <- fmt.Sprintf("name_%d", i)
		go func() {

			cusr := user.CreateUser{
				Name:            <-c,
				Password:        "testpassword",
				PasswordConfirm: "testpassword",
			}

			_, err := utests.client.Post("/users", cusr, nil)
			if err != nil {
				t.Errorf("concur: %v", err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	close(c)

	var outusrs []user.User
	httpCode, err := utests.client.Get("/users", &outusrs)

	expected := http.StatusOK
	switch {
	case httpCode != expected:
		t.Error(testEnv.FailedLog(testGoalLog, "httpStatus", expected, httpCode))
	case err != nil:
		t.Error(testEnv.ErrorLog(testGoalLog, err))
	case len(outusrs) != uamount+1:
		t.Error(testEnv.FailedLog(testGoalLog, "Amount of users", uamount+1, len(outusrs)))
	default:
		t.Log(testEnv.SuccessLog(testGoalLog))
	}
}
