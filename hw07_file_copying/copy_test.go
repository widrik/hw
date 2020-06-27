package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("File does not exist", func(t *testing.T) {
		err := Copy("blabla.txt", "/testdata", 0, 0)
		require.NotNil(t, err)
	})

	t.Run("Empty file name", func(t *testing.T) {
		err := Copy("", "/testdata", 0, 0)
		require.NotNil(t, err)
	})

	t.Run("Error offset exceeds file size", func(t *testing.T) {
		tmpfile, err := ioutil.TempFile("", "example")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		content := []byte("")
		if _, err := tmpfile.Write(content); err != nil {
			log.Fatal(err)
		}

		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}

		err = Copy(tmpfile.Name(), "", 10000, 0)
		require.EqualError(t, err, ErrOffsetExceedsFileSize.Error())
	})
}
