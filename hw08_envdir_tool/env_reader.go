package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Environment map[string]string

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	result := make(Environment)
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return result, err
	}

	for _, f := range files {
		if !f.Mode().IsRegular() || strings.Contains(f.Name(), "=") {
			continue
		}

		line, err := readLine(path.Join(dir, f.Name()))
		if err != nil {
			return result, err
		}
		result[f.Name()] = line
	}
	return result, nil
}

func readLine(filePath string) (string, error) {
	fd, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer fd.Close()

	reader := bufio.NewReader(fd)

	lineBytes, _, err := reader.ReadLine()
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	return trimLine(lineBytes), nil
}

func trimLine(lineBytes []byte) string {
	line := string(bytes.ReplaceAll(lineBytes, []byte{0x00}, []byte{'\n'}))
	return strings.TrimRight(line, "\t ")
}
