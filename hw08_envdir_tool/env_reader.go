package main

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrDirectoryNotExists = errors.New("directory not exists")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	fi, err := os.Stat(dir)
	if err != nil {
		return nil, ErrDirectoryNotExists
	}

	mode := fi.Mode()
	if !mode.IsDir() {
		return nil, ErrDirectoryNotExists
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var lineBreak string
	if runtime.GOOS == "windows" {
		lineBreak = "\r\n"
	} else {
		lineBreak = "\n"
	}

	environment := make(map[string]EnvValue)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if idx := strings.Index(name, "="); idx != -1 {
			continue
		}

		if file.Size() == 0 {
			environment[name] = EnvValue{Value: "", NeedRemove: true}
			continue
		}

		fromPath := path.Join(dir, file.Name())
		src, err := os.Open(fromPath)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(src)
		scanner.Scan()
		value := strings.TrimRight(scanner.Text(), " \t")
		vb := bytes.ReplaceAll([]byte(value), []byte{0x00}, []byte(lineBreak))
		value = string(vb)

		environment[name] = EnvValue{Value: value, NeedRemove: false}

		_ = src.Close()
	}

	return environment, nil
}
