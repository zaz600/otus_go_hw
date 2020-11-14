package main

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSrcFileNotFound       = errors.New("src file was not wound")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	src, err := openSrc(fromPath, offset, limit)
	if err != nil {
		return err
	}
	defer src.file.Close()

	dst, err := openDst(toPath)
	if err != nil {
		return err
	}
	defer dst.file.Close()

	fmt.Printf("'%s' -> '%s' %d bytes\n", fromPath, toPath, src.copySize)
	fmt.Printf("Offset: %d, limit: %d\n", offset, limit)

	p := pb.New64(src.copySize)
	p.Start()

	var chunkSize int64 = 1024
	var total int64 = 0
	for {
		if chunkSize > src.copySize {
			chunkSize = src.copySize
		}
		// fmt.Printf("chunk=%d\n", chunkSize)

		n, err := io.CopyN(dst.file, src.file, chunkSize)
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("error while coping '%s' to '%s': %w", fromPath, toPath, err)
		}

		if n == 0 {
			break
		}

		total += n
		if total >= src.copySize {
			break
		}

		p.Add64(n)
		time.Sleep(100 * time.Microsecond)
	}

	p.Finish()

	fmt.Println("Copy complete")

	return nil
}
