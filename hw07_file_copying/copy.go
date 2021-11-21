package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	src, err := os.Open(fromPath)
	defer func() {
		_ = src.Close()
	}()
	if err != nil {
		return err
	}

	srcInfo, err := src.Stat()
	if err != nil {
		return fmt.Errorf("can't get file information %s. Nothing copying", fromPath)
	}

	mode := srcInfo.Mode()
	if !mode.IsRegular() {
		return ErrUnsupportedFile
	}

	if offset < 0 {
		return errors.New("offset must be  more or equal 0")
	}

	if size := srcInfo.Size(); size < offset {
		return ErrOffsetExceedsFileSize
	}

	if offset > 0 {
		src.Seek(offset, io.SeekStart)
	}

	dst, err := os.Create(toPath)
	defer func() {
		_ = dst.Close()
	}()
	if err != nil {
		return err
	}

	pbLimit := limit
	var reader io.Reader
	if limit == 0 {
		reader = src
		srcInfo, _ := src.Stat()
		pbLimit = srcInfo.Size()
	} else {
		reader = io.LimitReader(src, limit)
	}

	io.Copy(dst, reader)

	bar := pb.Full.Start64(pbLimit)
	barReader := bar.NewProxyReader(reader)
	io.Copy(dst, barReader)
	bar.Finish()

	return nil
}
