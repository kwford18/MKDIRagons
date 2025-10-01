package abilities

import (
	"github.com/kwford18/MKDIRagons/internal/class"
	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/templates"
)

func BuildSavingThrows(base *templates.TemplateCharacter, abilities AbilityScore, charClass *class.Class) AbilityScore {
	profBonus := base.ProficiencyBonus()

	// start from modifiers for all abilities
	saves := AbilityScore{
		Strength:     abilities.Modifier(core.Strength),
		Dexterity:    abilities.Modifier(core.Dexterity),
		Constitution: abilities.Modifier(core.Constitution),
		Intelligence: abilities.Modifier(core.Intelligence),
		Wisdom:       abilities.Modifier(core.Wisdom),
		Charisma:     abilities.Modifier(core.Charisma),
	}

	// add proficiency bonus for proficient saves
	for _, bonus := range charClass.SavingThrows {
		switch bonus.Name {
		case "STR":
			saves.Strength += profBonus
		case "DEX":
			saves.Dexterity += profBonus
		case "CON":
			saves.Constitution += profBonus
		case "WIS":
			saves.Wisdom += profBonus
		case "INT":
			saves.Intelligence += profBonus
		case "CHA":
			saves.Charisma += profBonus
		}
	}

	return saves
}
