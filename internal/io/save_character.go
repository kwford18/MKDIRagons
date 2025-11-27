package io

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kwford18/MKDIRagons/internal/character"
)

func SaveJSON(character *character.Character, dirPath string) error {
	// Get the path of the currently running executable binary
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get the directory containing the executable
	exeDir := filepath.Dir(exePath)

	// Combine the Executable's directory with the requested dirPath
	// If dirPath is relative (e.g. "characters/"), it becomes "/path/to/exe/characters/"
	// If dirPath is absolute (e.g. "/tmp/chars"), filepath.Join will respect the absolute path and ignore exeDir
	fullDirPath := filepath.Join(exeDir, dirPath)
	cleanDir := filepath.Clean(fullDirPath)

	// Ensure directory exists
	if err := os.MkdirAll(cleanDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", cleanDir, err)
	}

	// Generate lowercase file name
	fileName := strings.ToLower(character.Name) + ".json"

	// Join directory and filename
	filePath := filepath.Join(cleanDir, fileName)

	// Convert character struct to pretty JSON
	data, err := json.MarshalIndent(character, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal character: %w", err)
	}

	// Create file if new, Write Only, Truncate if exists
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	// Write the data
	if _, err := file.Write(data); err != nil {
		return err
	}

	fmt.Printf("✓ Character saved to: %s\n", filePath)
	return nil
}
