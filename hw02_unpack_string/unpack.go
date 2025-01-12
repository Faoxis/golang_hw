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
	lastWasEscaped := false

	for _, symbol := range str {
		switch {
		// Если текущий символ - число
		case isNumber(symbol) && !lastWasEscaped:
			// Если за числом идет число, то это ошибка
			if lastWasNumber {
				return "", ErrInvalidString
			}
			lastWasNumber = true

			// Вычисляем число и добавляем нужное количество символов
			number := int(symbol - '0')
			if number == 0 {
				newStr = newStr[:len(newStr)-1]
			} else {
				for i := 0; i < number-1; i++ {
					newStr = append(newStr, lastSymbol)
				}
			}

			// Число - не экранирующий символ. Отмечаем
			lastWasEscaped = false

		// Если текущий символ - символ экранирования и встречается в первых раз
		case isEscaped(symbol) && !lastWasEscaped:
			lastWasEscaped = true // Просто отмечаем это

		// Во всех остальных случаях просто переносим символ из исходной строки в результирующую
		default:
			if lastWasEscaped && !isNumber(symbol) && !isEscaped(symbol) {
				return "", ErrInvalidString
			}

			lastSymbol = symbol
			newStr = append(newStr, symbol)

			// И отмечаем, что текущий символ - не число и не символ экранирования
			lastWasNumber = false
			lastWasEscaped = false
		}
	}

	// Если закончили на экранирующем символе
	if lastWasEscaped {
		return "", ErrInvalidString
	}

	return string(newStr), nil
}

func isEscaped(symbol rune) bool {
	return symbol == '\\'
}

func isNumber(symbol rune) bool {
	return symbol <= '9' && symbol >= '0'
}
