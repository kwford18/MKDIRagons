package character

import (
	"fmt"

	"github.com/kwford18/MKDIRagons/internal/abilities"
	"github.com/kwford18/MKDIRagons/internal/class"
	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/internal/inventory"
	"github.com/kwford18/MKDIRagons/internal/race"
	"github.com/kwford18/MKDIRagons/internal/skills"
	"github.com/kwford18/MKDIRagons/internal/spells"
	"github.com/kwford18/MKDIRagons/internal/stats"
)

type Character struct {
	Name          string                 `json:"name"`
	Level         int                    `json:"level"`
	Race          race.Race              `json:"race"`
	Class         class.Class            `json:"class"`
	Stats         stats.Stats            `json:"stats"`
	Proficiencies []string               `json:"proficiencies"`
	AbilityScores abilities.AbilityScore `json:"ability_scores"`
	Skills        skills.SkillList       `json:"skills"`
	Inventory     inventory.Inventory    `json:"inventory"`
	Spells        [][]spells.Spell       `json:"spells"`
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

func (c *Character) GetSkillAbility(name string) core.Ability {
	switch name {
	case "Athletics":
		return core.Strength
	case "Acrobatics", "SleightOfHand", "Stealth":
		return core.Dexterity
	case "Arcana", "History", "Investigation", "Nature", "Religion":
		return core.Intelligence
	case "AnimalHandling", "Insight", "Medicine", "Perception", "Survival":
		return core.Wisdom
	case "Deception", "Intimidation", "Performance", "Persuasion":
		return core.Charisma
	default:
		return core.Ability(0) // POTENTIAL CHANGE
	}
}

func (c *Character) Print() {
	// Name, Level, Race, Class, Basic Stats
	fmt.Printf("Name: %s\n", c.Name)
	fmt.Printf("Level: %d\n", c.Level)
	c.Race.Print()
	c.Class.Print()
	c.Stats.Print()

	fmt.Println()

	// Equipment
	c.Inventory.Print()

	fmt.Println()

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

	fmt.Println()

	// Ability Scores
	c.AbilityScores.Print()

	// Proficiencies
	fmt.Println("Proficiencies:")
	for _, prof := range c.Proficiencies {
		fmt.Printf("	- %s\n", prof)
	}

	fmt.Println()

	// Skills
	c.Skills.Print()
}
