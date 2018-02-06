package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/olst/libgo/api"
)

func main() {

	//jsonCtx := api.NewJSONctx()
	sqliteCtx := api.NewSqliteCtx("./library.db")
	defer sqliteCtx.CloseDB()

	router := mux.NewRouter()
	router.StrictSlash(false)

	//initHandlers(router, jsonCtx)
	initHandlers(router, sqliteCtx)
	log.Fatal(http.ListenAndServe(":8081", router))
}

type commonContext interface {
	BookIndex(w http.ResponseWriter, r *http.Request)
	GetBook(w http.ResponseWriter, r *http.Request)
	AddBook(w http.ResponseWriter, r *http.Request)
	DeleteBook(w http.ResponseWriter, r *http.Request)
	EditBook(w http.ResponseWriter, r *http.Request)
}

func initHandlers(router *mux.Router, context commonContext) {
	router.HandleFunc("/books/", context.BookIndex).Methods("GET")
	router.HandleFunc("/books/{id}", context.GetBook).Methods("GET")
	router.HandleFunc("/books/", AddBook).Methods("POST")
	router.HandleFunc("/books/{id}", context.DeleteBook).Methods("DELETE")
	router.HandleFunc("/books/{id}", context.EditBook).Methods("PUT")
}
