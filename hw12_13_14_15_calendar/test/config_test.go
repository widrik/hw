package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/config"
)

func TestConfig(t *testing.T) {
	t.Run("Empty path", func(t *testing.T) {
		_, err := config.Init("")
		require.Equal(t, config.ErrFilePathEmpty, err)
	})

	t.Run("Wrong path", func(t *testing.T) {
		_, err := config.Init("blabla")
		require.Equal(t, config.ErrReadFile, err)
	})

	t.Run("Wrong json data", func(t *testing.T) {
		_, err := config.Init("testdata/config/wrongData.json")
		require.Equal(t, config.ErrReadFile, err)
	})

	t.Run("Empty json data", func(t *testing.T) {
		_, err := config.Init("testdata/config/empty.json")
		require.Equal(t, config.ErrReadFile, err)
	})

	t.Run("Correct data", func(t *testing.T) {
		configuration, err := config.Init("testdata/config/correctData.json")
		require.Equal(t, nil, err)
		require.Equal(t, configuration.HTTPServer.Host, "test_host")
	})

}
