package character_test

import (
	"fmt"
	"github.com/kwford18/MKDIRagons/internal/character"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestCharacter_ProficiencyBonus(t *testing.T) {
	tests := []struct {
		level    int
		expected int
	}{
		// Tier 1 (Levels 1-4)
		{1, 2},
		{4, 2},
		// Tier 2 (Levels 5-8)
		{5, 3},
		{8, 3},
		// Tier 3 (Levels 9-12)
		{9, 4},
		{12, 4},
		// Tier 4 (Levels 13-16)
		{13, 5},
		{16, 5},
		// Tier 5 (Levels 17-20)
		{17, 6},
		{20, 6},
		// Edge Cases / Invalid Levels
		{0, -1},
		{21, -1},
		{-5, -1},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level %d", tt.level), func(t *testing.T) {
			c := &character.Character{Level: tt.level}
			assert.Equal(t, tt.expected, c.ProficiencyBonus())
		})
	}
}

func TestCharacter_GetSkillAbility(t *testing.T) {
	c := &character.Character{}

	tests := []struct {
		skill    string
		expected core.Ability
	}{
		// Strength
		{"Athletics", core.Strength},

		// Dexterity
		{"Acrobatics", core.Dexterity},
		{"SleightOfHand", core.Dexterity},
		{"Stealth", core.Dexterity},

		// Intelligence
		{"Arcana", core.Intelligence},
		{"History", core.Intelligence},
		{"Investigation", core.Intelligence},
		{"Nature", core.Intelligence},
		{"Religion", core.Intelligence},

		// Wisdom
		{"AnimalHandling", core.Wisdom},
		{"Insight", core.Wisdom},
		{"Medicine", core.Wisdom},
		{"Perception", core.Wisdom},
		{"Survival", core.Wisdom},

		// Charisma
		{"Deception", core.Charisma},
		{"Intimidation", core.Charisma},
		{"Performance", core.Charisma},
		{"Persuasion", core.Charisma},

		// Unknown/Default
		{"NotASkill", core.Ability(0)},
		{"", core.Ability(0)},
	}

	for _, tt := range tests {
		t.Run(tt.skill, func(t *testing.T) {
			result := c.GetSkillAbility(tt.skill)
			assert.Equal(t, tt.expected, result, "Failed for skill: %s", tt.skill)
		})
	}
}
