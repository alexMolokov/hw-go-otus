package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("copy empty file", func(t *testing.T) {
		tempDir := os.TempDir()
		src, _ := os.CreateTemp(tempDir, "hw07")
		defer func() {
			_ = os.Remove(src.Name())
		}()

		dstFileName := "hw07-copy-empty-file.txt"
		dstPath := path.Join(tempDir, dstFileName)
		defer func() {
			_ = os.Remove(dstPath)
		}()
		err := Copy(src.Name(), dstPath, 0, 0)
		require.Nil(t, err)

		file, err := os.Open(dstPath)
		require.False(t, os.IsNotExist(err))

		fileInfo, err := file.Stat()
		require.Equal(t, int64(0), fileInfo.Size())
	})

	t.Run("unsupported file", func(t *testing.T) {
		tempDir := os.TempDir()

		dstFileName := "hw08-copy-empty-file.txt"
		dstPath := path.Join(tempDir, dstFileName)
		defer func() {
			_ = os.Remove(dstPath)
		}()
		err := Copy("/dev/random", dstPath, 0, 100)
		require.NotNil(t, err)
		require.True(t, err == ErrUnsupportedFile)
	})
}
