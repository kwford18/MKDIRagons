package builder

import (
	"github.com/kwford18/MKDIRagons/models"
	"github.com/kwford18/MKDIRagons/templates"
)

// Build AbilityScore struct using values from base TemplateCharacter
func buildAbilityScores(base *templates.TemplateCharacter, race models.Race) models.AbilityScore {
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

	return models.AbilityScore{
		Strength:     base.AbilityScores.Strength,
		Dexterity:    base.AbilityScores.Dexterity,
		Constitution: base.AbilityScores.Constitution,
		Wisdom:       base.AbilityScores.Wisdom,
		Intelligence: base.AbilityScores.Intelligence,
		Charisma:     base.AbilityScores.Charisma,
	}
}

func buildSkill(base *templates.TemplateCharacter, name string) models.Skill {
	bonus := base.AbilityScores.Modifier(base.GetSkillAbility(name))
	proficient := false
	expert := false

	for _, prof := range base.Proficiencies {
		if prof == name {
			bonus += base.ProficiencyBonus()
			proficient = true
			break
		}
	}
	for _, exp := range base.Expertise {
		if exp == name {
			bonus += base.ProficiencyBonus()
			expert = true
			break
		}
	}

	return models.Skill{
		Name:       name,
		Bonus:      bonus,
		Ability:    base.GetSkillAbility(name),
		Proficient: proficient,
		Expertise:  expert,
	}
}

func buildSkillList(base *templates.TemplateCharacter) models.SkillList {
	return models.SkillList{
		Athletics:      buildSkill(base, "Athletics"),
		Acrobatics:     buildSkill(base, "Acrobatics"),
		SleightOfHand:  buildSkill(base, "SleightOfHand"),
		Stealth:        buildSkill(base, "Stealth"),
		Arcana:         buildSkill(base, "Arcana"),
		History:        buildSkill(base, "History"),
		Investigation:  buildSkill(base, "Investigation"),
		Nature:         buildSkill(base, "Nature"),
		Religion:       buildSkill(base, "Religion"),
		AnimalHandling: buildSkill(base, "AnimalHandling"),
		Insight:        buildSkill(base, "Insight"),
		Medicine:       buildSkill(base, "Medicine"),
		Perception:     buildSkill(base, "Perception"),
		Survival:       buildSkill(base, "Survival"),
		Deception:      buildSkill(base, "Deception"),
		Intimidation:   buildSkill(base, "Intimidation"),
		Performance:    buildSkill(base, "Performance"),
		Persuasion:     buildSkill(base, "Persuasion"),
	}
}
