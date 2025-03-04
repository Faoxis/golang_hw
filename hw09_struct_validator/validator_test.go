package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	Test struct {
		in          interface{}
		expectedErr error
	}

	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Number struct {
		Value string `validate:"in:one,two,three,four,five,six,seven,eight,nine,ten"`
	}

	NoValidate struct {
		Value  string
		Number int
	}

	WrongRegexp struct {
		Value string `validate:"regexp:[a-z"`
	}

	OnlyNumberRegexp struct {
		Value string `validate:"regexp:^\\d+$"`
	}

	OnlyLetterRegexp struct {
		Value string `validate:"regexp:^[a-zA-Z]+$"`
	}

	MinMaxInt struct {
		Value int `validate:"min:1|max:10"`
	}

	InInt struct {
		Value int `validate:"in:1,2,3"`
	}
	Slice struct {
		IntValues    []int    `validate:"min:1|max:5"`
		StringValues []string `validate:"regexp:^\\d+$|in:1,2,3"`
	}

	InvalidValidationType struct {
		RunValue rune `validate:"regexp:^\\d+$|in:1,2,3"`
	}
)

func TestLenInString(t *testing.T) {
	tests := []Test{
		// len = 5
		{
			App{
				Version: "1.0.0",
			},
			nil,
		},
		// len = 3
		{
			App{
				Version: "1.0",
			},
			ValidationErrors{
				{Field: "Version", Err: errors.New("invalid length")},
			},
		},
		// len = 6 - error
		{
			App{
				Version: "1.0.0.0",
			},
			ValidationErrors{
				{Field: "Version", Err: errors.New("invalid length")},
			},
		},
	}

	test(tests, t)
}

func TestRegexpString(t *testing.T) {
	tests := []Test{
		{
			WrongRegexp{
				Value: "some value",
			},
			errors.New("invalid regexp"),
		},
		{
			OnlyNumberRegexp{
				Value: "1",
			},
			nil,
		},
		{
			OnlyNumberRegexp{
				Value: "one",
			},
			ValidationErrors{
				{Field: "Value", Err: errors.New("does not match regexp ^\\d+$")},
			},
		},
		{
			OnlyLetterRegexp{
				Value: "one",
			},
			nil,
		},
		{
			OnlyLetterRegexp{
				Value: "1",
			},
			ValidationErrors{
				{Field: "Value", Err: errors.New("does not match regexp ^[a-zA-Z]+$")},
			},
		},
	}

	test(tests, t)
}

func TestInString(t *testing.T) {
	tests := []Test{
		{
			Number{
				Value: "one",
			},
			nil,
		},
		{
			Number{
				Value: "forty-two",
			},
			ValidationErrors{
				{Field: "Value", Err: errors.New("forty-two is not in possible values")},
			},
		},
	}

	test(tests, t)
}

func TestMinMaxInt(t *testing.T) {
	tests := []Test{
		{
			MinMaxInt{
				Value: 1,
			},
			nil,
		},
		{
			MinMaxInt{
				Value: 10,
			},
			nil,
		},
		{
			MinMaxInt{
				Value: 5,
			},
			nil,
		},
		{
			MinMaxInt{
				Value: 0,
			},
			ValidationErrors{
				{Field: "Value", Err: errors.New("min value validation error")},
			},
		},
		{
			MinMaxInt{
				Value: 11,
			},
			ValidationErrors{
				{Field: "Value", Err: errors.New("max value validation error")},
			},
		},
	}

	test(tests, t)
}

func TestInInt(t *testing.T) {
	tests := []Test{
		{
			InInt{
				Value: 1,
			},
			nil,
		},
		{
			InInt{
				Value: 3,
			},
			nil,
		},
		{
			InInt{
				Value: 42,
			},
			ValidationErrors{
				{Field: "Value", Err: errors.New("value in range validation error")},
			},
		},
	}

	test(tests, t)
}

func TestSliceValidate(t *testing.T) {
	tests := []Test{
		{
			Slice{
				IntValues:    []int{1, 2, 3},
				StringValues: []string{"1", "2", "3"},
			},
			nil,
		},
		{
			Slice{
				IntValues:    []int{1, 7, 6, 5},
				StringValues: []string{"1", "2", "3"},
			},
			ValidationErrors{
				{Field: "IntValues", Err: errors.New("max value validation error")},
				{Field: "IntValues", Err: errors.New("max value validation error")},
			},
		},
		{
			Slice{
				IntValues:    []int{1, 2, 3},
				StringValues: []string{"1", "2", "six", "6", "3"},
			},
			ValidationErrors{
				{Field: "StringValues", Err: errors.New("does not match regexp ^\\d+$")},
				{Field: "StringValues", Err: errors.New("six is not in possible values")},
				{Field: "StringValues", Err: errors.New("6 is not in possible values")},
			},
		},
	}

	test(tests, t)
}

func TestCommonValidate(t *testing.T) {
	tests := []Test{
		{
			NoValidate{
				Value:  "some value",
				Number: 42,
			},
			nil,
		},
	}

	test(tests, t)
}

func TestInvalidUserData(t *testing.T) {
	tests := []Test{
		{
			InvalidValidationType{
				RunValue: '1',
			},
			errors.New("unsupported type: int32"),
		},
	}
	test(tests, t)
}

//nolint:thelper
func test(tests []Test, t *testing.T) {
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			// Place your code here.
			err := Validate(tt.in)
			if tt.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedErr.Error())
			}
			_ = tt
		})
	}
}
