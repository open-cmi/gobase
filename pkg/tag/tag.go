package tag

import (
	"reflect"
)

func GetStructTagColumn(v interface{}, tag string) []string {
	var columns []string = []string{}
	dataValue := reflect.ValueOf(v)
	if dataValue.Kind() != reflect.Struct {
		return columns
	}

	t := dataValue.Type()
	for i := 0; i < t.NumField(); i++ {
		field := dataValue.Type().Field(i)
		tagColumn, ok := field.Tag.Lookup(tag)
		if ok {
			columns = append(columns, tagColumn)
		}
	}

	return columns
}
