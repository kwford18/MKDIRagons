package templates

import (
	"fmt"

	"github.com/kwford18/MKDIRagons/models"
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
	Level         int                   `toml:"level"`
	Race          string                `toml:"race"`
	Subrace       string                `toml:"subrace,omitempty"`
	Class         string                `toml:"class"`
	Subclass      string                `toml:"subclass,omitempty"`
	AbilityScores TemplateAbilityScores `toml:"ability_scores"`
	Proficiencies []string              `toml:"proficiencies"`
	Expertise     []string              `toml:"expertise,omitempty"`
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

func (t *TemplateCharacter) ProficiencyBonus() int {
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

func (t *TemplateCharacter) GetSkillAbility(name string) models.Ability {
	switch name {
	case "Athletics":
		return models.Strength
	case "Acrobatics", "SleightOfHand", "Stealth":
		return models.Dexterity
	case "Arcana", "History", "Investigation", "Nature", "Religion":
		return models.Intelligence
	case "AnimalHandling", "Insight", "Medicine", "Perception", "Survival":
		return models.Wisdom
	case "Deception", "Intimidation", "Performance", "Persuasion":
		return models.Charisma
	default:
		return models.Ability(0) // POTENTIAL CHANGE
	}
}

// Takes an ability and returns the modifier
func (ab *TemplateAbilityScores) Modifier(a models.Ability) int {
	var score int
	switch a {
	case models.Strength:
		score = ab.Strength
	case models.Dexterity:
		score = ab.Dexterity
	case models.Constitution:
		score = ab.Constitution
	case models.Intelligence:
		score = ab.Intelligence
	case models.Wisdom:
		score = ab.Wisdom
	case models.Charisma:
		score = ab.Charisma
	}
	return (score - 10) / 2
}
