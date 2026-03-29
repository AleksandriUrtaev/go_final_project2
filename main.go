package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AleksandriUrtaev/go_final_project/pkg/api"
	"github.com/AleksandriUrtaev/go_final_project/pkg/db"
)

func main() {

	// БД
	database, err := db.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	port := "7540"
	if p := os.Getenv("TODO_PORT"); p != "" {
		port = p
	}

	addr := ":" + port

	//
	mux := http.NewServeMux()

	//API
	api.Init(mux)

	db.SetDB(database)
	//
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)

	//
	log.Println("Server started at http://localhost" + addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
