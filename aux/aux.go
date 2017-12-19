package aux

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
)

// Auxiliary functions

// CheckError func
func CheckError(err error) error {
	if err != nil {
		log.Print("Error:", err)
	}
	return err
}

// GetFileData - parse file
func GetFileData(path string) ([]byte, error) {
	jsonFile, err := os.Open(path)
	if CheckError(err) != nil {
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

// GetBookByUUID - get a book by UUID
func GetBookByUUID(uuid string, input interface{}) ([]byte, int) {
	inputValRef := reflect.ValueOf(input)
	switch inputValRef.Type().Kind() {
	case reflect.Slice:
		books := reflect.ValueOf(input)
		for i := 0; i < books.Len(); i++ {
			v := reflect.Indirect(books.Index(i)).FieldByName("ID")

			if !v.IsNil() && v.IsValid() && uuid == v.String() {
				b := books.Index(i).Interface()
				jsonBook, err := json.Marshal(b)
				CheckError(err)
				return jsonBook, i
			}
		}

		log.Printf("Book %s not found", uuid)
		return nil, -1
	}

	log.Println("Not slice!!!")
	return nil, -1
}

// WriteSuccess func
func WriteSuccess(w http.ResponseWriter, status int, data []byte) {
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

// WriteError func
func WriteError(w http.ResponseWriter, err int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(err)
	log.Print("Error:", err)
}
