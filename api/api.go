package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"ostyhar/library/aux"
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
	dbData []byte
	dbPath string
	books  []Book
}

// BookIndex - get all books
func (ctx *AppCtx) BookIndex(w http.ResponseWriter, r *http.Request) {
	aux.WriteSuccess(w, aux.GetFileData(ctx.dbPath))
}

// GetBook - get book by id
func (ctx *AppCtx) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := json.Unmarshal(aux.GetFileData(ctx.dbPath), &ctx.books)
	aux.CheckError(err)
	for _, book := range ctx.books {
		if id == book.ID {
			b := book
			jsonBook, err := json.Marshal(b)
			aux.CheckError(err)
			aux.WriteSuccess(w, jsonBook)
			return
		}
	}
	aux.WriteError(w, http.StatusNotFound)
}

// AddBook - add a new book
func (ctx *AppCtx) AddBook(w http.ResponseWriter, r *http.Request) {
	err := json.Unmarshal(aux.GetFileData(ctx.dbPath), &ctx.books)
	aux.CheckError(err)

	decoder := json.NewDecoder(r.Body)
	book := new(Book)
	err = decoder.Decode(&book)
	aux.CheckError(err)

	uuid, err := exec.Command("uuidgen").Output()
	aux.CheckError(err)

	stringUUID := string(uuid)
	book.ID = strings.TrimSuffix(stringUUID, "\n")

	ctx.books = append(ctx.books, *book)
	booksBytes, err := json.MarshalIndent(ctx.books, "", "    ")
	aux.CheckError(err)

	err = ioutil.WriteFile(ctx.dbPath, booksBytes, 0644)
}
