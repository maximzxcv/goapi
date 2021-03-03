package main

import (
	"goapi/app/api/handlers"
	"goapi/app/api/middle"
	"goapi/business/mid"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
)

const (
	serverAddress = "127.0.0.1:8080"
)

func main() {
	if err := mid.LoadConfig("."); err != nil {
		log.Fatalf("Could configfile: %s", err)
	}
	dbConfig, _ := mid.GetDbConfig()
	log.Printf("config.json is loaded")

	//	log.Println("main: Initializing database support")
	db, err := open(dbConfig)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer func() {
		log.Printf("main: Database Stopping : %s", db.DriverName())
		//	db.Close()
	}()

	api := http.Server{
		Addr:    serverAddress,
		Handler: handlers.API(db, middle.LoggMiddle()),
	}

	log.Printf("API is running on %v", serverAddress)
	if err := api.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}

//TODO : duplicated code
func open(dbConfig *mid.DbConfig) (*sqlx.DB, error) {

	return sqlx.Open("postgres", dbConfig.ConnectinString())
}
