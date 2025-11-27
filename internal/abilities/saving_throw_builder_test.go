package abilities_test

import (
	"testing"

	"github.com/kwford18/MKDIRagons/internal/abilities"
	"github.com/kwford18/MKDIRagons/internal/class"
	"github.com/kwford18/MKDIRagons/internal/reference"
	"github.com/kwford18/MKDIRagons/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// BuildSavingThrowsTestSuite defines the test suite for BuildSavingThrows
type BuildSavingThrowsTestSuite struct {
	suite.Suite
	base          *template.Character
	abilityScores abilities.AbilityScores
}

// SetupTest runs before each test
func (suite *BuildSavingThrowsTestSuite) SetupTest() {
	suite.base = &template.Character{
		Level: 1,
	}
	suite.abilityScores = abilities.AbilityScores{
		Strength:     16, // +3 modifier
		Dexterity:    14, // +2 modifier
		Constitution: 15, // +2 modifier
		Intelligence: 12, // +1 modifier
		Wisdom:       10, // +0 modifier
		Charisma:     8,  // -1 modifier
	}
}

// TestBuildSavingThrowsNoProficiencies tests with no proficient saves
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsNoProficiencies() {
	charClass := &class.Class{
		SavingThrows: []reference.Reference{},
	}

	saves := abilities.BuildSavingThrows(suite.base, suite.abilityScores, charClass)

	// Should just be modifiers with no proficiency bonus
	assert.Equal(suite.T(), 3, saves.Strength)     // +3 modifier
	assert.Equal(suite.T(), 2, saves.Dexterity)    // +2 modifier
	assert.Equal(suite.T(), 2, saves.Constitution) // +2 modifier
	assert.Equal(suite.T(), 1, saves.Intelligence) // +1 modifier
	assert.Equal(suite.T(), 0, saves.Wisdom)       // +0 modifier
	assert.Equal(suite.T(), -1, saves.Charisma)    // -1 modifier
}

// TestBuildSavingThrowsStrengthProficiency tests STR save proficiency
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsStrengthProficiency() {
	charClass := &class.Class{
		SavingThrows: []reference.Reference{
			{Name: "STR"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, suite.abilityScores, charClass)

	// STR gets modifier + proficiency bonus (+2 at level 1)
	assert.Equal(suite.T(), 5, saves.Strength)  // +3 modifier + 2 proficiency
	assert.Equal(suite.T(), 2, saves.Dexterity) // No proficiency
}

// TestBuildSavingThrowsDexterityProficiency tests DEX save proficiency
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsDexterityProficiency() {
	charClass := &class.Class{
		SavingThrows: []reference.Reference{
			{Name: "DEX"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, suite.abilityScores, charClass)

	assert.Equal(suite.T(), 3, saves.Strength)
	assert.Equal(suite.T(), 4, saves.Dexterity) // +2 modifier + 2 proficiency
}

// TestBuildSavingThrowsConstitutionProficiency tests CON save proficiency
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsConstitutionProficiency() {
	charClass := &class.Class{
		SavingThrows: []reference.Reference{
			{Name: "CON"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, suite.abilityScores, charClass)

	assert.Equal(suite.T(), 4, saves.Constitution) // +2 modifier + 2 proficiency
}

// TestBuildSavingThrowsIntelligenceProficiency tests INT save proficiency
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsIntelligenceProficiency() {
	charClass := &class.Class{
		SavingThrows: []reference.Reference{
			{Name: "INT"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, suite.abilityScores, charClass)

	assert.Equal(suite.T(), 3, saves.Intelligence) // +1 modifier + 2 proficiency
}

// TestBuildSavingThrowsWisdomProficiency tests WIS save proficiency
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsWisdomProficiency() {
	charClass := &class.Class{
		SavingThrows: []reference.Reference{
			{Name: "WIS"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, suite.abilityScores, charClass)

	assert.Equal(suite.T(), 2, saves.Wisdom) // +0 modifier + 2 proficiency
}

// TestBuildSavingThrowsCharismaProficiency tests CHA save proficiency
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsCharismaProficiency() {
	charClass := &class.Class{
		SavingThrows: []reference.Reference{
			{Name: "CHA"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, suite.abilityScores, charClass)

	assert.Equal(suite.T(), 1, saves.Charisma) // -1 modifier + 2 proficiency
}

// TestBuildSavingThrowsMultipleProficiencies tests multiple save proficiencies
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsMultipleProficiencies() {
	charClass := &class.Class{
		SavingThrows: []reference.Reference{
			{Name: "STR"},
			{Name: "CON"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, suite.abilityScores, charClass)

	assert.Equal(suite.T(), 5, saves.Strength)     // +3 + 2
	assert.Equal(suite.T(), 2, saves.Dexterity)    // +2 (no proficiency)
	assert.Equal(suite.T(), 4, saves.Constitution) // +2 + 2
	assert.Equal(suite.T(), 1, saves.Intelligence) // +1 (no proficiency)
	assert.Equal(suite.T(), 0, saves.Wisdom)       // +0 (no proficiency)
	assert.Equal(suite.T(), -1, saves.Charisma)    // -1 (no proficiency)
}

// TestBuildSavingThrowsFighterClass tests Fighter (STR, CON)
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsFighterClass() {
	charClass := &class.Class{
		Name: "Fighter",
		SavingThrows: []reference.Reference{
			{Name: "STR"},
			{Name: "CON"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, suite.abilityScores, charClass)

	assert.Equal(suite.T(), 5, saves.Strength)
	assert.Equal(suite.T(), 4, saves.Constitution)
	assert.Equal(suite.T(), 2, saves.Dexterity) // No proficiency
}

// TestBuildSavingThrowsWizardClass tests Wizard (INT, WIS)
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsWizardClass() {
	charClass := &class.Class{
		Name: "Wizard",
		SavingThrows: []reference.Reference{
			{Name: "INT"},
			{Name: "WIS"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, suite.abilityScores, charClass)

	assert.Equal(suite.T(), 3, saves.Intelligence)
	assert.Equal(suite.T(), 2, saves.Wisdom)
	assert.Equal(suite.T(), 3, saves.Strength) // No proficiency
}

// TestBuildSavingThrowsRogueClass tests Rogue (DEX, INT)
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsRogueClass() {
	charClass := &class.Class{
		Name: "Rogue",
		SavingThrows: []reference.Reference{
			{Name: "DEX"},
			{Name: "INT"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, suite.abilityScores, charClass)

	assert.Equal(suite.T(), 4, saves.Dexterity)
	assert.Equal(suite.T(), 3, saves.Intelligence)
}

// TestBuildSavingThrowsHigherLevel tests proficiency bonus at higher levels
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsHigherLevel() {
	suite.base.Level = 5 // Proficiency bonus +3

	charClass := &class.Class{
		SavingThrows: []reference.Reference{
			{Name: "WIS"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, suite.abilityScores, charClass)

	assert.Equal(suite.T(), 3, saves.Wisdom) // +0 modifier + 3 proficiency
}

// TestBuildSavingThrowsLevel9 tests proficiency bonus at level 9
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsLevel9() {
	suite.base.Level = 9 // Proficiency bonus +4

	charClass := &class.Class{
		SavingThrows: []reference.Reference{
			{Name: "STR"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, suite.abilityScores, charClass)

	assert.Equal(suite.T(), 7, saves.Strength) // +3 modifier + 4 proficiency
}

// TestBuildSavingThrowsLevel17 tests proficiency bonus at level 17
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsLevel17() {
	suite.base.Level = 17 // Proficiency bonus +6

	charClass := &class.Class{
		SavingThrows: []reference.Reference{
			{Name: "CHA"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, suite.abilityScores, charClass)

	assert.Equal(suite.T(), 5, saves.Charisma) // -1 modifier + 6 proficiency
}

// TestBuildSavingThrowsInvalidSaveName tests unknown save name
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsInvalidSaveName() {
	charClass := &class.Class{
		SavingThrows: []reference.Reference{
			{Name: "INVALID"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, suite.abilityScores, charClass)

	// Invalid name should be ignored, should just be modifiers
	assert.Equal(suite.T(), 3, saves.Strength)
	assert.Equal(suite.T(), 2, saves.Dexterity)
}

// TestBuildSavingThrowsNegativeModifiers tests with negative modifiers
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsNegativeModifiers() {
	lowAbilities := abilities.AbilityScores{
		Strength:     8,  // -1
		Dexterity:    6,  // -2
		Constitution: 8,  // -1
		Intelligence: 10, // +0
		Wisdom:       10, // +0
		Charisma:     8,  // -1
	}

	charClass := &class.Class{
		SavingThrows: []reference.Reference{
			{Name: "DEX"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, lowAbilities, charClass)

	assert.Equal(suite.T(), -1, saves.Strength) // -1 modifier, no proficiency
	assert.Equal(suite.T(), 0, saves.Dexterity) // -2 modifier + 2 proficiency
	assert.Equal(suite.T(), -1, saves.Constitution)
}

// TestBuildSavingThrowsMaxAbilities tests with maxed abilities
func (suite *BuildSavingThrowsTestSuite) TestBuildSavingThrowsMaxAbilities() {
	maxAbilities := abilities.AbilityScores{
		Strength:     20, // +5
		Dexterity:    20, // +5
		Constitution: 20, // +5
		Intelligence: 20, // +5
		Wisdom:       20, // +5
		Charisma:     20, // +5
	}

	charClass := &class.Class{
		SavingThrows: []reference.Reference{
			{Name: "STR"},
			{Name: "DEX"},
		},
	}

	saves := abilities.BuildSavingThrows(suite.base, maxAbilities, charClass)

	assert.Equal(suite.T(), 7, saves.Strength)  // +5 + 2
	assert.Equal(suite.T(), 7, saves.Dexterity) // +5 + 2
	assert.Equal(suite.T(), 5, saves.Constitution)
}

func TestBuildSavingThrowsTestSuite(t *testing.T) {
	suite.Run(t, new(BuildSavingThrowsTestSuite))
}

// TestBuildSavingThrowsTableDriven uses table-driven tests
func TestBuildSavingThrowsTableDriven(t *testing.T) {
	testCases := []struct {
		name          string
		level         int
		strScore      int
		proficiencies []string
		expectedSTR   int
	}{
		{
			name:  "Level 1, STR 16, proficient",
			level: 1, strScore: 16,
			proficiencies: []string{"STR"},
			expectedSTR:   5, // +3 modifier + 2 proficiency
		},
		{
			name:  "Level 1, STR 16, not proficient",
			level: 1, strScore: 16,
			proficiencies: []string{},
			expectedSTR:   3, // +3 modifier only
		},
		{
			name:  "Level 5, STR 18, proficient",
			level: 5, strScore: 18,
			proficiencies: []string{"STR"},
			expectedSTR:   7, // +4 modifier + 3 proficiency
		},
		{
			name:  "Level 17, STR 20, proficient",
			level: 17, strScore: 20,
			proficiencies: []string{"STR"},
			expectedSTR:   11, // +5 modifier + 6 proficiency
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			base := &template.Character{Level: tc.level}
			abilityScores := abilities.AbilityScores{Strength: tc.strScore}

			var savingThrows []reference.Reference
			for _, prof := range tc.proficiencies {
				savingThrows = append(savingThrows, reference.Reference{Name: prof})
			}
			charClass := &class.Class{SavingThrows: savingThrows}

			saves := abilities.BuildSavingThrows(base, abilityScores, charClass)

			assert.Equal(t, tc.expectedSTR, saves.Strength)
		})
	}
}

// BenchmarkBuildSavingThrowsNoProficiency benchmarks with no proficiencies
func BenchmarkBuildSavingThrowsNoProficiency(b *testing.B) {
	base := &template.Character{Level: 1}
	abilityScores := abilities.AbilityScores{
		Strength: 16, Dexterity: 14, Constitution: 15,
		Intelligence: 12, Wisdom: 10, Charisma: 8,
	}
	charClass := &class.Class{SavingThrows: []reference.Reference{}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = abilities.BuildSavingThrows(base, abilityScores, charClass)
	}
}

// BenchmarkBuildSavingThrowsWithProficiencies benchmarks with proficiencies
func BenchmarkBuildSavingThrowsWithProficiencies(b *testing.B) {
	base := &template.Character{Level: 1}
	abilityScores := abilities.AbilityScores{
		Strength: 16, Dexterity: 14, Constitution: 15,
		Intelligence: 12, Wisdom: 10, Charisma: 8,
	}
	charClass := &class.Class{
		SavingThrows: []reference.Reference{
			{Name: "STR"},
			{Name: "CON"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = abilities.BuildSavingThrows(base, abilityScores, charClass)
	}
}

// ExampleBuildSavingThrows demonstrates basic usage
func ExampleBuildSavingThrows() {
	base := &template.Character{Level: 1}
	abilityScores := abilities.AbilityScores{
		Strength:     16,
		Dexterity:    14,
		Constitution: 15,
		Intelligence: 12,
		Wisdom:       10,
		Charisma:     8,
	}
	fighter := &class.Class{
		SavingThrows: []reference.Reference{
			{Name: "STR"},
			{Name: "CON"},
		},
	}

	saves := abilities.BuildSavingThrows(base, abilityScores, fighter)
	_ = saves.Strength     // 5 (+3 modifier + 2 proficiency)
	_ = saves.Constitution // 4 (+2 modifier + 2 proficiency)
	_ = saves.Dexterity    // 2 (+2 modifier, no proficiency)
}
