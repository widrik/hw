package main

import (
	"github.com/stretchr/testify/require"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/config"
	"testing"

)

func TestConfig(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		err := InitConfig("../../tests/testdata/config.json")
		require.Nil(t, err)
	})

	t.Run("validation error", func(t *testing.T) {
		err := InitConfig("../../tests/testdata/wrong_host.json")
		require.Equal(t, ErrWrongServerHost, err)
	})

	t.Run("nonexistent config", func(t *testing.T) {
		err := InitConfig("config.json")
		require.Equal(t, ErrCannotReadConfig, err)
	})

	t.Run("wrong config structure", func(t *testing.T) {
		err := InitConfig("../../tests/testdata/bad_structure.json")
		require.Equal(t, ErrCannotParseConfig, err)
	})
}
