package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type StringValidator struct{}

func (v StringValidator) Validate(fieldName string, fieldValue reflect.Value, validates []string) error {
	return validateStringField(validates, fieldName, fieldValue.String())
}

func validateStringField(validates []string, fieldType, fieldValue string) error {
	validationErrors := ValidationErrors{}
	for _, validate := range validates {
		if validate == "" {
			break
		}
		var err error
		switch {
		case strings.HasPrefix(validate, "in"):
			inBlock := strings.ReplaceAll(validate, "in:", "")
			ins := strings.Split(inBlock, ",")
			err = runInValidator(fieldType, fieldValue, ins)
		case strings.HasPrefix(validate, "regexp"):
			re := strings.ReplaceAll(validate, "regexp:", "")
			err = runRegexpValidator(fieldType, fieldValue, re)
		case strings.HasPrefix(validate, "len"):
			length := strings.ReplaceAll(validate, "len:", "")
			err = runLenValidator(fieldType, fieldValue, length)
		default:
			err = errors.New("invalid validation parameter")
		}
		if err != nil {
			var validationError ValidationError
			if errors.As(err, &validationError) {
				validationErrors = append(validationErrors, validationError)
			} else {
				return err
			}
		}
	}
	return validationErrors
}

func runLenValidator(name string, s string, length string) error {
	l, err := strconv.Atoi(length)
	if err != nil {
		return errors.New("invalid length format")
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
		return errors.New("invalid regexp")
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
