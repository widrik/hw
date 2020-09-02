package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	t.Run("validate correct in validator", func(t *testing.T) {
		res, err := createValidators("string", "`validate:\"in:admin,stuff\"`")
		require.NoError(t, err)
		require.EqualValues(t, []Validator{
			{
				Type:  "in",
				Value: []string{"admin", "stuff"},
			},
		}, res)
	})

	t.Run("empty data", func(t *testing.T) {
		t.Run("empty tag", func(t *testing.T) {
			validators, err := createValidators("", "")
			require.NoError(t, err)
			require.Equal(t, len(validators), 0)
		})

		t.Run("tag", func(t *testing.T) {
			_, err := createValidators("", "``")
			require.NoError(t, err)
		})
		t.Run("validate tag", func(t *testing.T) {
			_, err := createValidators("", "`validate:\"\"`")
			require.NoError(t, err)

		})
	})


}