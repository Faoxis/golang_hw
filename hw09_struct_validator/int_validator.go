package hw09structvalidator

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type IntValidator struct{}

func (v IntValidator) Validate(fieldName string, fieldValue reflect.Value, validates []string) ValidationErrors {
	return validateIntField(validates, fieldName, fieldValue.Int())
}

func validateIntField(validates []string, name string, value int64) ValidationErrors {
	validationErrors := ValidationErrors{}
	for _, validate := range validates {
		if validate == "" {
			break
		}
		if strings.HasPrefix(validate, "max") {
			maxVal := strings.ReplaceAll(validate, "max:", "")
			err := runMaxValidator(name, value, maxVal)
			if err != nil {
				var validationError ValidationError
				if errors.As(err, &validationError) {
					validationErrors = append(validationErrors, validationError)
				}
			}
		}
		if strings.HasPrefix(validate, "min") {
			maxVal := strings.ReplaceAll(validate, "min:", "")
			err := runMinValidator(name, value, maxVal)
			if err != nil {
				var validationError ValidationError
				if errors.As(err, &validationError) {
					validationErrors = append(validationErrors, validationError)
				}
			}
		}
		if strings.HasPrefix(validate, "in") {
			in := strings.ReplaceAll(validate, "in:", "")
			ins := strings.Split(in, ",")
			err := runIntInValidator(name, value, ins)
			if err != nil {
				var validationError ValidationError
				if errors.As(err, &validationError) {
					validationErrors = append(validationErrors, validationError)
				}
			}
		}
	}
	return validationErrors
}

func runIntInValidator(name string, value int64, ins []string) error {
	insInt := make([]int64, 0, len(ins))
	for _, in := range ins {
		v, err := strconv.ParseInt(in, 10, 64)
		if err != nil {
			return NewValidationError(
				name,
				errors.New("invalid in value format"),
			)
		}
		insInt = append(insInt, v)
	}
	for _, in := range insInt {
		if value == in {
			return nil
		}
	}
	return NewValidationError(
		name,
		errors.New("value in range validation error"),
	)
}

func runMinValidator(name string, value int64, minVal string) error {
	minValue, err := strconv.ParseInt(minVal, 10, 64)
	if err != nil {
		return NewValidationError(
			name,
			errors.New("invalid min value format"),
		)
	}
	if value < minValue {
		return NewValidationError(
			name,
			errors.New("min value validation error"),
		)
	}
	return nil
}

func runMaxValidator(name string, value int64, maxVal string) error {
	maxValue, err := strconv.ParseInt(maxVal, 10, 64)
	if err != nil {
		return NewValidationError(
			name,
			errors.New("invalid max value format"),
		)
	}
	if value > maxValue {
		return NewValidationError(
			name,
			errors.New("max value validation error"),
		)
	}
	return nil
}
