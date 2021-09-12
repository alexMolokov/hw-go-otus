package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}

	result := strings.Builder{}
	in := []rune(str)

	if unicode.IsDigit(in[0]) {
		return "", ErrInvalidString
	}

	isDigit := false
	for i := 0; i < len(in); i++ {
		if unicode.IsDigit(in[i]) {
			if isDigit {
				return "", ErrInvalidString
			}
			repeat, _ := strconv.Atoi(string(in[i]))
			result.WriteString(strings.Repeat(string(in[i-1]), repeat))
			isDigit = true
			continue
		}

		if !isDigit && i > 0 {
			result.WriteRune(in[i-1])
		}

		isDigit = false
	}

	if !isDigit {
		result.WriteRune(in[len(in)-1])
	}

	return result.String(), nil
}
