package hw02unpackstring

import (
	"errors"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}

	newStr := make([]rune, 0)

	var lastSymbol rune
	lastWasNumber := true
	escaped := false

	for _, v := range str {
		switch {
		// Если текущий символ - число
		case v <= '9' && v >= '0' && !escaped:
			// Если за числом идет число, то это ошибка
			if lastWasNumber {
				return "", ErrInvalidString
			}
			lastWasNumber = true

			// Вычисляем число и добавляем нужное количество символов
			number := int(v - '0')
			if number == 0 {
				newStr = newStr[:len(newStr)-1]
			} else {
				for i := 0; i < number-1; i++ {
					newStr = append(newStr, lastSymbol)
				}
			}

			// Число - не экранирующий символ. Отмечаем
			escaped = false

		// Если текущий символ - символ экранирования и встречается в первых раз
		case v == '\\' && !escaped:
			escaped = true // Просто отмечаем это

		// Во всех остальных случаях просто переносим символ из исходной строки в результирующую
		default:
			if escaped {
				return "", ErrInvalidString
			}

			lastSymbol = v
			newStr = append(newStr, v)

			// И отмечаем, что текущий символ - не число и не символ экранирования
			lastWasNumber = false
			escaped = false
		}
	}

	// Если закончили на экранирующем символе
	if escaped {
		return "", ErrInvalidString
	}

	return string(newStr), nil
}
