package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func IsValidString(input string) bool {
	// число
	if _, err := strconv.Atoi(input); err == nil {
		return false
	}

	// начинается с цифры
	if _, err := strconv.Atoi(input[:1]); err == nil {
		return false
	}

	// есть подстрока-числа
	if regexp.MustCompile(`\d{2,}`).FindStringIndex(input) != nil {
		return false
	}

	return true
}

func Unpack(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}

	if !IsValidString(input) {
		return "", ErrInvalidString
	}

	var result strings.Builder
	var prevChar string

	for _, currRune := range input {
		currChar := string(currRune)

		if unicode.IsDigit(currRune) {
			repeatCount, _ := strconv.Atoi(currChar)
			repeatCount--
			result.WriteString(strings.Repeat(prevChar, repeatCount))
		} else {
			result.WriteString(currChar)
		}

		prevChar = currChar
	}

	return result.String(), nil
}
