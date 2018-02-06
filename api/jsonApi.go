package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"sync"

	"github.com/olst/libgo/model"
	uuid "github.com/satori/go.uuid"

	"github.com/gorilla/mux"
)

// JSONctx - application JSON context
type JSONctx struct {
	sync.RWMutex
	DbPath string
	Books  []model.Book
}

// NewJSONctx - creates a new JSON context
func NewJSONctx() *JSONctx {
	c := new(JSONctx)
	c.DbPath = "./library.json"
	c.Books = make([]model.Book, 1)
	return c
}

func (jsonCtx *JSONctx) Close() error { return nil }

// BookIndex - get all books
func (jsonCtx *JSONctx) BookIndex(w http.ResponseWriter, r *http.Request) {
	data, err := getFileData(jsonCtx.DbPath)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	writeSuccess(w, http.StatusOK, data)
}

// GetBook - get book by id
func (jsonCtx *JSONctx) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "No id specified", http.StatusBadRequest)
	}

	jsonCtx.RLock()
	defer jsonCtx.RUnlock()

	data, err := getFileData(jsonCtx.DbPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &jsonCtx.Books)
	if err != nil {
		return
	}

	book, err := getBookByUUIDjson(id, jsonCtx.Books)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeSuccess(w, http.StatusOK, book)
}

// AddBook - add a new book
func (jsonCtx *JSONctx) AddBook(w http.ResponseWriter, r *http.Request) {
	jsonCtx.Lock()
	defer jsonCtx.Unlock()

	data, err := getFileData(jsonCtx.DbPath)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &jsonCtx.Books)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(r.Body)
	book := new(model.Book)
	err = decoder.Decode(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	book.ID = uuid.NewV4().String()

	jsonCtx.Books = append(jsonCtx.Books, *book)
	booksBytes, err := json.MarshalIndent(jsonCtx.Books, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile(jsonCtx.DbPath, booksBytes, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteBook - delete book by id
func (jsonCtx *JSONctx) DeleteBook(w http.ResponseWriter, r *http.Request) {
	jsonCtx.Lock()
	defer jsonCtx.Unlock()

	data, err := getFileData(jsonCtx.DbPath)
	if err != nil {
		http.Error(w, "Error: Couldn't open a file",
			http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &jsonCtx.Books)
	if err != nil {
		http.Error(w, "Error: Couldn't get books",
			http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Error getting ID",
			http.StatusInternalServerError)
		return
	}

	idx, err := getBookIndex(id, jsonCtx.Books)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		log.Print("Error:", err)
	}

	jsonCtx.Books = append(jsonCtx.Books[:idx], jsonCtx.Books[idx+1:]...)
	booksBytes, err := json.MarshalIndent(jsonCtx.Books, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile(jsonCtx.DbPath, booksBytes, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeSuccess(w, http.StatusNoContent, nil)
}

// EditBook - edit a book by id
func (jsonCtx *JSONctx) EditBook(w http.ResponseWriter, r *http.Request) {
	jsonCtx.Lock()
	defer jsonCtx.Unlock()

	data, err := getFileData(jsonCtx.DbPath)
	if err != nil {
		http.Error(w, "Error: Couldn't open a file",
			http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &jsonCtx.Books)
	if err != nil {
		http.Error(w, "Error: Couldn't get books",
			http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Error getting ID",
			http.StatusInternalServerError)
		return
	}

	idx, err := getBookIndex(id, jsonCtx.Books)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		log.Print("Error:", err)
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&jsonCtx.Books[idx])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	booksBytes, err := json.MarshalIndent(jsonCtx.Books, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile(jsonCtx.DbPath, booksBytes, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeSuccess(w, http.StatusNoContent, nil)
}

func getFileData(path string) ([]byte, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := jsonFile.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	return ioutil.ReadAll(jsonFile)
}

func getBookByUUIDjson(uuid string, input interface{}) ([]byte, error) {
	inputValRef := reflect.ValueOf(input)
	switch inputValRef.Type().Kind() {
	case reflect.Slice:
		books := reflect.ValueOf(input)
		for i := 0; i < books.Len(); i++ {
			currentBookUUID := reflect.Indirect(books.Index(i)).FieldByName("ID")
			if uuid == currentBookUUID.String() {
				b := books.Index(i).Interface()
				jsonBook, err := json.Marshal(b)
				if err != nil {
					return nil, err
				}
				return jsonBook, nil
			}

		}
		log.Printf("Book %s not found", uuid)
		return nil, nil
	}
	return nil, errors.New("you should use slices")
}

func getBookIndex(uuid string, input interface{}) (int, error) {
	inputValRef := reflect.ValueOf(input)
	if inputValRef.Type().Kind() != reflect.Slice {
		return -1, errors.New("you should use slices")
	}

	for i := 0; i < inputValRef.Len(); i++ {
		currentBookUUID := reflect.Indirect(inputValRef.Index(i)).FieldByName("ID")
		if uuid == currentBookUUID.String() {
			return i, nil
		}
	}
	return -1, nil
}

func writeSuccess(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	n, err := w.Write(data)
	if err != nil {
		log.Print(err)
	}
	if len(data) != n {
		log.Print("WAAAAT!?")
	}
}
