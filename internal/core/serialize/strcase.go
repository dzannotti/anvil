package utils

import (
	"reflect"

	"github.com/iancoleman/strcase"
)

func ToSnakeCaseMap(obj any) map[string]any {
	result := make(map[string]any)
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		snakeCaseKey := strcase.ToSnake(field.Name)
		result[snakeCaseKey] = value
	}

	return result
}
