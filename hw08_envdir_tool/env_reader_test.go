package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadFilename(t *testing.T) {
	_, err := ReadDir("testdata/badenv")
	assert.NotNil(t, err)

	var badFilenameError *BadFilenameError
	assert.True(t, errors.As(err, &badFilenameError), "File contains '=' cannot be parsed")
}

func TestReadDir(t *testing.T) {
	// Place your code here
	testData := []struct {
		title           string
		filename        string
		expectedContent string
		expectedDelete  bool
	}{
		{"Test usual positive case", "BAR", "bar", false},
		{"Test file with space", "EMPTY", "", false},
		{"Test bad file", "FOO", "   foo\nwith new line", false},
		{"Test file with \"", "HELLO", "\"hello\"", false},
		{"Test empty file", "UNSET", "", true},
	}
	envs, err := ReadDir("testdata/env")
	assert.Nil(t, err)

	for _, data := range testData {
		t.Run(data.title, func(t *testing.T) {
			env := envs[data.filename]
			assert.Equal(t, data.expectedContent, env.Value)
			assert.Equal(t, data.expectedDelete, env.NeedRemove)
		})
	}
}
