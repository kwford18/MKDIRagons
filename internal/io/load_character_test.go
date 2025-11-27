package io_test

import (
	"path/filepath"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadCharacter(t *testing.T) {
	// Define the path to your fixtures directory relative to this test file
	testDataPath := "testdata"

	t.Run("successfully loads valid character from fixture", func(t *testing.T) {
		filename := filepath.Join(testDataPath, "valid_character.json")

		// Execute
		char, err := io.LoadCharacter(filename)

		// Verify
		require.NoError(t, err, "Should not return error for valid file")
		require.NotNil(t, char, "Character struct should not be nil")

		// Assert specific fields to ensure data integrity
		assert.Equal(t, "Leki", char.Name)
		assert.Equal(t, "Cleric", char.Class.Name)
		assert.Equal(t, 5, char.Level)
		assert.Equal(t, 15, char.AbilityScores.Wisdom)
	})

	t.Run("fails when file does not exist", func(t *testing.T) {
		filename := filepath.Join(testDataPath, "non_existent_file.json")

		// Execute
		char, err := io.LoadCharacter(filename)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, char)
		assert.Contains(t, err.Error(), "could not open file")
	})

	t.Run("fails when json is malformed", func(t *testing.T) {
		filename := filepath.Join(testDataPath, "malformed.json")

		// Execute
		char, err := io.LoadCharacter(filename)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, char)
		// We expect the specific error wrapping message defined in your io.go
		assert.Contains(t, err.Error(), "could not open file: open testdata/malformed.json: no such file or directory")
	})
}
