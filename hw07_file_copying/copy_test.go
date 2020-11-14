package main

import (
	"crypto/md5"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	inputFile := "testdata/input.txt"
	unknownFile := "testdata/not_exists_file.txt"

	type testCase struct {
		name          string
		inputFile     string
		goldenFile    string
		offset        int64
		limit         int64
		expectedError error
	}

	tests := []testCase{
		{
			name:       "offset 0 limit 0",
			inputFile:  inputFile,
			goldenFile: "testdata/out_offset0_limit0.txt",
			offset:     0,
			limit:      0,
		},
		{
			name:       "offset 0 limit 10",
			inputFile:  inputFile,
			goldenFile: "testdata/out_offset0_limit10.txt",
			offset:     0,
			limit:      10,
		},
		{
			name:       "offset 0 limit 1000",
			inputFile:  inputFile,
			goldenFile: "testdata/out_offset0_limit1000.txt",
			offset:     0,
			limit:      1000,
		},
		{
			name:       "offset 0 limit 10000",
			inputFile:  inputFile,
			goldenFile: "testdata/out_offset0_limit10000.txt",
			offset:     0,
			limit:      10000,
		},
		{
			name:       "offset 100 limit 1000",
			inputFile:  inputFile,
			goldenFile: "testdata/out_offset100_limit1000.txt",
			offset:     100,
			limit:      1000,
		},
		{
			name:       "offset 6000 limit 1000",
			inputFile:  inputFile,
			goldenFile: "testdata/out_offset6000_limit1000.txt",
			offset:     6000,
			limit:      1000,
		},
		{
			name:          "offset > file size",
			inputFile:     inputFile,
			offset:        300000,
			expectedError: ErrOffsetExceedsFileSize,
		},
		{
			name:          "not found source file",
			inputFile:     unknownFile,
			expectedError: ErrSrcFileNotFound,
		},
		{
			name:          "unsupported file",
			inputFile:     "/dev/random",
			expectedError: ErrUnsupportedFile,
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			dst, err := ioutil.TempFile("/tmp", "zaz600")
			require.NoError(t, err)
			defer os.Remove(dst.Name())

			err = Copy(tst.inputFile, dst.Name(), tst.offset, tst.limit)
			if tst.expectedError != nil {
				require.Error(t, err)
				require.True(t, errors.Is(err, tst.expectedError))
			} else {
				require.NoError(t, err)
				defer dst.Close()

				goldenFileMD5, err := hashFileMD5(tst.goldenFile)
				require.NoError(t, err)

				dstMD5, err := hashFileMD5(dst.Name())
				require.NoError(t, err)

				require.Equal(t, goldenFileMD5, dstMD5)
			}
		})
	}
}

func hashFileMD5(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
