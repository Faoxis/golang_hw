package hw09structvalidator

import (
	"errors"
	"reflect"
	"strings"
)

type FieldValidator interface {
	Validate(fieldName string, fieldValue reflect.Value, validates []string) error
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

		validator, exists := validators[fieldValue.Kind()]
		if !exists {
			return errors.New("unsupported type: " + fieldValue.Kind().String())
		}
		foundError := validator.Validate(fieldType.Name, fieldValue, validates)
		var foundValidationErrors ValidationErrors
		ok := errors.As(foundError, &foundValidationErrors)
		if !ok {
			return foundError
		}
		validationErrors = append(validationErrors, foundValidationErrors...)
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}
