package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func getSrcFile() (*os.File, error) {
	src, err := os.Open(from)
	if err != nil {
		if os.IsNotExist(err) {
			return src, fmt.Errorf("file not exists %s. Nothing copying", from)
		}

		return src, fmt.Errorf("error open file %s", from)
	}

	srcInfo, err := src.Stat()
	if err != nil {
		return src, fmt.Errorf("can't get file information %s. Nothing copying", from)
	}

	size := srcInfo.Size()
	if offset < 0 || size < offset {
		return src, fmt.Errorf("offset must be less than file %d bytes and more than 0", size)
	}

	src.Seek(offset, 0)

	return src, nil
}

func getDstFile() (*os.File, error) {
	dst, err := os.OpenFile(to, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return dst, fmt.Errorf("can't open file to write %s", to)
	}
	return dst, nil
}

func main() {
	flag.Parse()

	src, err := getSrcFile()
	defer func() {
		_ = src.Close()
	}()
	if err != nil {
		log.Fatal(err)
	}

	dst, err := getDstFile()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = dst.Close()
	}()

	if limit < 0 {
		err := fmt.Errorf("param limit must be more or equal 0. Get %d", limit)
		log.Fatal(err)
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

	bar := pb.Full.Start64(pbLimit)
	barReader := bar.NewProxyReader(reader)
	io.Copy(dst, barReader)
	bar.Finish()
}
