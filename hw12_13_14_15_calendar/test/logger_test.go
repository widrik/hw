package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/logging"
)

func TestLogger(t *testing.T) {
	t.Run("Empty file", func(t *testing.T) {
		err := logging.Init("debug", "")
		require.Error(t, err)
	})

	t.Run("Not valid level", func(t *testing.T) {
		err := logging.Init("blabla", "testdata/test.txt")
		require.Error(t, err)
	})

	t.Run("Ok", func(t *testing.T) {
		err := logging.Init("error", "testdata/test.txt")
		require.NoError(t, err)
	})
}
