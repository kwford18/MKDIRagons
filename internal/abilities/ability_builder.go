package abilities

import (
	"github.com/kwford18/MKDIRagons/internal/race"
	"github.com/kwford18/MKDIRagons/template"
)

// BuildAbilityScores builds struct using values from base TemplateCharacter
func BuildAbilityScores(base *template.Character, race race.Race) AbilityScore {
	// Apply Racial bonus
OuterLoop:
	for _, ability := range race.AbilityBonuses {
		switch ability.AbilityScore.Name {
		case "STR":
			base.AbilityScores.Strength += ability.Bonus
		case "DEX":
			base.AbilityScores.Dexterity += ability.Bonus
		case "CON":
			base.AbilityScores.Constitution += ability.Bonus
		case "WIS":
			base.AbilityScores.Wisdom += ability.Bonus
		case "INT":
			base.AbilityScores.Intelligence += ability.Bonus
		case "CHA":
			base.AbilityScores.Charisma += ability.Bonus
		default:
			break OuterLoop
		}
	}

	return AbilityScore{
		Strength:     base.AbilityScores.Strength,
		Dexterity:    base.AbilityScores.Dexterity,
		Constitution: base.AbilityScores.Constitution,
		Wisdom:       base.AbilityScores.Wisdom,
		Intelligence: base.AbilityScores.Intelligence,
		Charisma:     base.AbilityScores.Charisma,
	}
}
