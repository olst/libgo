package aux

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Auxiliary functions

// CheckError func
func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
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

// WriteSuccess func
func WriteSuccess(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// WriteError func
func WriteError(w http.ResponseWriter, err int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(err)
}
