package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

func main() {
	// db connection
	db, err := sql.Open("postgres", "postgres://root:secret@localhost:5434/prototype?sslmode=disable")
	if err != nil { 
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil { 
		panic(err.Error())
	}

	fmt.Println("Successfully connected to database!")

	// chi router 
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
