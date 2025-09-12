package templates

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

func GenerateEmptyTOML() error {
	// Spell 2D array
	levels := make([][]string, 10)
	for i := range levels {
		levels[i] = []string{}
	}

	template := TemplateCharacter{
		Name:     "",
		Race:     "",
		Subrace:  "",
		Class:    "",
		Subclass: "",
		Inventory: TemplateInventory{
			Armor:   []string{},
			Weapons: []string{},
			Items:   []string{},
		},
		Spells: TemplateSpells{
			Level: levels,
		},
	}

	// Create template directory if it doesn't exist already
	if err := os.MkdirAll("templates", 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return err
	}

	// Create file if it doesn't exit (or truncate if it does)
	file, err := os.Create("templates/template.toml")
	if err != nil {
		fmt.Printf("Error creating template file: %v\n", err)
		return err
	}
	defer file.Close()

	// Encode struct as TOML
	if err := toml.NewEncoder(file).Encode(template); err != nil {
		fmt.Printf("Error encoding toml file: %v\n", err)
		return err
	}

	return nil
}
