package templates

import (
	"fmt"

	"github.com/kwford18/MKDIRagons/fetch"
)

// Character represents the TOML structure
type TemplateCharacter struct {
	Name      string `toml:"name"`
	Race      string `toml:"race"`
	Subrace   string `toml:"subrace"`
	Class     string `toml:"class"`
	Subclass  string `toml:"subclass"`
	Equipment string `toml:"Equipment"`
	Spells    string `toml:"spells"`
}

type Character struct {
	Name      string
	Race      fetch.Race
	Class     fetch.Class
	Equipment fetch.Equipment
	Spells    fetch.Spell
}

func (c *Character) Print() {
	fmt.Printf("Name: %s\n", c.Name)
	fmt.Printf("Race: %s\n", c.Race.Name)
	fmt.Printf("Class: %s\n", c.Class.Name)
	fmt.Printf("Equipment: %v\n", c.Equipment.Name)
}

func (t *TemplateCharacter) Print() {
	fmt.Printf("Name: %s\n", t.Name)
	fmt.Printf("Race: %s\n", t.Race)
	fmt.Printf("Subrace: %s\n", t.Subrace)
	fmt.Printf("Class: %s\n", t.Class)
	fmt.Printf("Subclass: %s\n", t.Subclass)
	fmt.Printf("Equipment: %v\n", t.Equipment)
	fmt.Printf("Spells: %v\n", t.Spells)
}
