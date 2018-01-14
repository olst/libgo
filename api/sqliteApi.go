package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/olst/libgo/model"
	uuid "github.com/satori/go.uuid"
)

var tableName = "books_table"

// SqliteCtx - application sqlite context
type SqliteCtx struct {
	sync.RWMutex
	db     *sql.DB
	dbName string
}

// CloseDB ...
func (sqliteCtx *SqliteCtx) CloseDB() {
	sqliteCtx.db.Close()
}

// NewSqliteCtx - creates a new sqlite context
func NewSqliteCtx(dbName string) *SqliteCtx {
	c := new(SqliteCtx)
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		panic(err)
	}
	c.db = db
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s "+
		"(id TEXT PRIMARY KEY, title TEXT, genres TEXT,"+
		"pages INT, price REAL)", tableName)
	statement, err := c.db.Prepare(query)
	if err != nil {
		panic(err)
	}
	statement.Exec()
	return c
}

// BookIndex - get all books
func (sqliteCtx *SqliteCtx) BookIndex(w http.ResponseWriter, r *http.Request) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := sqliteCtx.db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	book := new(model.Book)

	for rows.Next() {
		err = rows.Scan(&book.ID, &book.Title, &book.Genres,
			&book.Pages, &book.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//jsonBook, err := json.Marshal(book)
		jsonBook, err := json.MarshalIndent(book, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(jsonBook)
	}
}

func getBookByUUIDsql(sqliteCtx *SqliteCtx, w http.ResponseWriter, r *http.Request) []byte {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "No id specified", http.StatusBadRequest)
		return nil
	}
	query := fmt.Sprintf("SELECT * from %s WHERE id='%s'", tableName, id)
	book := new(model.Book)
	err := sqliteCtx.db.QueryRow(query).Scan(&book.ID, &book.Title,
		&book.Genres, &book.Pages, &book.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	jsonBook, err := json.MarshalIndent(book, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return jsonBook
}

// GetBook by ID
func (sqliteCtx *SqliteCtx) GetBook(w http.ResponseWriter, r *http.Request) {
	sqliteCtx.Lock()
	defer sqliteCtx.Unlock()
	jsonBook := getBookByUUIDsql(sqliteCtx, w, r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonBook)
}

// AddBook ...
func (sqliteCtx *SqliteCtx) AddBook(w http.ResponseWriter, r *http.Request) {
	sqliteCtx.Lock()
	defer sqliteCtx.Unlock()
	decoder := json.NewDecoder(r.Body)
	book := new(model.Book)
	err := decoder.Decode(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	book.ID = uuid.NewV4().String()

	query := fmt.Sprintf("INSERT INTO %s"+
		"(id, title, genres, pages, price) values(?,?,?,?,?)", tableName)

	statement, err := sqliteCtx.db.Prepare(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = statement.Exec(&book.ID, &book.Title, &book.Genres,
		&book.Pages, &book.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonBook, err := json.MarshalIndent(book, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonBook)
}

// DeleteBook ...
func (sqliteCtx *SqliteCtx) DeleteBook(w http.ResponseWriter, r *http.Request) {
	sqliteCtx.Lock()
	defer sqliteCtx.Unlock()
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Error getting ID",
			http.StatusInternalServerError)
		return
	}

	jsonBook := getBookByUUIDsql(sqliteCtx, w, r)

	query := fmt.Sprintf("DELETE from %s WHERE id='%s'", tableName, id)
	statement, err := sqliteCtx.db.Prepare(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = statement.Exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonBook)
}

// EditBook ...
func (sqliteCtx *SqliteCtx) EditBook(w http.ResponseWriter, r *http.Request) {
	sqliteCtx.Lock()
	defer sqliteCtx.Unlock()
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Error getting ID",
			http.StatusInternalServerError)
		return
	}
	query := fmt.Sprintf("SELECT * from %s WHERE id='%s'", tableName, id)
	book := new(model.Book)
	err := sqliteCtx.db.QueryRow(query).Scan(&book.ID, &book.Title,
		&book.Genres, &book.Pages, &book.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query = fmt.Sprintf("UPDATE %s SET title='%s', genres='%s', pages=%d, "+
		"price=%f WHERE id='%s'", tableName, book.Title, book.Genres, book.Pages,
		book.Price, id)

	statement, err := sqliteCtx.db.Prepare(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = statement.Exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonBook, err := json.MarshalIndent(book, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(jsonBook)
}
