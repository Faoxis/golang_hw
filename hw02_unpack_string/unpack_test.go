package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

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
		{input: "d2ф3", expected: "ddффф"},
		{input: "🎉2🎊3", expected: "🎉🎉🎊🎊🎊"},
		{input: "a2b0c3", expected: "aaccc"},
		{input: `\52`, expected: "55"},
		{input: "ф2ы3", expected: "ффыыы"},
		{input: "aaф0b", expected: "aab"},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},

		{input: `\\\\\\\\`, expected: `\\\\`},     // Четное количество экранирующих символов
		{input: `\\0`, expected: ""},              // Экранированный ноль удаляет предыдущий символ
		{input: `a5b3🙃2`, expected: `aaaaabbb🙃🙃`}, // Разные типы символов в одной строке
		{input: `🌟3a0`, expected: `🌟🌟🌟`},          // Эмодзи с повторением и удалением
		{input: `\1\2\3`, expected: `123`},        // Экранированные цифры
		{input: `a\2b\3c`, expected: `a2b3c`},     // Чередование экранирования и букв
		{input: `a\\`, expected: `a\`},            // Экранирование в конце строки
		{input: `a1b1c1`, expected: `abc`},        // Повторение на 1 не меняет строку
		{input: `d1e1f1`, expected: `def`},
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
	invalidStrings := []string{
		"3abc",
		"45",
		"aaa10b",
		`qwe\s45`,
		`\\\`, // Строка заканчивается экранирующим симоволом
		`a\`,  // Завершение строки на экранирующем символе

	}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
