package ttesting

import (
	"fmt"
	"goapi/business/mid"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	Success = "\u2713"
	Failed  = "\u2717"
)

// TestUnit is an environment to run test
type TestUnit struct {
	dRunner *DockerRunner
	Db      *sqlx.DB
}

// NewUnit constructor
func NewUnit() (*TestUnit, error) {
	dbConfig := mid.NewTestConfig()
	log.Print("DB connection: ", dbConfig.ConnectinString())

	dRunner := NewDockerRunner(dbConfig)

	if err := dRunner.Start(); err != nil {
		return nil, errors.Wrap(err, "Could not run docker: %s")
	}

	db, err := open(dbConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to DB on docker: %s")
	}

	if err := dRunner.CreateDbSchema(db); err != nil {
		return nil, errors.Wrap(err, "Could not create DB schema on docker: %s")
	}

	return &TestUnit{
		dRunner: dRunner,
		Db:      db,
	}, nil

}

// Teardown remove container
func (tunit *TestUnit) Teardown() {
	if err := tunit.dRunner.Stop(); err != nil {
		log.Fatalf("Failed to stop container: %s", err)
	}
}

//TODO : duplicated code
func open(dbConfig *mid.DbConfig) (*sqlx.DB, error) {
	return sqlx.Open("postgres", dbConfig.ConnectinString())
}

type failedResult struct {
	test     string
	field    string
	expected interface{}
	actual   interface{}
}

func FailedLog(test string, field string, expected interface{}, actual interface{}) string {
	return fmt.Sprintf("\t%s\tTest :\t%s: %s is not as expected. (ex/act -> %v/%v)", Failed, test, field, expected, actual)
}

func ErrorLog(test string, err error) string {
	return fmt.Sprintf("\t%s\tTest :\t%s: %s.", Failed, test, err)
}

func SuccessLog(test string) string {
	return fmt.Sprintf("\t%s\tTest :\t%s", Success, test)
}
