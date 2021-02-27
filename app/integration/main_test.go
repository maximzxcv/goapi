package integration

import (
	"fmt"
	"goapi/business/mid"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
	"github.com/spf13/viper"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	// Create a new pool for Docker containers

	pool, err := dockertest.NewPool("") // pool is the place to run container
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	log.Println("pool", pool)

	// Pull an image, create a container based on it and set all necessary parameters
	// opts := dockertest.RunOptions{
	// 	Repository:   "mdillon/postgis",
	// 	Tag:          "latest",
	// 	Env:          []string{"POSTGRES_PASSWORD=goapitestpass"},
	// 	ExposedPorts: []string{"5432"},
	// 	PortBindings: map[docker.Port][]docker.PortBinding{
	// 		"5432": {
	// 			{HostIP: "0.0.0.0", HostPort: "5433"},
	// 		},
	// 	},
	// }
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=goapitestpass",
			"listen_addresses = '*'",
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: "5433"},
			},
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not run in Docker: %s", err)
	}

	log.Println("resource", resource)

	// resource, err := pool.RunWithOptions(&opts)
	// if err != nil {
	// 	log.Fatalf("Could not start resource: %s", err)
	// }

	config, err := mid.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could configfile: %s", err)
	}

	log.Print(config)

	// Exponential retry to connect to database while it is booting
	if err := pool.Retry(func() error {
		//databaseConnStr := fmt.Sprintf("host=localhost port=5433 user=postgres dbname=postgres password=goapitestpass sslmode=disable")
		psqlInfo := fmt.Sprintf("host=localhost port=5433 user=postgres password=goapitestpass dbname=postgres sslmode=disable")
		db, err := sqlx.Open("postgres", psqlInfo)
		if err != nil {
			log.Println("Database not ready yet (it is booting up, wait for a few tries)...")
			return err
		}

		log.Println("ping", db.Ping())
		// Tests if database is reachable
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	/*  "host": "localhost",
	    "port": "5432",
	    "user": "postgres",
	    "pass": "0819",
	    "dbname": "goapi"*/
	// Run the Docker container

	// Delete the docker container
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(0)
}

//TODO : duplicated code
func open() (*sqlx.DB, error) {

	conf := func(key string) string {
		return viper.GetString(`database.` + key)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf("host"), conf("port"), conf("user"), conf("pass"), conf("dbname"))
	fmt.Println(psqlInfo)
	return sqlx.Open("postgres", psqlInfo)
}

// func TestMain(m *testing.M) {
// 	pool, resource := initDB()
// 	code := m.Run()
// 	closeDB(pool, resource)
// 	os.Exit(code)
// }

// func initDB() (*dockertest.Pool, *dockertest.Resource) {
// 	pgURL := initPostgres()
// 	pgPass, _ := pgURL.User.Password()

// 	runOpts := dockertest.RunOptions{
// 		Repository: "postgres",
// 		Tag:        "latest",
// 		Env: []string{
// 			"POSTGRES_USER=" + pgURL.User.Username(),
// 			"POSTGRES_PASSWORD=" + pgPass,
// 			"POSTGRES_DB=" + pgURL.Path,
// 		},
// 	}

// 	pool, err := dockertest.NewPool("")
// 	if err != nil {
// 		log.WithError(err).Fatal("Could not connect to docker")
// 	}

// 	resource, err := pool.RunWithOptions(&runOpts)
// 	if err != nil {
// 		log.WithError(err).Fatal("Could start postgres container")
// 	}

// 	pgURL.Host = resource.Container.NetworkSettings.IPAddress

// 	// Docker layer network is different on Mac
// 	// if runtime.GOOS == "darwin" {
// 	// 	pgURL.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
// 	// }

// 	DockerDBConn = &dockerDBConn{}
// 	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
// 	if err := pool.Retry(func() error {
// 		DockerDBConn.Conn, err = sql.Open("postgres", pgURL.String())
// 		if err != nil {
// 			return err
// 		}
// 		return DockerDBConn.Conn.Ping()
// 	}); err != nil {
// 		phrase := fmt.Sprintf("Could not connect to docker: %s", err)
// 		log.Error(phrase)
// 	}

// 	DockerDBConn.initMigrations()

// 	return pool, resource
// }

// func closeDB(pool *dockertest.Pool, resource *dockertest.Resource) {
// 	if err := pool.Purge(resource); err != nil {
// 		phrase := fmt.Sprintf("Could not purge resource: %s", err)
// 		log.Error(phrase)
// 	}
// }

// func (db dockerDBConn) initMigrations() {
// 	driver, err := postgres.WithInstance(db.Conn, &postgres.Config{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	migrate, err := migrate.NewWithDatabaseInstance(
// 		"file://testdata/migrations/",
// 		"mydatabase", driver)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	err = migrate.Steps(2)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func initPostgres() *url.URL {
// 	pgURL := &url.URL{
// 		Scheme: "postgres",
// 		User:   url.UserPassword("postgres-test", "0819-test"),
// 		Path:   "goapi-test",
// 	}
// 	q := pgURL.Query()
// 	q.Add("sslmode", "disable")
// 	pgURL.RawQuery = q.Encode()

// 	return pgURL
// }
