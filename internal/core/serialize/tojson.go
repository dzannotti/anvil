package serialize

import (
	"encoding/json"
	"reflect"

	"github.com/iancoleman/strcase"
)

func ToSnakeCaseObject(obj any) any {
	if obj == nil {
		return nil
	}

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		return ToSnakeCaseStruct(val)
	case reflect.Map:
		return ToSnakeCaseMap(val)
	case reflect.Slice, reflect.Array:
		return ToSnakeCaseSlice(val)
	default:
		return obj
	}
}

func ToSnakeCaseStruct(val reflect.Value) map[string]any {
	result := make(map[string]any)
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue // Skip unexported fields
		}
		fieldName := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			fieldName = jsonTag
		}
		snakeName := strcase.ToSnake(fieldName)
		result[snakeName] = ToSnakeCaseObject(val.Field(i).Interface())
	}
	return result
}

func ToSnakeCaseMap(val reflect.Value) map[string]any {
	result := make(map[string]any)
	for _, key := range val.MapKeys() {
		if key.Kind() == reflect.String {
			snakeKey := strcase.ToSnake(key.String())
			result[snakeKey] = ToSnakeCaseObject(val.MapIndex(key).Interface())
		}
	}
	return result
}

func ToSnakeCaseSlice(val reflect.Value) []any {
	result := make([]any, val.Len())
	for i := 0; i < val.Len(); i++ {
		result[i] = ToSnakeCaseObject(val.Index(i).Interface())
	}
	return result
}

func ToJSON(obj any) []byte {
	data := ToSnakeCaseObject(obj)
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return jsonData
}
