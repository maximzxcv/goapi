package integration

import (
	"fmt"
	"goapi/app/api/handlers"
	"goapi/business/apiclient"
	"goapi/business/auth"
	"goapi/business/data/call"
	"goapi/business/data/user"
	testEnv "goapi/testing"
	"log"
	"net/http"
	"sync"
	"testing"
)

type callTests struct {
	app    http.Handler
	client apiclient.Client
}

func TestCall(t *testing.T) {
	tunit, err := testEnv.NewUnit()
	if err != nil {
		log.Fatalf("Failed to run test: %s", err)
	}

	t.Cleanup(tunit.Teardown)

	app := handlers.API(tunit.Db)
	client, err := apiclient.BuildClient(app)

	if err != nil {
		log.Fatalf("Failed to run test: %s", err)
	}

	ctests := callTests{
		app,
		client,
	}

	if err := ctests.singupClient("callTesterAccount", "testpassword"); err != nil {
		log.Fatalf("Failed to singup: %s", err)
	}

	if err := ctests.client.Login("callTesterAccount", "testpassword"); err != nil {
		log.Fatalf("Failed to set authorize: %s", err)
	}

	t.Log("Call READ functionality")
	{
		ctests.getCallsList200(t)
	}
}

func (ctests *callTests) singupClient(u string, p string) error {
	signup := auth.Signup{
		Username: u,
		Password: p,
	}

	_, err := ctests.client.UnauthorizedCall(http.MethodPost, "/singup", &signup, nil)

	return err
}

// get multiple calls
func (ctests *callTests) getCallsList200(t *testing.T) {
	testGoalLog := "getCallsList200: Should be able to get LIST of calls for logged in user."

	uamount := 1

	_, _ = ctests.client.Get("/users", nil)

	var outusrs []call.Call
	httpCode, err := ctests.client.Get("/calls", &outusrs)

	expected := http.StatusOK
	switch {
	case httpCode != expected:
		t.Error(testEnv.FailedLog(testGoalLog, "httpStatus", expected, httpCode))
	case err != nil:
		t.Error(testEnv.ErrorLog(testGoalLog, err))
	case len(outusrs) != uamount:
		t.Error(testEnv.FailedLog(testGoalLog, "Amount of users", uamount, len(outusrs)))
	}

	uamount = 22

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

			_, err := ctests.client.Post("/users", cusr, nil)
			if err != nil {
				t.Errorf("concur: %v", err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
	close(c)

	httpCode, err = ctests.client.Get("/calls", &outusrs)

	switch {
	case err != nil:
		t.Error(testEnv.ErrorLog(testGoalLog, err))
	case len(outusrs) != uamount+2: // + /users and /calls
		t.Error(testEnv.FailedLog(testGoalLog, "Amount of users", uamount+2, len(outusrs)))
	default:
		t.Log(testEnv.SuccessLog(testGoalLog))
	}
}
