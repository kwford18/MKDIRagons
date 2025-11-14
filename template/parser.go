package template

import (
	"fmt"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
)

// Helper to validate a single score
func validateScore(name string, score int) error {
	if score < 0 {
		return fmt.Errorf("ability score for %s is too low -> %d. Must be in range [0, 20]", name, score)
	} else if score > 20 {
		return fmt.Errorf("ability score for %s is too high -> %d. Must be in range [0, 20]", name, score)
	}
	return nil
}

// Validate all ability scores
func (t AbilityScores) Validate() error {
	if err := validateScore("Strength", t.Strength); err != nil {
		return err
	}
	if err := validateScore("Dexterity", t.Dexterity); err != nil {
		return err
	}
	if err := validateScore("Constitution", t.Constitution); err != nil {
		return err
	}
	if err := validateScore("Intelligence", t.Intelligence); err != nil {
		return err
	}
	if err := validateScore("Wisdom", t.Wisdom); err != nil {
		return err
	}
	if err := validateScore("Charisma", t.Charisma); err != nil {
		return err
	}
	return nil
}

// verifyTOML checks relevant fields parsed into TemplateCharacter and ensures they are valid
func verifyTOML(t Character) error {
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

// TomlParse Parses the provided TOML file into the Template struct
func TomlParse(fileName string) (Character, error) {
	var t Character
	_, err := toml.DecodeFile(fileName, &t)
	if err != nil {
		return t, fmt.Errorf("failed to parse file: %w", err)
	}

	err = verifyTOML(t)
	if err != nil {
		return Character{}, err
	}

	return t, nil
}
