package main

import (
	"fmt"
	"io"
	"os"
)

type src struct {
	file     *os.File
	copySize int64
}
type SrcOpenError struct {
	path string
	Err  error
}

func (e SrcOpenError) Error() string {
	return fmt.Sprintf("Error open %s: %s", e.path, e.Err.Error())
}

func (e SrcOpenError) Unwrap() error {
	return e.Err
}

func openSrc(path string, offset, limit int64) (src, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return src{}, ErrSrcFileNotFound
		}

		return src{}, &SrcOpenError{path: path, Err: err}
	}

	if !fi.Mode().IsRegular() {
		return src{}, ErrUnsupportedFile
	}

	copySize, err := getCopySize(fi, offset, limit)
	if err != nil {
		return src{}, &SrcOpenError{path: path, Err: err}
	}

	fh, err := os.Open(path)
	if err != nil {
		return src{}, &SrcOpenError{path: path, Err: err}
	}

	_, err = fh.Seek(offset, io.SeekStart)
	if err != nil {
		defer fh.Close()
		return src{}, &SrcOpenError{path: path, Err: err}
	}

	return src{
		file:     fh,
		copySize: copySize,
	}, nil
}

func getCopySize(fi os.FileInfo, offset, limit int64) (int64, error) {
	if offset >= fi.Size() {
		return 0, ErrOffsetExceedsFileSize
	}

	copyLen := fi.Size() - offset
	if limit != 0 && limit < copyLen {
		copyLen = limit
	}
	return copyLen, nil
}
