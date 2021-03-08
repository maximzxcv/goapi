package testing

import (
	"goapi/business/mid"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/pkg/errors"
)

type DockerRunner struct {
	config    *mid.DbConfig
	pool      *dockertest.Pool
	container *dockertest.Resource
}

// NewDockerRunner ....
func NewDockerRunner(config *mid.DbConfig) *DockerRunner {
	return &DockerRunner{
		config: config,
	}
}

// Start prepare environment to run tests
func (runner *DockerRunner) Start() error {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return errors.Wrap(err, "Could not connect to Docker: %s")
	}
	runner.pool = pool

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_USER=" + runner.config.User,
			"POSTGRES_PASSWORD=" + runner.config.Password,
			"listen_addresses = '*'",
		},
		ExposedPorts: []string{"5432"}, //
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: runner.config.Port}, // "5433"},
			},
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		errors.Wrap(err, "Could not run in Docker: %s")
	}

	runner.container = resource

	log.Printf("Postgres (DB) is running on container:%s", resource.Container.Name)
	return nil
}

// CreateDbSchema create all tables (+ migration) and data if required
func (runner *DockerRunner) CreateDbSchema(db *sqlx.DB) error {
	if err := runner.pool.Retry(func() error {
		ping := db.Ping()
		log.Println("DB ping", ping == nil)
		return ping
	}); err != nil {
		return errors.Wrap(err, "Could not connect to Docker: %s")
	}

	Migrate(db)
	log.Println("DB schema ready to use")
	return nil
}

// Stop - Delete the docker container
func (runner *DockerRunner) Stop() error {
	if err := runner.pool.Purge(runner.container); err != nil {
		errors.Wrap(err, "Could not purge resource: %s")
	}
	log.Printf("Container:%s is deleted.", runner.container.Container.Name)
	return nil
}
