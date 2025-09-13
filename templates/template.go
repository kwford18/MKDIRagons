package templates

import (
	"fmt"
)

// Template character parsed from a TOML file

type TemplateAbilityScores struct {
	Strength     int
	Dexterity    int
	Constitution int
	Wisdom       int
	Intelligence int
	Charisma     int
}

type TemplateInventory struct {
	Weapons []string
	Armor   []string
	Items   []string
}

type TemplateSpells struct {
	Level [][]string
}

type TemplateCharacter struct {
	Name          string                `toml:"name"`
	Race          string                `toml:"race"`
	Subrace       string                `toml:"subrace"`
	Class         string                `toml:"class"`
	Subclass      string                `toml:"subclass"`
	AbilityScores TemplateAbilityScores `toml:"ability_scores"`
	Proficiencies []string              `toml:"proficiencies"`
	Inventory     TemplateInventory     `toml:"inventory"`
	Spells        TemplateSpells        `toml:"spells"`
}

func (t *TemplateCharacter) Print() {
	fmt.Printf("Name: %s\n", t.Name)
	fmt.Printf("Race: %s\n", t.Race)
	fmt.Printf("Subrace: %s\n", t.Subrace)
	fmt.Printf("Class: %s\n", t.Class)
	fmt.Printf("Subclass: %s\n", t.Subclass)
	fmt.Printf("Ability Scores: %v\n", t.AbilityScores)
	fmt.Printf("Proficiencies: %v\n", t.Proficiencies)
	fmt.Printf("Inventory: %v\n", t.Inventory)
	fmt.Printf("Spells: %v\n", t.Spells)
}
