package builder

import (
	"fmt"
	"strings"

	"github.com/kwford18/MKDIRagons/templates"
)

func FormatKeys(input string) (string, error) {
	if len(input) < 3 {
		return "undefined", fmt.Errorf("input is not valid: %s", input)
	}
	lower := strings.ToLower(input)
	key := lower[:3]

	return key, nil
}

func buildAbilityScores(base *templates.TemplateCharacter) map[string]int {
	ability_scores := make(map[string]int, 10)

	ability_scores["str"] = base.AbilityScores.Strength
	ability_scores["dex"] = base.AbilityScores.Dexterity
	ability_scores["con"] = base.AbilityScores.Constitution
	ability_scores["wis"] = base.AbilityScores.Wisdom
	ability_scores["int"] = base.AbilityScores.Intelligence
	ability_scores["cha"] = base.AbilityScores.Charisma

	return ability_scores
}

// func buildSkillBonuses(character *templates.Character) map[string]int {

// }
