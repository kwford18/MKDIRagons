package models

import "fmt"

type Character struct {
	Name          string
	Race          Race
	Class         Class
	Proficiencies []string
	AbilityScores map[string]int
	Inventory     Inventory
	Spells        [][]Spell
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
