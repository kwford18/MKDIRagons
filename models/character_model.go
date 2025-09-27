package models

import "fmt"

type Character struct {
	Name          string
	Level         int
	Race          Race
	Class         Class
	Proficiencies []string
	AbilityScores AbilityScore
	Skills        SkillList
	Inventory     Inventory
	Spells        [][]Spell
}

func (c *Character) ProficiencyBonus() int {
	switch c.Level {
	case 1, 2, 3, 4:
		return 2
	case 5, 6, 7, 8:
		return 3
	case 9, 10, 11, 12:
		return 4
	case 13, 14, 15, 16:
		return 5
	case 17, 18, 19, 20:
		return 6
	default:
		return -1
	}
}

func (c *Character) GetSkillAbility(name string) Ability {
	switch name {
	case "Athletics":
		return Strength
	case "Acrobatics", "SleightOfHand", "Stealth":
		return Dexterity
	case "Arcana", "History", "Investigation", "Nature", "Religion":
		return Intelligence
	case "AnimalHandling", "Insight", "Medicine", "Perception", "Survival":
		return Wisdom
	case "Deception", "Intimidation", "Performance", "Persuasion":
		return Charisma
	default:
		return Ability(0) // POTENTIAL CHANGE
	}
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
	fmt.Printf("    - Strength:     %d\n", c.AbilityScores.Strength)
	fmt.Printf("    - Dexterity:    %d\n", c.AbilityScores.Dexterity)
	fmt.Printf("    - Constitution: %d\n", c.AbilityScores.Constitution)
	fmt.Printf("    - Wisdom:       %d\n", c.AbilityScores.Wisdom)
	fmt.Printf("    - Intelligence: %d\n", c.AbilityScores.Intelligence)
	fmt.Printf("    - Charisma:     %d\n", c.AbilityScores.Charisma)

	// Proficiencies
	fmt.Println("Proficiencies:")
	for _, prof := range c.Proficiencies {
		fmt.Printf("	- %s\n", prof)
	}
}
