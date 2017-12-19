package main

import (
	"log"
	"net/http"

	"github.com/olst/libgo/api"

	"github.com/gorilla/mux"
)

func main() {
	ctx := api.AppCtx{
		Books:  make([]api.Book, 1),
		DbPath: "./library.json",
	}
	router := mux.NewRouter()
	router.HandleFunc("/books", ctx.BookIndex).Methods("GET")
	router.HandleFunc("/books/{id}", ctx.GetBook).Methods("GET")
	router.HandleFunc("/books", ctx.AddBook).Methods("POST")
	router.HandleFunc("/books/{id}", ctx.DeleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8081", router))
}
