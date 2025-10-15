package templates

import (
	"fmt"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
)

// Helper to validate a single score
func validateScore(name string, score int) error {
	if score < 0 || score > 20 {
		return fmt.Errorf("invalid ability score for %s -> %d. Must be in range [0, 20]", name, score)
	}
	return nil
}

// Validate all ability scores
func (ta *TemplateAbilityScores) Validate() error {
	if err := validateScore("Strength", ta.Strength); err != nil {
		return err
	}
	if err := validateScore("Dexterity", ta.Dexterity); err != nil {
		return err
	}
	if err := validateScore("Constitution", ta.Constitution); err != nil {
		return err
	}
	if err := validateScore("Intelligence", ta.Intelligence); err != nil {
		return err
	}
	if err := validateScore("Wisdom", ta.Wisdom); err != nil {
		return err
	}
	if err := validateScore("Charisma", ta.Charisma); err != nil {
		return err
	}
	return nil
}

// verifyTOML checks relevant fields parsed into TemplateCharacter and ensures they are valid
func verifyTOML(t TemplateCharacter) error {
	// Validate character level
	if t.Level < 1 || t.Level > 20 {
		return fmt.Errorf("invalid level")
	}

	// Validate 5e 2014 race
	valid5eRace := []string{
		"dragonborn",
		"dwarf",
		"elf",
		"gnome",
		"half-elf",
		"half-orc",
		"halfling",
		"human",
		"tiefling",
	}
	baseRace := strings.ToLower(t.Race)
	if !slices.Contains(valid5eRace, baseRace) {
		return fmt.Errorf("no valid 5e 2014 race provided")
	}

	// Validate 5e 2014 class
	validClass := []string{
		"barbarian",
		"bard",
		"cleric",
		"druid",
		"fighter",
		"monk",
		"paladin",
		"ranger",
		"rogue",
		"sorcerer",
		"warlock",
		"wizard",
	}
	baseClass := strings.ToLower(t.Class)
	if !slices.Contains(validClass, baseClass) {
		return fmt.Errorf("no valid 5e 2014 class provided")
	}

	// Validate ability scores
	if err := t.AbilityScores.Validate(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// TomlParse parses the provided TOML file into the Template struct
func TomlParse(fileName string) (TemplateCharacter, error) {
	var t TemplateCharacter
	_, err := toml.DecodeFile(fileName, &t)
	if err != nil {
		return t, fmt.Errorf("failed to parse file: %w", err)
	}

	err = verifyTOML(t)
	if err != nil {
		return TemplateCharacter{}, err
	}

	return t, nil
}
