package dev02

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	errInvalidEsc          = errors.New("Invalid escape position")
	errInvalidStringFormat = errors.New("Invalid string format")
)

// UnpackString performs primitive unpacking of a string containing repeating characters / runes
// with support for escape sequences.
// The function takes a string as input and returns the unpacked string and an error, if applicable.
func UnpackString(s string) (string, error) {
	var result strings.Builder    // Using strings.Builder for efficient string concatenation
	var escape bool               // Flag to track escape sequences
	var nextRuneIsCountDigit bool // keeping track of the type of the next rune to skip on purpose

	if s == "" { // If the string is empty, return an empty string
		return "", nil
	}
	// Iterate over each character in the input string
	for i, r := range s {
		if r == '\\' && !escape {
			// check that the escape symbol is not first nor last rune and also followed by a digit or another escape symbol
			if i == 0 || i >= len(s)-1 || (!unicode.IsDigit(rune(s[i+1])) && s[i+1] != '\\') {
				return "", errInvalidEsc
			} else {
				escape = true
			}
		} else if nextRuneIsCountDigit { // skip if the rune is a repeatCounter
			nextRuneIsCountDigit = false
			continue
		} else if i < len(s)-1 && unicode.IsDigit(rune(s[i+1])) {
			repeatCount, err := strconv.Atoi(string(s[i+1]))
			if err != nil {
				return "", errInvalidStringFormat
			}
			result.WriteString(strings.Repeat(string(r), repeatCount)) // write the rune repeatCount times
			escape = false                                             // resetting escape
			nextRuneIsCountDigit = true
		} else {
			result.WriteRune(r) // write a single rune to result
			escape = false      // resetting escape
			nextRuneIsCountDigit = false
		}

	}

	// Return the final unpacked string
	return result.String(), nil
}
