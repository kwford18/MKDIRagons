package skills

import (
	"github.com/kwford18/MKDIRagons/template"
)

func BuildSkill(base *template.Character, name string) Skill {
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
		Athletics:      BuildSkill(base, "Athletics"),
		Acrobatics:     BuildSkill(base, "Acrobatics"),
		SleightOfHand:  BuildSkill(base, "SleightOfHand"),
		Stealth:        BuildSkill(base, "Stealth"),
		Arcana:         BuildSkill(base, "Arcana"),
		History:        BuildSkill(base, "History"),
		Investigation:  BuildSkill(base, "Investigation"),
		Nature:         BuildSkill(base, "Nature"),
		Religion:       BuildSkill(base, "Religion"),
		AnimalHandling: BuildSkill(base, "AnimalHandling"),
		Insight:        BuildSkill(base, "Insight"),
		Medicine:       BuildSkill(base, "Medicine"),
		Perception:     BuildSkill(base, "Perception"),
		Survival:       BuildSkill(base, "Survival"),
		Deception:      BuildSkill(base, "Deception"),
		Intimidation:   BuildSkill(base, "Intimidation"),
		Performance:    BuildSkill(base, "Performance"),
		Persuasion:     BuildSkill(base, "Persuasion"),
	}
}
