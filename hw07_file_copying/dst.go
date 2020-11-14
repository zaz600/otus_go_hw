package main

import (
	"fmt"
	"io"
	"os"
)

type dst struct {
	file *os.File
}

type DstOpenError struct {
	path string
	Err  error
}

func (e DstOpenError) Error() string {
	return fmt.Sprintf("Error open %s: %s", e.path, e.Err.Error())
}

func (e DstOpenError) Unwrap() error {
	return e.Err
}

func createDst(path string) (dst, error) {
	fh, err := os.Create(path)
	if err != nil {
		return dst{}, &DstOpenError{path: path, Err: err}
	}

	return dst{file: fh}, nil
}

func openDst(path string) (dst, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return dst{}, &DstOpenError{path: path, Err: err}
		}

		return createDst(path)
	}

	if !fi.Mode().IsRegular() {
		return dst{}, ErrUnsupportedFile
	}

	fh, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return dst{}, &DstOpenError{path: path, Err: err}
	}

	err = fh.Truncate(0)
	if err != nil {
		defer fh.Close()
		return dst{}, &DstOpenError{path: path, Err: err}
	}

	_, err = fh.Seek(0, io.SeekCurrent)
	if err != nil {
		defer fh.Close()
		return dst{}, &DstOpenError{path: path, Err: err}
	}

	return dst{file: fh}, nil
}
