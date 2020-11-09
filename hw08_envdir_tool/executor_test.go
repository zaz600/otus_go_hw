package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := make(Environment)
	returnCode := RunCmd([]string{"./testdata/echo.sh"}, env)
	require.Equal(t, 0, returnCode)
}

func TestRunCmdSetEnv(t *testing.T) {
	env := Environment{"FOO": "bar"}
	returnCode := RunCmd([]string{"echo", "abcd"}, env)
	require.Equal(t, 0, returnCode)

	val, ok := os.LookupEnv("FOO")
	require.True(t, ok)
	require.Equal(t, "bar", val)
}

func TestRunCmdUnSetEnv(t *testing.T) {
	os.Setenv("FOO", "bar")

	env := Environment{"FOO": ""}
	returnCode := RunCmd([]string{"echo", "abcd"}, env)
	require.Equal(t, 0, returnCode)

	_, ok := os.LookupEnv("FOO")
	require.False(t, ok)
}

func TestRunCmdNotExistsArgs(t *testing.T) {
	env := make(Environment)
	returnCode := RunCmd([]string{"/bin/bash", "abcdefg"}, env)
	require.Equal(t, 127, returnCode)
}
