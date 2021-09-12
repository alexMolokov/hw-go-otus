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

	in := []rune(str)
	repeat := make([]int, len(in))

	for i := 0; i < len(in); i++ {
		if !unicode.IsDigit(in[i]) {
			repeat[i] = 1
			continue
		}
		if i == 0 || unicode.IsDigit(in[i-1]) {
			return "", ErrInvalidString
		}
		count, _ := strconv.Atoi(string(in[i]))
		repeat[i-1] = count
		repeat[i] = 0
	}

	result := strings.Builder{}
	for i := 0; i < len(repeat); i++ {
		if repeat[i] == 0 {
			continue
		}
		result.WriteString(strings.Repeat(string(in[i]), repeat[i]))
	}

	return result.String(), nil
}
