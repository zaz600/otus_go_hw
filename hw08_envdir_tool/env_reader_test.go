package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	type testCase struct {
		name          string
		envDir        string
		expected      Environment
		expectedError error
	}

	tests := []testCase{
		{
			name:   "positive case",
			envDir: "./testdata/env",
			expected: Environment{
				"BAR":   "bar",
				"FOO":   "   foo\nwith new line",
				"HELLO": `"hello"`,
				"UNSET": "",
			},
		},
		{
			name:          "no env dir",
			envDir:        "./testdata/no_env_dir",
			expectedError: os.ErrNotExist,
		},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			env, err := ReadDir(tst.envDir)

			if tst.expectedError == nil {
				require.NoError(t, err)
				require.Equal(t, tst.expected, env)
			} else {
				require.Error(t, err)
				require.True(t, errors.Is(err, tst.expectedError))
			}
		})
	}
}
