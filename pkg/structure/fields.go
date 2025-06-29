package structure

import (
	"reflect"
)

func GetStructFields(v interface{}, tag string) []string {
	var fields []string = []string{}
	dataValue := reflect.ValueOf(v)
	if dataValue.Kind() != reflect.Struct {
		return fields
	}

	t := dataValue.Type()
	for i := 0; i < t.NumField(); i++ {
		field := dataValue.Type().Field(i)
		tagColumn, ok := field.Tag.Lookup(tag)
		if ok {
			fields = append(fields, tagColumn)
		}
	}

	return fields
}
