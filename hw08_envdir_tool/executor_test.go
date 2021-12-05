package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("set remove env", func(t *testing.T) {
		_ = os.Setenv("DELETED", "deleted")
		environment := make(map[string]EnvValue)
		environment["HELLO_WORLD"] = EnvValue{Value: "hello world", NeedRemove: false}
		environment["DELETED"] = EnvValue{Value: "", NeedRemove: true}

		code := RunCmd([]string{"env"}, environment)
		require.Equal(t, 0, code)
		_, ok := os.LookupEnv("DELETED")
		require.False(t, ok)
		world, ok := os.LookupEnv("HELLO_WORLD")
		require.True(t, ok)
		require.Equal(t, "hello world", world)
	})

	t.Run("mkdir", func(t *testing.T) {
		environment := make(map[string]EnvValue)

		tempDir := os.TempDir()
		dir := path.Join(tempDir, "bbb", "ccc")
		defer func() {
			_ = os.RemoveAll(path.Join(tempDir, "bbb"))
		}()

		code := RunCmd([]string{"mkdir", "-p", dir}, environment)
		require.Equal(t, 0, code)

		fi, err := os.Stat(dir)
		require.Nil(t, err)
		require.True(t, fi.IsDir())
	})
}
