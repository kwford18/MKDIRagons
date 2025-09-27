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
		Level:    1,
		Race:     "",
		Subrace:  "",
		Class:    "",
		Subclass: "",
		AbilityScores: TemplateAbilityScores{
			Strength:     10,
			Dexterity:    10,
			Constitution: 10,
			Wisdom:       10,
			Intelligence: 10,
			Charisma:     10,
		},
		Proficiencies: []string{},
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
		return fmt.Errorf("error creating directory: %w", err)
	}

	// Create file if it doesn't exit (or truncate if it does)
	file, err := os.Create("templates/template.toml")
	if err != nil {
		return fmt.Errorf("error creating template file: %w", err)
	}
	defer file.Close()

	// Encode struct as TOML
	if err := toml.NewEncoder(file).Encode(template); err != nil {
		return fmt.Errorf("error encoding toml file: %w", err)
	}

	return nil
}
