package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCmd(t *testing.T) {
	// Place your code here
	testEnvironment := Environment{
		"FOO":       {"foo", false},
		"BAR":       {"bar", false},
		"TO_DELETE": {"", true},
	}

	// Захват stdout
	var stdoutBuf, stderrBuf bytes.Buffer
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	// Создание pipe для подмены stdout и stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	// Выполнение команды

	// testEnv()
	RunCmd([]string{"printenv"}, testEnvironment)

	// Закрытие записи и чтение буфера
	wOut.Close()
	wErr.Close()
	stdoutBuf.ReadFrom(rOut)
	stderrBuf.ReadFrom(rErr)

	// Восстановление оригинальных stdout и stderr
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Проверка результатов
	output := stdoutBuf.String()
	assert.Contains(t, output, "foo")
	assert.Contains(t, output, "bar")
	assert.NotContains(t, output, "to_delete")

	if stderrBuf.Len() != 0 {
		t.Errorf("Expected no stderr, got '%s'", stderrBuf.String())
	}
}
