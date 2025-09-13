package templates

import (
	"fmt"

	"github.com/kwford18/MKDIRagons/fetch"
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

// Character represents the TOML structure
type Character struct {
	Name          string
	Race          fetch.Race
	Class         fetch.Class
	Proficiencies []string
	AbilityScores map[string]int
	Inventory     fetch.Inventory
	Spells        [][]fetch.Spell
}

func (c *Character) Print() {
	// Name, Race, Class
	fmt.Printf("Name: %s\n", c.Name)
	fmt.Printf("Race: %s\n", c.Race.Name)
	fmt.Printf("Class: %s\n", c.Class.Name)

	// Equipment
	fmt.Println("Equipment:")

	fmt.Printf("    - Armor: \n")
	for _, armor := range c.Inventory.Armor {
		fmt.Printf("	- %s\n", armor.Name)
	}

	fmt.Printf("    - Weapons: \n")
	for _, weapons := range c.Inventory.Weapons {
		fmt.Printf("	- %s\n", weapons.Name)
	}

	fmt.Printf("    - Equipment: \n")
	for _, items := range c.Inventory.Items {
		fmt.Printf("	- %s\n", items.Name)
	}

	// Spells
	fmt.Println("Spells:")
	for i, level := range c.Spells {
		if len(level) == 0 {
			continue
		}
		fmt.Printf("  Level %d:\n", i)
		for _, spell := range level {
			fmt.Printf("    - %s\n", spell.Name)
		}
	}

	// Ability Scores
	fmt.Println("Ability Scores:")
	fmt.Printf("    - Strength:     %d\n", c.AbilityScores["str"])
	fmt.Printf("    - Dexterity:    %d\n", c.AbilityScores["dex"])
	fmt.Printf("    - Constitution: %d\n", c.AbilityScores["con"])
	fmt.Printf("    - Wisdom:       %d\n", c.AbilityScores["wis"])
	fmt.Printf("    - Intelligence: %d\n", c.AbilityScores["int"])
	fmt.Printf("    - Charisma:     %d\n", c.AbilityScores["cha"])

	// Proficiencies
	fmt.Println("Proficiencies:")
	for _, prof := range c.Proficiencies {
		fmt.Printf("	- %s\n", prof)
	}
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
