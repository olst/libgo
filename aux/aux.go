package aux

import (
	"errors"
	"reflect"
)

// GetStructFields - get book fields
func GetStructFields(input interface{}) ([]string, error) {
	strctFields := make([]string, 1)
	strct := reflect.ValueOf(input).Elem()

	switch strct.Type().Kind() {
	case reflect.Struct:
		for i := 0; i < strct.NumField(); i++ {
			currentStructField := strct.Type().Field(i).Name
			strctFields = append(strctFields, currentStructField)
		}
		return strctFields, nil
	default:
		return strctFields, errors.New("You should use structs")
	}
}
