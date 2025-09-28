package main

import (
	"fmt"
	"log"

	"github.com/kwford18/MKDIRagons/internal/character"
	"github.com/kwford18/MKDIRagons/internal/io"

	"github.com/kwford18/MKDIRagons/templates"
)

func main() {
	err := templates.GenerateEmptyTOML()
	if err != nil {
		log.Fatalf("Error generating empty TOML: %v\n", err)
	}

	fileName, err := io.FileArgs()
	if err != nil {
		log.Fatalf("Incorrect file arguments: %v\n", err)
	}

	base, err := templates.TomlParse(fileName)
	if err != nil {
		log.Fatalf("Failed to parse file: %v\n", err)
	}

	// Print base character
	fmt.Println("TOML BASE:")
	base.Print()
	fmt.Println()

	char, err := character.BuildCharacter(&base)
	if err != nil {
		log.Fatalf("Error building character: %v\n", err)
	}

	// Print built character
	fmt.Println("BUILT CHARACTER:")
	char.Print()

	if err := io.SaveJSON(char); err != nil {
		log.Fatalf("Failed to save character as JSON\n")
	}
}
