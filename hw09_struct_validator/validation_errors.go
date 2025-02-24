package hw09structvalidator

import "strings"

type ValidationError struct {
	Field string
	Err   error
}

func NewValidationError(field string, err error) ValidationError {
	return ValidationError{
		Field: field,
		Err:   err,
	}
}

func (v ValidationError) Error() string {
	return v.Field + ": " + v.Err.Error() + "\n"
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	stringBuilder := strings.Builder{}
	for _, err := range v {
		stringBuilder.WriteString(err.Error() + "\n")
	}
	return stringBuilder.String()
}
