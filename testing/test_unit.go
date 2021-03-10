package testing

import (
	"context"
	"fmt"
	"goapi/app/api/handlers"
	"goapi/business/mid"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	Success       = "\u2713"
	Failed        = "\u2717"
	ServerAddress = "127.0.0.1:8080"
)

// TestUnit is an environment to run test
type TestUnit struct {
	dRunner       *DockerRunner
	Db            *sqlx.DB
	ServerAddress string
	api           *http.Server
	stopAPI       func()
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

	api := http.Server{
		Addr:    ServerAddress,
		Handler: handlers.API(db), //, middle.LoggMiddle(), middle.CallMiddle()),
	}

	go func() {

	}()

	return &TestUnit{
		dRunner:       dRunner,
		Db:            db,
		ServerAddress: ServerAddress,
		api:           &api,
	}, nil

}

func (tunit *TestUnit) RunApi(ctx context.Context) { //} (wg *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(ctx)
	tunit.stopAPI = cancel
	go func() {
		select {
		case <-ctx.Done():
			log.Printf("API is stopping on %v", ServerAddress)
			if err := tunit.api.Shutdown(context.Background()); err != nil {
				log.Fatalf("Failed to stop API server: %s", err)
			}
		}
	}()
	go func() {
		log.Printf("API is running on %v", ServerAddress)
		if err := tunit.api.ListenAndServe(); err != nil {
			if errors.Cause(err) != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()
}

// Teardown remove container
func (tunit *TestUnit) Teardown() {
	if err := tunit.dRunner.Stop(); err != nil {
		log.Fatalf("Failed to stop container: %s", err)
	}
	if tunit.stopAPI != nil {
		tunit.stopAPI()
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
