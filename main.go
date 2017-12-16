package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
)

// Book - represents a book
type Book struct {
	ID     string   `json:"id,omitempty"`
	Title  string   `json:"title,omitempty"`
	Ganres []string `json:"ganres,omitempty"`
	Pages  int      `json:"pages,omitempty"`
	Price  float32  `json:"price,omitempty"`
}

// AppCtx - application context
type AppCtx struct {
	books []Book
}

func main() {
	ctx := AppCtx{books: make([]Book, 1)}
	router := mux.NewRouter()
	router.HandleFunc("/books", BookIndex).Methods("GET")
	router.HandleFunc("/book/{id}", ctx.GetBook).Methods("GET")
	router.HandleFunc("/books", ctx.AddBook).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", router))
}

// BookIndex - get all books
func BookIndex(w http.ResponseWriter, r *http.Request) {
	bytes := getFileData("./library.json")
	writeSuccess(w, bytes)
}

// GetBook - get book by id
func (ctx *AppCtx) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	bytes := getFileData("./library.json")
	err := json.Unmarshal(bytes, &ctx.books)
	checkError(err)
	for _, book := range ctx.books {
		if id == book.ID {
			b := book
			jsonBook, err := json.Marshal(b)
			checkError(err)
			writeSuccess(w, jsonBook)
			return
		}
	}
	writeError(w, http.StatusNotFound)
}

// AddBook - add a new book
func (ctx *AppCtx) AddBook(w http.ResponseWriter, r *http.Request) {
	bytes := getFileData("./library.json")
	err := json.Unmarshal(bytes, &ctx.books)
	checkError(err)

	decoder := json.NewDecoder(r.Body)
	book := new(Book)
	err = decoder.Decode(&book)
	checkError(err)

	uuid, err := exec.Command("uuidgen").Output()
	checkError(err)

	stringUUID := string(uuid)
	book.ID = strings.TrimSuffix(stringUUID, "\n")

	ctx.books = append(ctx.books, *book)
	booksBytes, err := json.MarshalIndent(ctx.books, "", "    ")
	checkError(err)

	err = ioutil.WriteFile("./library.json", booksBytes, 0644)
}

// Auxiliary functions

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func getFileData(path string) []byte {
	jsonFile, err := os.Open(path)
	checkError(err)
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	checkError(err)
	return bytes
}

func writeSuccess(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func writeError(w http.ResponseWriter, err int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(err)
}
