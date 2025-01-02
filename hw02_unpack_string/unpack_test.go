package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

//		// Base cases
//		{"Basic string", "a4bc2d5e", "aaaabccddddde", nil},
//		{"Single char no repeat", "abc", "abc", nil},
//		{"No input", "", "", nil},
//
//		// Edge cases
//		{"Zero after char", "a0bc2d", "bccd", nil},
//		{"Zero leading to empty string", "a0b0c0", "", nil},
//		{"Number at start", "3abc", "", ErrInvalidString},
//		{"Consecutive numbers", "ab12c", "", ErrInvalidString},
//
//		// Escape character cases
//		{"Escaped digit", "a\\4bc2d5", "a4bccddddd", nil},
//		{"Escaped triple backslash", "a\\\\3b", "a\\\\\\b", nil},
//		{"Unfinished escape", "a\\", "a\\", ErrInvalidString},
//
//		// Complex cases
//		{"Mixed escape and normal", "a\\2b3c\\0d", "a2bbbc0d", nil},
//		{"Only escaped digits", "\\1\\2\\3", "123", nil},
//		{"Trailing escape", "abc\\", "", ErrInvalidString},
//		{"Last symbol is escaped escape", "abc\\\\", "abc\\", nil},

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "🙃0", expected: ""},
		{input: "aaф0b", expected: "aab"},

		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},

		// Escape character cases
		{"a\\4bc2d5", "a4bccddddd"},
		{"a\\\\3b", "a\\\\\\b"},
		{"a\\", ""},

		// Complex cases
		{"a\\2b3c\\0d", "a2bbbc0d"},
		{"\\1\\2\\3", "123"},
		{"abc\\", ""}, // error
		{"abc\\\\", "abc\\"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

//func TestUnpack(t *testing.T) {
//	cases := []struct {
//		name   string
//		input  string
//		output string
//		err    error
//	}{
//		// Base cases
//		{"Basic string", "a4bc2d5e", "aaaabccddddde", nil},
//		{"Single char no repeat", "abc", "abc", nil},
//		{"No input", "", "", nil},
//
//		// Edge cases
//		{"Zero after char", "a0bc2d", "bccd", nil},
//		{"Zero leading to empty string", "a0b0c0", "", nil},
//		{"Number at start", "3abc", "", ErrInvalidString},
//		{"Consecutive numbers", "ab12c", "", ErrInvalidString},
//
//		// Escape character cases
//		{"Escaped digit", "a\\4bc2d5", "a4bccddddd", nil},
//		{"Escaped triple backslash", "a\\\\3b", "a\\\\\\b", nil},
//		{"Unfinished escape", "a\\", "a\\", ErrInvalidString},
//
//		// Complex cases
//		{"Mixed escape and normal", "a\\2b3c\\0d", "a2bbbc0d", nil},
//		{"Only escaped digits", "\\1\\2\\3", "123", nil},
//		{"Trailing escape", "abc\\", "", ErrInvalidString},
//		{"Last symbol is escaped escape", "abc\\\\", "abc\\", nil},
//	}
//
//	for _, testCase := range cases {
//		t.Run(testCase.name, func(t *testing.T) {
//			result, err := Unpack(testCase.input)
//			if testCase.err == nil {
//				assert.Equal(
//					t,
//					testCase.output,
//					result,
//					"For input (%q) expected (%q), got (%q, %v)",
//					testCase.input, testCase.output, result, err,
//				)
//			} else {
//				assert.ErrorIs(
//					t,
//					testCase.err,
//					err,
//					"For input (%q) expected error (%v), but got string (%q, %v)",
//					testCase.input, testCase.output, result, err,
//				)
//			}
//		})
//	}
//}
