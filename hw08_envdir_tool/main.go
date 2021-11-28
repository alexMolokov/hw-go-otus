package main

import (
	"errors"
	"log"
	"os"
)

const (
	minParam = 3
)

var ErrUnSupportedFormat = errors.New("unsupported format")

func main() {
	args := make([]string, len(os.Args))
	copy(args, os.Args)

	if len(args) < minParam {
		log.Fatal(ErrUnSupportedFormat)
	}

	env, err := ReadDir(args[1])
	if err != nil {
		log.Fatal(err)
	}

	RunCmd(args[2:], env)
}
