package io

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/kwford18/MKDIRagons/internal/character"
)

func SaveJSON(character *character.Character) error {
	// Ensure ./characters exists
	dir := "characters"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	fileName := character.Name + ".json"

	// Build file path inside ./characters
	filePath := filepath.Join(dir, fileName)

	// Convert to pretty JSON
	data, err := json.MarshalIndent(character, "", "  ")
	if err != nil {
		return err
	}

	// Open file with truncate + write-only + create
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write JSON
	if _, err := file.Write(data); err != nil {
		return err
	}

	// Successfully saved file
	return nil
}
