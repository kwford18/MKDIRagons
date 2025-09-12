package templates

import (
	"fmt"

	"github.com/kwford18/MKDIRagons/fetch"
)

// Template character parsed from a TOML file
type TemplateInventory struct {
	Weapons []string
	Armor   []string
	Items   []string
}

type TemplateSpells struct {
	Level [][]string
}

type TemplateCharacter struct {
	Name      string            `toml:"name"`
	Race      string            `toml:"race"`
	Subrace   string            `toml:"subrace"`
	Class     string            `toml:"class"`
	Subclass  string            `toml:"subclass"`
	Inventory TemplateInventory `toml:"inventory"`
	Spells    TemplateSpells    `toml:"spells"`
}

// Character represents the TOML structure
type Character struct {
	Name      string
	Race      fetch.Race
	Class     fetch.Class
	Inventory fetch.Inventory
	Spells    [][]fetch.Spell
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
}

func (t *TemplateCharacter) Print() {
	fmt.Printf("Name: %s\n", t.Name)
	fmt.Printf("Race: %s\n", t.Race)
	fmt.Printf("Subrace: %s\n", t.Subrace)
	fmt.Printf("Class: %s\n", t.Class)
	fmt.Printf("Subclass: %s\n", t.Subclass)
	fmt.Printf("Inventory: %v\n", t.Inventory)
	fmt.Printf("Spells: %v\n", t.Spells)
}
