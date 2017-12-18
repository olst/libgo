package main

import (
	"log"
	"net/http"
	"ostyhar/library/api"

	"github.com/gorilla/mux"
)

func main() {
	ctx := api.AppCtx{
		books:  make([]api.Book, 1),
		dbPath: "./library.json",
	}
	router := mux.NewRouter()
	router.HandleFunc("/books", ctx.BookIndex).Methods("GET")
	router.HandleFunc("/book/{id}", ctx.GetBook).Methods("GET")
	router.HandleFunc("/books", ctx.AddBook).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", router))
}
