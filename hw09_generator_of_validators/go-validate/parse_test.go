package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Run("wrong type", func(t *testing.T) {
		require.Equal(t, "1", "1")
	})

}