package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"debug", Debug},
		{"DEBUG", Debug},
		{"info", Info},
		{"INFO", Info},
		{"warn", Warn},
		{"WARN", Warn},
		{"error", Error},
		{"ERROR", Error},
		{"invalid", Info}, // default
	}

	for _, tt := range tests {
		result := parseLevel(tt.input)
		if result != tt.expected {
			t.Errorf("parseLevel(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	outputWriter = &buf
	defer func() { outputWriter = nil }() // reset
	f()
	return buf.String()
}

// Внедряем writer в логгер для тестов
var outputWriter *bytes.Buffer

func (l *Logger) logfTestable(lvl Level, tag, msg string) {
	if lvl < l.level {
		return
	}
	if outputWriter != nil {
		outputWriter.WriteString(tag + ": " + msg + "\n")
	}
}

func TestLoggerFiltering(t *testing.T) {
	l := &Logger{level: Warn}
	var buf bytes.Buffer
	l.SetWriter(&buf)

	l.Debug("debug")
	l.Info("info")
	l.Warn("warn")
	l.Error("error")

	output := buf.String()

	if strings.Contains(output, "debug") || strings.Contains(output, "info") {
		t.Error("unexpected log entries at level < Warn")
	}
	if !strings.Contains(output, "WARN") || !strings.Contains(output, "ERROR") {
		t.Error("expected WARN and ERROR entries missing")
	}
}

func TestLoggerMethods(t *testing.T) {
	var buf bytes.Buffer
	l := &Logger{level: Debug}
	l.SetWriter(&buf)

	l.Debug("debug line")
	l.Info("info line")
	l.Warn("warn line")
	l.Error("error line")

	out := buf.String()

	for _, expected := range []string{"DEBUG", "INFO", "WARN", "ERROR"} {
		if !strings.Contains(out, expected) {
			t.Errorf("expected %s in log", expected)
		}
	}
}
