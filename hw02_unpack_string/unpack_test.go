package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "correct", input: "a4bc", expected: "aaaabc"},
		{name: "correct", input: "a4bc2d5e", expected: "aaaabccddddde"},
		{name: "no number", input: "abccd", expected: "abccd"},
		{name: "empty string", input: "", expected: ""},
		{name: "digit is 0 ", input: "aaa0b", expected: "aab"},
		{name: "one symbol", input: "a", expected: "a"},
		{name: "repeat one time", input: "a1", expected: "a"},
		{name: "capital", input: "A2s3K2", expected: "AAsssKK"},
		{name: "cyrillic", input: "ы2я4", expected: "ыыяяяя"},
		{name: "minus", input: "a-5", expected: "a-----"},
		{name: "whitespace in", input: "a2 2", expected: "aa  "},
		{name: "empty null digit", input: "a0", expected: ""},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "5"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
