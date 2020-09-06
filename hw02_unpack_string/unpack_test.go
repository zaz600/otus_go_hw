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
			input:    "a4b5",
			expected: "aaaabbbbb",
		},
		{
			input:    "5",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "世界4bc2d5e界",
			expected: "世界界界界bccddddde界",
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
			input:    "1",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "aaa0b",
			expected: "aab",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpackWithEscape(t *testing.T) {
	// t.Skip() // Remove if task with asterisk completed

	for _, tst := range [...]test{
		{
			input:    `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			input:    `qwe\45`,
			expected: `qwe44444`,
		},
		{
			input:    `qwe2\45`,
			expected: `qwee44444`,
		},
		{
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			input:    `qwe\\\3`,
			expected: `qwe\3`,
		},
		{
			input:    `qwe3\`,
			expected: ``,
			err:      ErrInvalidString,
		},
		{
			input:    `qwe3\\`,
			expected: `qweee\`,
		},
		{
			input:    `\\`,
			expected: `\`,
		},
		{
			input:    `\\\`,
			expected: ``,
			err:      ErrInvalidString,
		},
		{
			input:    `\`,
			expected: ``,
			err:      ErrInvalidString,
		},
		{
			input:    `\x`,
			expected: ``,
			err:      ErrInvalidString,
		},
		{
			input:    `a\x`,
			expected: ``,
			err:      ErrInvalidString,
		},
		{
			input:    `a\\\x`,
			expected: ``,
			err:      ErrInvalidString,
		},
		{
			input:    `a\\x`,
			expected: `a\x`,
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func BenchmarkUnpack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Unpack(`qwe\\\3`)
	}
}
