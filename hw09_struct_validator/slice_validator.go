package hw09structvalidator

import "reflect"

type SliceValidator struct{}

func (sv SliceValidator) Validate(fieldName string, fieldValue reflect.Value, validates []string) ValidationErrors {
	validationErrors := ValidationErrors{}
	for i := 0; i < fieldValue.Len(); i++ {
		elem := fieldValue.Index(i)
		if elem.Kind() == reflect.String {
			validationErrors = append(validationErrors, validateStringField(validates, fieldName, elem.String())...)
		}
		if elem.Kind() == reflect.Int {
			validationErrors = append(validationErrors, validateIntField(validates, fieldName, elem.Int())...)
		}
	}
	return validationErrors
}
