package logger

import (
	"fmt"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/app"
	"io"
	"os"
	"strings"
	"time"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

type Logger struct {
	level  Level
	writer io.Writer
}

func parseLevel(level string) Level {
	switch strings.ToLower(level) {
	case "debug":
		return Debug
	case "info":
		return Info
	case "warn":
		return Warn
	case "error":
		return Error
	default:
		return Info
	}
}

func New(level string) app.Logger {
	return &Logger{
		level:  parseLevel(level),
		writer: os.Stdout,
	}
}

func (l *Logger) SetWriter(w io.Writer) {
	l.writer = w
}

func (l *Logger) logf(lvl Level, tag, msg string) {
	if lvl < l.level {
		return
	}
	line := fmt.Sprintf("[%s] %s: %s\n", time.Now().Format(time.RFC3339), tag, msg)
	_, _ = l.writer.Write([]byte(line))
}

func (l *Logger) Debug(msg string) {
	l.logf(Debug, "DEBUG", msg)
}

func (l *Logger) Info(msg string) {
	l.logf(Info, "INFO", msg)
}

func (l *Logger) Warn(msg string) {
	l.logf(Warn, "WARN", msg)
}

func (l *Logger) Error(msg string) {
	l.logf(Error, "ERROR", msg)
}
