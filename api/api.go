package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/olst/libgo/aux"

	"sync"

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
	sync.RWMutex
	DbPath string
	Books  []Book
}

// BookIndex - get all books
func (ctx *AppCtx) BookIndex(w http.ResponseWriter, r *http.Request) {
	data, err := aux.GetFileData(ctx.DbPath)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	aux.WriteSuccess(w, http.StatusOK, data)
}

// GetBook - get book by id
func (ctx *AppCtx) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "No id specified", http.StatusBadRequest)
	}

	ctx.Lock()
	defer ctx.Unlock()

	data, err := aux.GetFileData(ctx.DbPath)
	if aux.CheckError(err) != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &ctx.Books)
	if aux.CheckError(err) != nil {
		return
	}

	book, idx := aux.GetBookByUUID(id, ctx.Books)
	if idx > 0 {
		aux.WriteSuccess(w, http.StatusOK, book)
	}
}

// AddBook - add a new book
func (ctx *AppCtx) AddBook(w http.ResponseWriter, r *http.Request) {
	ctx.Lock()
	defer ctx.Unlock()

	data, err := aux.GetFileData(ctx.DbPath)
	if aux.CheckError(err) != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &ctx.Books)
	if aux.CheckError(err) != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(r.Body)
	book := new(Book)
	err = decoder.Decode(book)
	if aux.CheckError(err) != nil {
		http.Error(w, "Error", http.StatusBadRequest)
		return
	}

	// Better to use some package
	uuid, err := exec.Command("uuidgen").Output()
	// Skipped error check
	aux.CheckError(err)

	stringUUID := string(uuid)
	book.ID = strings.TrimSuffix(stringUUID, "\n")

	ctx.Books = append(ctx.Books, *book)
	booksBytes, err := json.MarshalIndent(ctx.Books, "", "    ")
	if aux.CheckError(err) != nil {
		http.Error(w, "Error on save", http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile(ctx.DbPath, booksBytes, 0644)
	aux.CheckError(err)
}

// DeleteBook - delete book by id
func (ctx *AppCtx) DeleteBook(w http.ResponseWriter, r *http.Request) {
	data, err := aux.GetFileData(ctx.DbPath)
	if aux.CheckError(err) != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &ctx.Books)
	if aux.CheckError(err) != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	_, idx := aux.GetBookByUUID(id, ctx.Books)
	if idx <= 0 {
		aux.WriteError(w, http.StatusNotFound)
	}

	ctx.Books = append(ctx.Books[:idx], ctx.Books[idx+1:]...)
	booksBytes, err := json.MarshalIndent(ctx.Books, "", "    ")
	aux.CheckError(err)

	err = ioutil.WriteFile(ctx.DbPath, booksBytes, 0644)
	aux.CheckError(err)

	aux.WriteSuccess(w, http.StatusNoContent, nil)
}
