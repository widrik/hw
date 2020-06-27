package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRunCmd(t *testing.T) {

	t.Run("Empty cmd and env", func(t *testing.T) {
		var stringsEmptySlice []string
		r := RunCmd(stringsEmptySlice, Environment{})
		require.Equal(t, -1, r)
	})

	t.Run("Empty env", func(t *testing.T) {
		r := RunCmd([]string{"ls"}, Environment{})
		require.Equal(t, 0, r)
	})
}
