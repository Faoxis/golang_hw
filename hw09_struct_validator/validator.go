package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

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

func Validate(v interface{}) error {
	// Place your code here.
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
		switch fieldValue.Kind() {
		case reflect.String:
			stringValidationErrors := validateStringField(validates, fieldType.Name, fieldValue.String())
			validationErrors = append(validationErrors, stringValidationErrors...)
		case reflect.Int:
			intValidationErrors := validateIntField(validates, fieldType.Name, fieldValue.Int())
			validationErrors = append(validationErrors, intValidationErrors...)
		case reflect.Slice:
			if fieldValue.Type().Elem().Kind() == reflect.String {
				for i := 0; i < fieldValue.Len(); i++ {
					stringValidationErrors := validateStringField(validates, fieldType.Name, fieldValue.Index(i).String())
					validationErrors = append(validationErrors, stringValidationErrors...)
				}
			}
			if fieldValue.Type().Elem().Kind() == reflect.Int {
				for i := 0; i < fieldValue.Len(); i++ {
					stringValidationErrors := validateIntField(validates, fieldType.Name, fieldValue.Index(i).Int())
					validationErrors = append(validationErrors, stringValidationErrors...)
				}
			}
		default:
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
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
	insInt := make([]int64, len(ins))
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

func validateStringField(validates []string, fieldType, fieldValue string) ValidationErrors {
	validationErrors := ValidationErrors{}
	for _, validate := range validates {
		if validate == "" {
			break
		}
		if strings.HasPrefix(validate, "in") {
			inBlock := strings.ReplaceAll(validate, "in:", "")
			ins := strings.Split(inBlock, ",")
			err := runInValidator(fieldType, fieldValue, ins)
			if err != nil {
				var validationError ValidationError
				if errors.As(err, &validationError) {
					validationErrors = append(validationErrors, validationError)
				}
			}
		}
		if strings.HasPrefix(validate, "regexp") {
			regexp := strings.ReplaceAll(validate, "regexp:", "")
			err := runRegexpValidator(fieldType, fieldValue, regexp)
			if err != nil {
				var validationError ValidationError
				if errors.As(err, &validationError) {
					validationErrors = append(validationErrors, validationError)
				}
			}
		}
		if strings.HasPrefix(validate, "len") {
			length := strings.ReplaceAll(validate, "len:", "")
			err := runLenValidator(fieldType, fieldValue, length)
			if err != nil {
				if validationError, ok := err.(ValidationError); ok {
					validationErrors = append(validationErrors, validationError)
				}
			}
		}
	}
	return validationErrors
}

func runLenValidator(name string, s string, length string) interface{} {
	l, err := strconv.Atoi(length)
	if err != nil {
		return NewValidationError(
			name,
			errors.New("invalid length format"),
		)
	}
	if len(s) != l {
		return NewValidationError(
			name,
			errors.New("invalid length"),
		)
	}
	return nil
}

func runRegexpValidator(name, value string, regexps string) error {
	reg, err := regexp.Compile(regexps)
	if err != nil {
		return NewValidationError(
			name,
			errors.New("invalid regexp"),
		)
	}
	if !reg.MatchString(value) {
		errMsg := fmt.Sprintf("does not match regexp %s", regexps)
		return NewValidationError(
			name,
			errors.New(errMsg),
		)
	}

	return nil
}

func runInValidator(fieldName, value string, possibleValues []string) error {
	contains := false
	for _, in := range possibleValues {
		if value == in {
			contains = true
		}
	}
	if !contains {
		errMsg := fmt.Sprintf("%s is not in possible values", value)
		return NewValidationError(
			fieldName,
			errors.New(errMsg),
		)
	}
	return nil
}
