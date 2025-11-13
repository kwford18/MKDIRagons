package template_test

import (
	"fmt"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TemplateCharacterTestSuite defines the test suite for TemplateCharacter
type TemplateCharacterTestSuite struct {
	suite.Suite
	character *template.Character
}

// SetupTest runs before each test
func (suite *TemplateCharacterTestSuite) SetupTest() {
	suite.character = &template.Character{
		Name:     "Gandalf",
		Level:    5,
		Race:     "human",
		Subrace:  "",
		Class:    "wizard",
		Subclass: "evocation",
		AbilityScores: template.AbilityScores{
			Strength:     10,
			Dexterity:    14,
			Constitution: 12,
			Intelligence: 18,
			Wisdom:       15,
			Charisma:     8,
		},
		Proficiencies: []string{"arcana", "history"},
		Expertise:     []string{},
		Inventory: template.Inventory{
			Weapons: []string{"dagger", "quarterstaff"},
			Armor:   []string{},
			Items:   []string{"spellbook", "component-pouch"},
		},
		Spells: template.Spells{
			Level: [][]string{
				{"fire-bolt", "mage-hand"},
				{"magic-missile", "shield"},
			},
		},
	}
}

// TestTemplateCharacterFields tests basic field assignment
func (suite *TemplateCharacterTestSuite) TestTemplateCharacterFields() {
	assert.Equal(suite.T(), "Gandalf", suite.character.Name)
	assert.Equal(suite.T(), 5, suite.character.Level)
	assert.Equal(suite.T(), "human", suite.character.Race)
	assert.Equal(suite.T(), "wizard", suite.character.Class)
	assert.Equal(suite.T(), "evocation", suite.character.Subclass)
}

// TestTemplateCharacterAbilityScores tests ability scores
func (suite *TemplateCharacterTestSuite) TestTemplateCharacterAbilityScores() {
	assert.Equal(suite.T(), 10, suite.character.AbilityScores.Strength)
	assert.Equal(suite.T(), 14, suite.character.AbilityScores.Dexterity)
	assert.Equal(suite.T(), 12, suite.character.AbilityScores.Constitution)
	assert.Equal(suite.T(), 18, suite.character.AbilityScores.Intelligence)
	assert.Equal(suite.T(), 15, suite.character.AbilityScores.Wisdom)
	assert.Equal(suite.T(), 8, suite.character.AbilityScores.Charisma)
}

// TestTemplateCharacterPrint tests Print doesn't panic
func (suite *TemplateCharacterTestSuite) TestTemplateCharacterPrint() {
	assert.NotPanics(suite.T(), func() {
		suite.character.Print()
	})
}

func TestTemplateCharacterTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateCharacterTestSuite))
}

// ========== ProficiencyBonus Tests ==========

// TestProficiencyBonusLevels tests proficiency bonus for all levels
func TestProficiencyBonusLevels(t *testing.T) {
	testCases := []struct {
		level    int
		expected int
	}{
		{1, 2}, {2, 2}, {3, 2}, {4, 2},
		{5, 3}, {6, 3}, {7, 3}, {8, 3},
		{9, 4}, {10, 4}, {11, 4}, {12, 4},
		{13, 5}, {14, 5}, {15, 5}, {16, 5},
		{17, 6}, {18, 6}, {19, 6}, {20, 6},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Level %d", tc.level), func(t *testing.T) {
			char := &template.Character{Level: tc.level}
			assert.Equal(t, tc.expected, char.ProficiencyBonus())
		})
	}
}

// TestProficiencyBonusInvalidLevel tests invalid levels
func TestProficiencyBonusInvalidLevel(t *testing.T) {
	testCases := []int{0, -1, 21, 100}

	for _, level := range testCases {
		t.Run(fmt.Sprintf("Level %d", level), func(t *testing.T) {
			char := &template.Character{Level: level}
			assert.Equal(t, -1, char.ProficiencyBonus())
		})
	}
}

// ========== GetSkillAbility Tests ==========

// TestGetSkillAbilityStrength tests Strength-based skills
func TestGetSkillAbilityStrength(t *testing.T) {
	char := &template.Character{}
	assert.Equal(t, core.Strength, char.GetSkillAbility("Athletics"))
}

// TestGetSkillAbilityDexterity tests Dexterity-based skills
func TestGetSkillAbilityDexterity(t *testing.T) {
	char := &template.Character{}

	testCases := []string{"Acrobatics", "SleightOfHand", "Stealth"}
	for _, skill := range testCases {
		t.Run(skill, func(t *testing.T) {
			assert.Equal(t, core.Dexterity, char.GetSkillAbility(skill))
		})
	}
}

// TestGetSkillAbilityIntelligence tests Intelligence-based skills
func TestGetSkillAbilityIntelligence(t *testing.T) {
	char := &template.Character{}

	testCases := []string{"Arcana", "History", "Investigation", "Nature", "Religion"}
	for _, skill := range testCases {
		t.Run(skill, func(t *testing.T) {
			assert.Equal(t, core.Intelligence, char.GetSkillAbility(skill))
		})
	}
}

// TestGetSkillAbilityWisdom tests Wisdom-based skills
func TestGetSkillAbilityWisdom(t *testing.T) {
	char := &template.Character{}

	testCases := []string{"AnimalHandling", "Insight", "Medicine", "Perception", "Survival"}
	for _, skill := range testCases {
		t.Run(skill, func(t *testing.T) {
			assert.Equal(t, core.Wisdom, char.GetSkillAbility(skill))
		})
	}
}

// TestGetSkillAbilityCharisma tests Charisma-based skills
func TestGetSkillAbilityCharisma(t *testing.T) {
	char := &template.Character{}

	testCases := []string{"Deception", "Intimidation", "Performance", "Persuasion"}
	for _, skill := range testCases {
		t.Run(skill, func(t *testing.T) {
			assert.Equal(t, core.Charisma, char.GetSkillAbility(skill))
		})
	}
}

// TestGetSkillAbilityInvalidSkill tests invalid skill name
func TestGetSkillAbilityInvalidSkill(t *testing.T) {
	char := &template.Character{}

	// Returns 0 (core.Strength) for invalid skills
	result := char.GetSkillAbility("InvalidSkill")
	assert.Equal(t, core.Ability(0), result)
	assert.Equal(t, core.Strength, result) // 0 is Strength
}

// TestGetSkillAbilityAllSkills tests all valid skills
func TestGetSkillAbilityAllSkills(t *testing.T) {
	char := &template.Character{}

	skillMap := map[string]core.Ability{
		// Strength
		"Athletics": core.Strength,
		// Dexterity
		"Acrobatics": core.Dexterity, "SleightOfHand": core.Dexterity, "Stealth": core.Dexterity,
		// Intelligence
		"Arcana": core.Intelligence, "History": core.Intelligence, "Investigation": core.Intelligence,
		"Nature": core.Intelligence, "Religion": core.Intelligence,
		// Wisdom
		"AnimalHandling": core.Wisdom, "Insight": core.Wisdom, "Medicine": core.Wisdom,
		"Perception": core.Wisdom, "Survival": core.Wisdom,
		// Charisma
		"Deception": core.Charisma, "Intimidation": core.Charisma,
		"Performance": core.Charisma, "Persuasion": core.Charisma,
	}

	for skill, expectedAbility := range skillMap {
		t.Run(skill, func(t *testing.T) {
			assert.Equal(t, expectedAbility, char.GetSkillAbility(skill))
		})
	}
}

// ========== TemplateAbilityScores Modifier Tests ==========

// TestTemplateAbilityScoresModifier tests modifier calculation
func TestTemplateAbilityScoresModifier(t *testing.T) {
	scores := template.AbilityScores{
		Strength:     16,
		Dexterity:    14,
		Constitution: 15,
		Intelligence: 12,
		Wisdom:       10,
		Charisma:     8,
	}

	testCases := []struct {
		ability  core.Ability
		expected int
	}{
		{core.Strength, 3},     // (16-10)/2 = 3
		{core.Dexterity, 2},    // (14-10)/2 = 2
		{core.Constitution, 2}, // (15-10)/2 = 2
		{core.Intelligence, 1}, // (12-10)/2 = 1
		{core.Wisdom, 0},       // (10-10)/2 = 0
		{core.Charisma, -1},    // (8-10)/2 = -1
	}

	for _, tc := range testCases {
		t.Run(tc.ability.String(), func(t *testing.T) {
			assert.Equal(t, tc.expected, scores.Modifier(tc.ability))
		})
	}
}

// TestTemplateAbilityScoresModifierEdgeCases tests edge cases
func TestTemplateAbilityScoresModifierEdgeCases(t *testing.T) {
	testCases := []struct {
		score    int
		expected int
	}{
		{1, -5},
		{3, -4},
		{8, -1},
		{9, -1},
		{10, 0},
		{11, 0},
		{12, 1},
		{18, 4},
		{20, 5},
		{30, 10},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Score %d", tc.score), func(t *testing.T) {
			scores := template.AbilityScores{Strength: tc.score}
			assert.Equal(t, tc.expected, scores.Modifier(core.Strength))
		})
	}
}

// ========== TemplateInventory Tests ==========

// TestTemplateInventory tests inventory structure
func TestTemplateInventory(t *testing.T) {
	inv := template.Inventory{
		Weapons: []string{"longsword", "shortbow"},
		Armor:   []string{"leather-armor", "shield"},
		Items:   []string{"rope", "torch"},
	}

	assert.Len(t, inv.Weapons, 2)
	assert.Len(t, inv.Armor, 2)
	assert.Len(t, inv.Items, 2)
	assert.Equal(t, "longsword", inv.Weapons[0])
	assert.Equal(t, "leather-armor", inv.Armor[0])
	assert.Equal(t, "rope", inv.Items[0])
}

// TestEmptyTemplateInventory tests empty inventory
func TestEmptyTemplateInventory(t *testing.T) {
	inv := template.Inventory{}

	assert.Nil(t, inv.Weapons)
	assert.Nil(t, inv.Armor)
	assert.Nil(t, inv.Items)
}

// ========== TemplateSpells Tests ==========

// TestTemplateSpells tests spell structure
func TestTemplateSpells(t *testing.T) {
	spells := template.Spells{
		Level: [][]string{
			{"fire-bolt", "mage-hand"},
			{"magic-missile", "shield"},
			{"misty-step"},
		},
	}

	assert.Len(t, spells.Level, 3)
	assert.Len(t, spells.Level[0], 2) // Cantrips
	assert.Len(t, spells.Level[1], 2) // Level 1
	assert.Len(t, spells.Level[2], 1) // Level 2
	assert.Equal(t, "fire-bolt", spells.Level[0][0])
	assert.Equal(t, "magic-missile", spells.Level[1][0])
}

// TestEmptyTemplateSpells tests empty spells
func TestEmptyTemplateSpells(t *testing.T) {
	spells := template.Spells{}
	assert.Nil(t, spells.Level)
}

// ========== Integration Tests ==========

// TestCompleteTemplateCharacter tests a fully populated character
func TestCompleteTemplateCharacter(t *testing.T) {
	char := &template.Character{
		Name:     "Aragorn",
		Level:    10,
		Race:     "human",
		Subrace:  "",
		Class:    "ranger",
		Subclass: "hunter",
		AbilityScores: template.AbilityScores{
			Strength:     18,
			Dexterity:    16,
			Constitution: 14,
			Intelligence: 10,
			Wisdom:       15,
			Charisma:     12,
		},
		Proficiencies: []string{"athletics", "survival", "nature"},
		Expertise:     []string{"survival"},
		Inventory: template.Inventory{
			Weapons: []string{"longsword", "longbow"},
			Armor:   []string{"studded-leather"},
			Items:   []string{"rope", "bedroll"},
		},
		Spells: template.Spells{
			Level: [][]string{
				{"hunters-mark"},
				{"pass-without-trace"},
			},
		},
	}

	assert.Equal(t, "Aragorn", char.Name)
	assert.Equal(t, 10, char.Level)
	assert.Equal(t, 4, char.ProficiencyBonus())
	assert.Equal(t, 18, char.AbilityScores.Strength)
	assert.Len(t, char.Proficiencies, 3)
	assert.Len(t, char.Expertise, 1)
	assert.Len(t, char.Inventory.Weapons, 2)
}

// TestMinimalTemplateCharacter tests minimal character
func TestMinimalTemplateCharacter(t *testing.T) {
	char := &template.Character{
		Name:  "Bob",
		Level: 1,
		Race:  "human",
		Class: "fighter",
	}

	assert.Equal(t, "Bob", char.Name)
	assert.Equal(t, 1, char.Level)
	assert.Equal(t, 2, char.ProficiencyBonus())
	assert.Empty(t, char.Subrace)
	assert.Empty(t, char.Subclass)
}

// ========== Benchmark Tests ==========

// BenchmarkProficiencyBonus benchmarks proficiency bonus calculation
func BenchmarkProficiencyBonus(b *testing.B) {
	char := &template.Character{Level: 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = char.ProficiencyBonus()
	}
}

// BenchmarkGetSkillAbility benchmarks skill ability lookup
func BenchmarkGetSkillAbility(b *testing.B) {
	char := &template.Character{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = char.GetSkillAbility("Athletics")
	}
}

// BenchmarkTemplateAbilityScoresModifier benchmarks modifier calculation
func BenchmarkTemplateAbilityScoresModifier(b *testing.B) {
	scores := template.AbilityScores{Strength: 16}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = scores.Modifier(core.Strength)
	}
}

// ========== Example Tests ==========

// ExampleTemplateCharacter_ProficiencyBonus demonstrates proficiency bonus
func ExampleTemplateCharacter_ProficiencyBonus() {
	char := &template.Character{Level: 5}
	bonus := char.ProficiencyBonus()
	_ = bonus // 3
}

// ExampleTemplateCharacter_GetSkillAbility demonstrates skill ability lookup
func ExampleTemplateCharacter_GetSkillAbility() {
	char := &template.Character{}
	ability := char.GetSkillAbility("Athletics")
	_ = ability // core.Strength
}

// ExampleTemplateAbilityScores_Modifier demonstrates modifier calculation
func ExampleTemplateAbilityScores_Modifier() {
	scores := template.AbilityScores{
		Strength: 16,
	}
	modifier := scores.Modifier(core.Strength)
	_ = modifier // 3
}
