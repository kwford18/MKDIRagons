package template

import (
	"fmt"
	"github.com/kwford18/MKDIRagons/internal/core"
	"math"
)

// Template character parsed from a TOML file

type AbilityScores struct {
	Strength     int
	Dexterity    int
	Constitution int
	Wisdom       int
	Intelligence int
	Charisma     int
}

type Inventory struct {
	Weapons []string
	Armor   []string
	Items   []string
}

type Spells struct {
	Level [][]string
}

type Character struct {
	Name          string        `toml:"name"`
	Level         int           `toml:"level"`
	Race          string        `toml:"race"`
	Subrace       string        `toml:"subrace,omitempty"`
	Class         string        `toml:"class"`
	Subclass      string        `toml:"subclass,omitempty"`
	AbilityScores AbilityScores `toml:"ability_scores"`
	Proficiencies []string      `toml:"proficiencies"`
	Expertise     []string      `toml:"expertise,omitempty"`
	Inventory     Inventory     `toml:"inventory"`
	Spells        Spells        `toml:"spells"`
}

func (t *Character) Print() {
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

func (t *Character) ProficiencyBonus() int {
	switch t.Level {
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

// GetSkillAbility takes a skill name and returns which ability score it uses
func (t *Character) GetSkillAbility(name string) core.Ability {
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

// Modifier takes an ability and returns the modifier
func (t AbilityScores) Modifier(a core.Ability) int {
	var score int
	switch a {
	case core.Strength:
		score = t.Strength
	case core.Dexterity:
		score = t.Dexterity
	case core.Constitution:
		score = t.Constitution
	case core.Intelligence:
		score = t.Intelligence
	case core.Wisdom:
		score = t.Wisdom
	case core.Charisma:
		score = t.Charisma
	}
	return int(math.Floor(float64(score-10) / 2.0))
}
