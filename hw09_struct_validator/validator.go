package hw09structvalidator

import (
	"errors"
	"reflect"
	"strings"
)

type FieldValidator interface {
	Validate(fieldName string, fieldValue reflect.Value, validates []string) ValidationErrors
}

var validators = map[reflect.Kind]FieldValidator{
	reflect.String: StringValidator{},
	reflect.Int:    IntValidator{},
	reflect.Slice:  SliceValidator{},
}

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	typeValue := reflect.TypeOf(v)

	if value.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}

	validationErrors := ValidationErrors{}

	for i := 0; i < typeValue.NumField(); i++ {
		fieldValue := value.Field(i)
		fieldType := typeValue.Field(i)

		validateTag := fieldType.Tag.Get("validate")
		validates := strings.Split(validateTag, "|")

		if validator, exists := validators[fieldValue.Kind()]; exists {
			validationErrors = append(validationErrors, validator.Validate(fieldType.Name, fieldValue, validates)...)
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}
