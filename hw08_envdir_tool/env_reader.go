package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type BadFilenameError struct {
	Message string
}

func (e BadFilenameError) Error() string {
	return e.Message
}

func NewBadFilenameError(message string) *BadFilenameError {
	return &BadFilenameError{
		Message: message,
	}
}

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading directory: %w", err)
	}

	env := make(Environment)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		envValue, err := parseFile(dir, entry)
		if err != nil {
			return Environment{}, err
		}
		env[entry.Name()] = envValue
	}

	return env, nil
}

func parseFile(directory string, entry os.DirEntry) (EnvValue, error) {
	filename := entry.Name()
	if strings.Contains(filename, "=") {
		errorMessage := fmt.Sprintf("file %s contains '='", filename)
		return EnvValue{}, NewBadFilenameError(errorMessage)
	}
	filePath := filepath.Join(directory, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return EnvValue{}, fmt.Errorf("opening file %s: %w", filename, err)
	}
	defer func() {
		e := file.Close()
		if e != nil {
			fmt.Printf("ERROR: %v\n", e)
		}
	}()

	scanner := bufio.NewScanner(file)
	var line string
	if scanner.Scan() {
		line = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return EnvValue{}, fmt.Errorf("scanning file %s: %w", filename, err)
	}

	remove := false
	if len(line) == 0 {
		remove = true
	}

	// Заменить нулевые байты на символы новой строки
	line = strings.ReplaceAll(line, "\x00", "\n")

	// Удалить завершающие пробелы и табуляции
	line = strings.TrimRight(line, " \t")

	return EnvValue{
		Value:      line,
		NeedRemove: remove,
	}, nil
}
