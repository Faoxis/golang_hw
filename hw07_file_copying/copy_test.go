package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	//nolint:depguard
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type offsetLimitTests struct {
	offset int64
	limit  int64
}

const inputFilePath = "./testdata/input.txt"

func TestCopy(t *testing.T) {
	tests := []offsetLimitTests{
		{
			offset: 0,
			limit:  0,
		},
		{
			offset: 0,
			limit:  10,
		},
		{
			offset: 0,
			limit:  1000,
		},
		{
			offset: 0,
			limit:  10000,
		},
		{
			offset: 100,
			limit:  1000,
		},
		{
			offset: 6000,
			limit:  1000,
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf(
			"Copying file with offset %v and limit %v",
			test.offset,
			test.limit,
		)

		t.Run(name, func(t *testing.T) {
			resultFilePath := fmt.Sprintf("/tmp/result%v.txt", uuid.New().String())
			err := Copy(inputFilePath, resultFilePath, test.offset, test.limit)
			assert.Nil(t, err)

			expectedFilePath := fmt.Sprintf(
				"./testdata/out_offset%v_limit%v.txt",
				test.offset,
				test.limit,
			)
			expectedData, err := os.ReadFile(expectedFilePath)
			assert.Nil(t, err)

			actualData, err := os.ReadFile(resultFilePath)
			assert.Nil(t, err)

			assert.Equal(t, expectedData, actualData, "Содержимое файлов должно совпадать")
		})
	}
}

func TestOffsetExceedsFileSize(t *testing.T) {
	data, err := os.ReadFile(inputFilePath)
	assert.Nil(t, err)

	invalidOffset := int64(len(data) + 1)
	err = Copy(inputFilePath, "/dev/null", invalidOffset, 0)
	assert.ErrorIs(t, err, ErrOffsetExceedsFileSize)
}

func TestInvalidFileName(t *testing.T) {
	err := Copy("/dev/urandom", "/dev/null", 0, 0)
	assert.ErrorIs(t, err, ErrUnsupportedFile)
}

func TestTheSamePathsToFile(t *testing.T) {
	absPath, err := filepath.Abs(inputFilePath)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	err = Copy(inputFilePath, absPath, 0, 0)
	assert.ErrorIs(t, err, ErrFromAndToPathsAreTheSame)
}
