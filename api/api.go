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
	DbPath string
	Books  []Book
}

// BookIndex - get all books
func (ctx *AppCtx) BookIndex(w http.ResponseWriter, r *http.Request) {
	aux.WriteSuccess(w, http.StatusOK, aux.GetFileData(ctx.DbPath))
}

// GetBook - get book by id
func (ctx *AppCtx) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := json.Unmarshal(aux.GetFileData(ctx.DbPath), &ctx.Books)
	aux.CheckError(err)
	book, idx := aux.GetBookByUUID(id, ctx.Books)
	if idx > 0 {
		aux.WriteSuccess(w, http.StatusOK, book)
	}
}

// AddBook - add a new book
func (ctx *AppCtx) AddBook(w http.ResponseWriter, r *http.Request) {
	err := json.Unmarshal(aux.GetFileData(ctx.DbPath), &ctx.Books)
	aux.CheckError(err)

	decoder := json.NewDecoder(r.Body)
	book := new(Book)
	err = decoder.Decode(&book)
	aux.CheckError(err)

	uuid, err := exec.Command("uuidgen").Output()
	aux.CheckError(err)

	stringUUID := string(uuid)
	book.ID = strings.TrimSuffix(stringUUID, "\n")

	ctx.Books = append(ctx.Books, *book)
	booksBytes, err := json.MarshalIndent(ctx.Books, "", "    ")
	aux.CheckError(err)

	err = ioutil.WriteFile(ctx.DbPath, booksBytes, 0644)
	aux.CheckError(err)
}

// DeleteBook - delete book by id
func (ctx *AppCtx) DeleteBook(w http.ResponseWriter, r *http.Request) {
	err := json.Unmarshal(aux.GetFileData(ctx.DbPath), &ctx.Books)
	aux.CheckError(err)

	vars := mux.Vars(r)
	id := vars["id"]

	_, idx := aux.GetBookByUUID(id, ctx.Books)
	if idx > 0 {
		ctx.Books = append(ctx.Books[:idx], ctx.Books[idx+1:]...)
		booksBytes, err := json.MarshalIndent(ctx.Books, "", "    ")
		aux.CheckError(err)

		err = ioutil.WriteFile(ctx.DbPath, booksBytes, 0644)
		aux.CheckError(err)

		aux.WriteSuccess(w, http.StatusNoContent, nil)
	} else {
		aux.WriteError(w, http.StatusNotFound)
	}

}
