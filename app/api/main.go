package main

import (
	"fmt"
	"goapi/app/api/handlers"
	"goapi/app/api/middle"
	"goapi/bal"
	"goapi/business/data/user"
	"goapi/business/mid"
	"log"
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
	"github.com/spf13/viper"
)

const (
	serverAddress = "127.0.0.1:8080"
)

func main() {
	logg := bal.NewLogg()
	config, err := mid.LoadConfig(".")
	if err != nil {
		log.Print(err)
	}
	log.Print(config)
	//	log.Println("main: Initializing database support")
	db, err := open()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer func() {
		log.Printf("main: Database Stopping : %s", db.DriverName())
		//	db.Close()
	}()

	router := httptreemux.NewContextMux()

	urep := user.NewRepository(db)
	uh := handlers.NewUserHandler(urep)
	router.Handler(http.MethodGet, "/users/:id", appHandler(uh.GetUserByID))
	router.Handler(http.MethodPut, "/users/:id", appHandler(uh.UpdateUser))
	router.Handler(http.MethodDelete, "/users/:id", appHandler(uh.DeleteUser))
	router.Handler(http.MethodGet, "/users", appHandler(uh.GetUsers))
	router.Handler(http.MethodPost, "/users", appHandler(uh.CreateUser))

	loggMiddle := middle.LoggMiddle(logg)
	cnfgrdRouter := loggMiddle(router)

	api := http.Server{
		Addr:    serverAddress,
		Handler: cnfgrdRouter,
	}

	log.Printf("API is running on %v", serverAddress)
	if err := api.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}

func open() (*sqlx.DB, error) {

	conf := func(key string) string {
		return viper.GetString(`database.` + key)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf("host"), conf("port"), conf("user"), conf("pass"), conf("dbname"))
	fmt.Println(psqlInfo)
	return sqlx.Open("postgres", psqlInfo)
}

// AppHandler .....
type appHandler func(http.ResponseWriter, *http.Request) *handlers.ErrorResponse

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if erresp := fn(w, r); erresp != nil {
		log.Printf("%+v", erresp)
		w.WriteHeader(erresp.Code)
	}
}
