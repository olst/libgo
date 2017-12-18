package aux

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
)

// Auxiliary functions

// CheckError func
func CheckError(err error) {
	if err != nil {
		log.Print("Error: " + err.Error())
	}
}

// GetFileData - parse file
func GetFileData(path string) []byte {
	jsonFile, err := os.Open(path)
	CheckError(err)
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	CheckError(err)
	return bytes
}

// GetBookByUUID - get a book by UUID
func GetBookByUUID(uuid string, input interface{}) ([]byte, int) {
	switch reflect.TypeOf(input).Kind() {
	case reflect.Slice:
		books := reflect.ValueOf(input)
		for i := 0; i < books.Len(); i++ {
			v := reflect.Indirect(books.Index(i)).FieldByName("ID")
			if uuid == v.String() {
				b := books.Index(i).Interface()
				jsonBook, err := json.Marshal(b)
				CheckError(err)
				return jsonBook, i
			}
		}
	default:
		log.Println("Internal Server Error")
		return nil, -1
	}
	log.Printf("Book %s not found", uuid)
	return nil, -1
}

// WriteSuccess func
func WriteSuccess(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(data)
}

// WriteError func
func WriteError(w http.ResponseWriter, err int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(err)
	log.Print("Error: " + strconv.Itoa(err))
}
