package user

import (
	"context"
	"goapi/foundation/dbase"
	testEnv "goapi/testing"
	"log"
	"strconv"
	"testing"

	"github.com/pkg/errors"
)

func TestUser(t *testing.T) {
	tunit, err := testEnv.NewUnit()
	if err != nil {
		log.Fatalf("Failed to run test: %s", err)
	}

	t.Cleanup(tunit.Teardown)

	rep := NewRepository(tunit.Db)
	ctx := context.Background()
	t.Log("User CRUD functionality")
	{
		testGoalLog := "Create: Should be able to create a user."
		curs := CreateUser{
			Name:     "Name",
			Password: "Password",
		}
		var fusr User
		if fusr, err = rep.Create(ctx, curs); err != nil {
			t.Error(testEnv.ErrorLog(testGoalLog, err))
		} else {
			t.Logf(testEnv.SuccessLog(testGoalLog))
		}

		// populate more users
		for i := 1; i < 5; i++ {
			curs.Name = curs.Name + strconv.Itoa(i)
			rep.Create(ctx, curs)
		}

		testGoalLog = "QueryByID: Should be able to query single user by Id."
		usr, err := rep.QueryByID(ctx, fusr.ID)
		if err != nil {
			t.Error(testEnv.ErrorLog(testGoalLog, err))
		}
		if usr.Name != "Name" {
			t.Error(testEnv.FailedLog(testGoalLog, "Name", usr.Name, "Name"))
		} else {
			t.Logf(testEnv.SuccessLog(testGoalLog))
		}

		testGoalLog = "Update: Should be able to update user by Id."
		var uusr = UpdateUser{
			Name: curs.Name,
		}
		uusr.Name = "updatedName"
		usr, err = rep.Update(ctx, fusr.ID, uusr)
		if err != nil {
			t.Error(testEnv.ErrorLog(testGoalLog, err))
		}
		if usr.Name != "updatedName" {
			t.Error(testEnv.FailedLog(testGoalLog, "Name", usr.Name, "updatedName"))
		} else {
			t.Logf(testEnv.SuccessLog(testGoalLog))
		}

		testGoalLog = "Query: Should be able to query multipleusers."
		usrs, err := rep.Query(ctx)
		if err != nil {
			t.Error(testEnv.ErrorLog(testGoalLog, err))
		}
		if len(usrs) != 5 {
			t.Error(testEnv.FailedLog(testGoalLog, "len(usrs)", 5, len(usrs)))
		} else {
			t.Logf(testEnv.SuccessLog(testGoalLog))
		}

		testGoalLog = "Delete: Should be able to delete user."
		if err := rep.Delete(ctx, fusr.ID); err != nil {
			t.Error(testEnv.ErrorLog(testGoalLog, err))
		}
		_, err = rep.QueryByID(ctx, fusr.ID)
		if errors.Cause(err) != dbase.ErrNotExist {
			t.Error(testEnv.ErrorLog(testGoalLog, err))
		} else {
			t.Logf(testEnv.SuccessLog(testGoalLog))
		}
	}
}
