package main

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestReadDir(t *testing.T) {

	t.Run("Empty dir", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "prefix")
		if err != nil {
			log.Fatal(err)
		}

		env, err := ReadDir(dir)
		require.Equal(t, len(env), 0)

		defer os.RemoveAll(dir)
	})

	t.Run("Wrong dir path", func(t *testing.T) {
		_, err := ReadDir("blabla")
		require.Error(t, err)
	})

	t.Run("Empty file", func(t *testing.T) {
		dir, err := ioutil.TempDir("dir", "prefix")
		file, err := ioutil.TempFile(dir, "filename.txt")
		if err != nil {
			log.Fatal(err)
		}

		env, err := ReadDir(dir)

		require.Equal(t, len(env), 0)

		defer os.Remove(file.Name())
		defer os.RemoveAll(dir)
	})



}
