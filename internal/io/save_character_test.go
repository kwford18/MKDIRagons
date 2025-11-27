package io_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/character"
	"github.com/kwford18/MKDIRagons/internal/class"
	"github.com/kwford18/MKDIRagons/internal/inventory"
	"github.com/kwford18/MKDIRagons/internal/io"
	"github.com/kwford18/MKDIRagons/internal/race"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveJSON(t *testing.T) {
	// 1. Get the location of the test binary (just like SaveJSON does)
	exePath, err := os.Executable()
	require.NoError(t, err)
	exeDir := filepath.Dir(exePath)

	// 2. Define the relative path we will pass to the function
	relativePath := filepath.Join("testdata", "characters")

	// 3. Define the ABSOLUTE path where we expect the file to actually land
	// SaveJSON = exeDir + relativePath
	fullOutputPath := filepath.Join(exeDir, relativePath)

	// Clean up the actual location where files are written
	t.Cleanup(func() {
		_ = os.RemoveAll(fullOutputPath)
	})

	t.Run("successfully saves complex character to json file", func(t *testing.T) {
		// Arrange
		char := &character.Character{
			Name:  "tester",
			Level: 5,
			Race: race.Race{
				Name:  "Human",
				Speed: 30,
			},
			Class: class.Class{
				Name:   "Paladin",
				HitDie: 10,
			},
			Inventory: inventory.Inventory{
				Weapons: []inventory.Weapon{
					{BaseEquipment: inventory.BaseEquipment{Name: "Longsword", Index: "longsword"}},
				},
				Armor: []inventory.Armor{
					{BaseEquipment: inventory.BaseEquipment{Name: "Plate Armor", Index: "plate-armor"}},
				},
			},
			Proficiencies: []string{"Athletics", "Intimidation"},
		}

		// Calculate expected file path for verification
		expectedFileName := "tester.json"
		expectedFilePath := filepath.Join(fullOutputPath, expectedFileName)

		// Act: Pass the RELATIVE path. SaveJSON will prepend the exeDir.
		err := io.SaveJSON(char, relativePath)

		// Assert
		require.NoError(t, err, "SaveJSON should not return an error")
		assert.FileExists(t, expectedFilePath, "File should exist at the calculated executable path")

		// Verify Integrity: Read from the FULL path
		content, err := os.ReadFile(expectedFilePath)
		require.NoError(t, err, "Should be able to read the saved file")

		var savedChar character.Character
		err = json.Unmarshal(content, &savedChar)
		require.NoError(t, err, "Saved file should contain valid JSON")

		// Verify deep equality
		assert.Equal(t, char.Name, savedChar.Name)
		assert.Equal(t, char.Level, savedChar.Level)
		assert.Equal(t, char.Race.Name, savedChar.Race.Name)
		assert.Equal(t, char.Class.Name, savedChar.Class.Name)
		require.NotEmpty(t, savedChar.Inventory.Weapons)
		assert.Equal(t, "Longsword", savedChar.Inventory.Weapons[0].Name)
	})

	t.Run("overwrites existing file", func(t *testing.T) {
		char := &character.Character{
			Name:  "overwrite",
			Level: 1,
			Class: class.Class{Name: "Fighter"},
		}
		expectedFilePath := filepath.Join(fullOutputPath, "overwrite.json")

		// Act 1: Save initial version
		err := io.SaveJSON(char, relativePath)
		require.NoError(t, err)

		// Act 2: Update and save again
		char.Level = 99
		char.Class.Name = "Wizard"
		err = io.SaveJSON(char, relativePath)
		require.NoError(t, err)

		// Assert: Verify the file contains the NEW data
		content, err := os.ReadFile(expectedFilePath)
		require.NoError(t, err)

		var savedChar character.Character
		err = json.Unmarshal(content, &savedChar)
		require.NoError(t, err)

		assert.Equal(t, 99, savedChar.Level)
		assert.Equal(t, "Wizard", savedChar.Class.Name)
	})
}
