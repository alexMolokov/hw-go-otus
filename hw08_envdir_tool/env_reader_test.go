package main

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("read not exists directory", func(t *testing.T) {
		tempDir := os.TempDir()
		dirName := path.Join(tempDir, "not_exists_directory")
		si, err := os.Stat(dirName)
		if err == nil && si.IsDir() {
			_ = os.Remove(dirName)
		}
		_, err = ReadDir(dirName)
		require.True(t, errors.Is(err, ErrDirectoryNotExists))
	})

	t.Run("symbol = in file name", func(t *testing.T) {
		tempDir := os.TempDir()
		dirName, err := os.MkdirTemp(tempDir, "exists_directory")
		if err == nil {
			defer func() {
				_ = os.RemoveAll(dirName)
			}()
		}

		_, _ = os.CreateTemp(dirName, "a=bc")
		result, err := ReadDir(dirName)
		require.Nil(t, err)
		require.Equal(t, 0, len(result))

		_, _ = os.CreateTemp(dirName, "abc")
		result, err = ReadDir(dirName)
		require.Nil(t, err)
		require.Equal(t, 1, len(result))
	})
}
