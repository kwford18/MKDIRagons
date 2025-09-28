package io

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kwford18/MKDIRagons/internal/character"
)

func LoadCharacter(filename string) (*character.Character, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	// Decode JSON into the struct
	var char character.Character
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&char); err != nil {
		return nil, fmt.Errorf("could not decode JSON: %w", err)
	}

	return &char, nil
}
