package hw02_unpack_string //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	input    string
	expected string
	err      error
}

func TestUnpack(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			input:    "abccd",
			expected: "abccd",
		},
		{
			input:    "3abc",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "45",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "aaa10b",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "",
			expected: "",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpacWithPunctuationMarks(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "!!!",
			expected: "!!!",
		},
		{
			input:    "!!?4",
			expected: "!!????",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpacWithPunctuationMarksAndLetters(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "a!a!a!a",
			expected: "a!a!a!a",
		},
		{
			input:    "!ab!?4",
			expected: "!ab!????",
		},
		{
			input:    "ab!?4",
			expected: "ab!????",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpacWithEmoji(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "🥰",
			expected: "🥰",
		},
		{
			input:    "🥰4",
			expected: "🥰🥰🥰🥰",
		},
		{
			input:    "ab🥰4",
			expected: "ab🥰🥰🥰🥰",
		},
		{
			input:    "ab5🥰4",
			expected: "abbbbb🥰🥰🥰🥰",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}
