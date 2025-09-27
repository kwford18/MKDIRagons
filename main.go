package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/kwford18/MKDIRagons/builder"
	"github.com/kwford18/MKDIRagons/models"
	"github.com/kwford18/MKDIRagons/templates"
)

// Builds the file path
func fileArgs() (string, error) {
	if len(os.Args) > 2 {
		return "Too many args", fmt.Errorf("too many args: %s", os.Args)
	} else if len(os.Args) < 2 {
		return "Too few args", fmt.Errorf("too few args: %s", os.Args)
	}

	// Load file
	path := "templates/toml-characters/"

	// Check if file extension was provided
	if filepath.Ext(os.Args[1]) != ".toml" && filepath.Ext(os.Args[1]) != "" {
		return "Incorrect file type", fmt.Errorf("incorrect file type: %s", filepath.Ext(os.Args[1]))
	}

	file := path + os.Args[1]
	return file, nil
}

// Parses the provided TOML file into the Template struct
func tomlParse(fileName string) (templates.TemplateCharacter, error) {
	var t templates.TemplateCharacter
	_, err := toml.DecodeFile(fileName, &t)
	if err != nil {
		return t, fmt.Errorf("failed to parse file: %w", err)
	}

	return t, nil
}

func saveJSON(character *models.Character) error {
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

func main() {
	err := templates.GenerateEmptyTOML()
	if err != nil {
		log.Fatalf("Error generating empty TOML: %v\n", err)
	}

	fileName, err := fileArgs()
	if err != nil {
		log.Fatalf("Incorrect file arguments: %v\n", err)
	}

	base, err := tomlParse(fileName)
	if err != nil {
		log.Fatalf("Failed to parse file: %v\n", err)
	}

	// Print base character
	fmt.Println("TOML BASE:")
	base.Print()
	fmt.Println()

	character, err := builder.BuildCharacter(&base)
	if err != nil {
		log.Fatalf("Error building character: %v\n", err)
	}

	// Print built character
	fmt.Println("BUILT CHARACTER:")
	character.Print()

	if err := saveJSON(character); err != nil {
		log.Fatalf("Failed to save character as JSON\n")
	}
}
