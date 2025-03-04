package hw09structvalidator

import (
	"errors"
	"reflect"
)

type SliceValidator struct{}

func (sv SliceValidator) Validate(fieldName string, fieldValue reflect.Value, validates []string) error {
	validationErrors := ValidationErrors{}
	for i := 0; i < fieldValue.Len(); i++ {
		elem := fieldValue.Index(i)
		var foundError error
		//nolint:exhaustive
		switch elem.Kind() {
		case reflect.String:
			foundError = validateStringField(validates, fieldName, elem.String())
		case reflect.Int:
			foundError = validateIntField(validates, fieldName, elem.Int())
		default:
			foundError = errors.New("unsupported type: " + elem.Kind().String())
		}
		var foundValidationErrors ValidationErrors
		ok := errors.As(foundError, &foundValidationErrors)
		if !ok {
			return foundError
		}
		validationErrors = append(validationErrors, foundValidationErrors...)
	}
	return validationErrors
}
