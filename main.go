package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/scheduler-prototype/handler"
	"github.com/scheduler-prototype/mgraph"
	"github.com/scheduler-prototype/repository"
)

func main() {
	// loading env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// db connection
	dbConnStr := os.Getenv("DB_CONN_STR")
	dbDriver := os.Getenv("DB_DRIVER")
	db, err := sql.Open(dbDriver, dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to database!")

	// initialize msgraph client
	client, err := mgraph.NewMGraphClient()
	if err != nil {
		log.Fatal(err)
	}

	// initialize reposotories
	repo := repository.NewRepository(db)

	// initialize handlers
	controller := handler.NewHandler(client, repo)

	// chi router
	r := chi.NewRouter()

	subRouter := chi.NewRouter()
	subRouter.Get("/calendarview", controller.MGraphGetCalendarView)
	subRouter.Post("/event/create", controller.MGraphCreateEvent)
	subRouter.Post("/calendarview/first-sync", controller.MGraphCalendarViewFirstSync)
	subRouter.Post("/calendarview/subscription", controller.MGraphCreateCalendarViewSubscription)
	subRouter.Post("/calendarview/subscription/verify", controller.MGraphVerifyCalendarViewSubscription)

	r.Mount("/mgraph", subRouter)

	http.ListenAndServe(":8080", r)
}
