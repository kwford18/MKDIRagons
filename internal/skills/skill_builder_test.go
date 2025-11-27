package skills_test

import (
	"github.com/kwford18/MKDIRagons/internal/skills"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/template"
	"github.com/stretchr/testify/assert"
)

// setupTestCharacter creates a character instance for testing.
func setupTestCharacter() *template.Character {
	// Creating a character with specific proficiencies for testing
	ab := template.AbilityScores{
		Strength:     10,
		Dexterity:    14,
		Constitution: 10,
		Wisdom:       10,
		Intelligence: 10,
		Charisma:     10,
	}

	char := &template.Character{
		Level:         1,
		AbilityScores: ab,
		Proficiencies: []string{"Stealth", "Acrobatics"},
		Expertise:     []string{"Acrobatics"},
	}

	return char
}

func TestBuildSkill_Logic(t *testing.T) {
	char := setupTestCharacter()

	// Test Base Case (No Proficiency)
	// Assuming Athletics is Strength based. Let's assume Str score 10 (+0 mod).
	// Result should be 0.
	athletics := skills.BuildSkill(char, "Athletics")
	assert.Equal(t, "Athletics", athletics.Name)
	assert.Equal(t, 0, athletics.Bonus)
	assert.False(t, athletics.Proficient)
	assert.False(t, athletics.Expertise)

	// Test Proficiency
	// Stealth (Dex): +2 (Mod) + 2 (Prof) = +4
	stealth := skills.BuildSkill(char, "Stealth")
	assert.Equal(t, "Stealth", stealth.Name)
	assert.Equal(t, 4, stealth.Bonus)
	assert.True(t, stealth.Proficient)
	assert.False(t, stealth.Expertise)

	// Test Expertise
	// Acrobatics (Dex): +2 (Mod) + 2 (Prof) + 2 (Expertise) = +6
	acrobatics := skills.BuildSkill(char, "Acrobatics")
	assert.Equal(t, "Acrobatics", acrobatics.Name)
	assert.Equal(t, 6, acrobatics.Bonus)
	assert.True(t, acrobatics.Proficient)
	assert.True(t, acrobatics.Expertise)
}

func TestBuildSkillList(t *testing.T) {
	char := setupTestCharacter()

	skillList := skills.BuildSkillList(char)

	// Verify the list was populated
	assert.NotEmpty(t, skillList.Athletics.Name)
	assert.NotEmpty(t, skillList.Stealth.Name)

	// Verify mappings logic propagated to the list
	// Check a proficient skill in the list
	assert.True(t, skillList.Stealth.Proficient, "Stealth in list should be proficient")

	// Check a non-proficient skill in the list
	assert.False(t, skillList.Arcana.Proficient, "Arcana in list should not be proficient")

	// Verify Abilities are mapped correctly (assuming D&D 5e defaults in GetSkillAbility)
	// Athletics -> Strength
	// Stealth -> Dexterity
	// Arcana -> Intelligence
	assert.Equal(t, core.Strength, skillList.Athletics.Ability)
	assert.Equal(t, core.Dexterity, skillList.Stealth.Ability)
	assert.Equal(t, core.Intelligence, skillList.Arcana.Ability)
}
