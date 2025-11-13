package skills

import (
	"github.com/kwford18/MKDIRagons/template"
)

func buildSkill(base *template.Character, name string) Skill {
	// Calculate bonus by getting the modifier of input skill
	bonus := base.AbilityScores.Modifier(base.GetSkillAbility(name))

	// proficiencies & expertise are false by default
	proficient := false
	expert := false

	// Check for proficiency & expertise
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

	return Skill{
		Name:       name,
		Bonus:      bonus,
		Ability:    base.GetSkillAbility(name),
		Proficient: proficient,
		Expertise:  expert,
	}
}

func BuildSkillList(base *template.Character) SkillList {
	return SkillList{
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
