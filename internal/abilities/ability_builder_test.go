package abilities_test

import (
	"testing"

	"github.com/kwford18/MKDIRagons/internal/abilities"
	"github.com/kwford18/MKDIRagons/internal/race"
	"github.com/kwford18/MKDIRagons/internal/reference"
	"github.com/kwford18/MKDIRagons/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// BuildAbilityScoresTestSuite defines the test suite for BuildAbilityScores
type BuildAbilityScoresTestSuite struct {
	suite.Suite
	base *template.Character
}

// SetupTest runs before each test
func (suite *BuildAbilityScoresTestSuite) SetupTest() {
	suite.base = &template.Character{
		AbilityScores: template.AbilityScores{
			Strength:     10,
			Dexterity:    12,
			Constitution: 14,
			Intelligence: 13,
			Wisdom:       15,
			Charisma:     8,
		},
	}
}

// TestBuildAbilityScoresNoRacialBonus tests building with no racial bonuses
func (suite *BuildAbilityScoresTestSuite) TestBuildAbilityScoresNoRacialBonus() {
	playerRace := race.Race{
		AbilityBonuses: []race.AbilityBonus{},
	}

	result := abilities.BuildAbilityScores(suite.base, playerRace)

	assert.Equal(suite.T(), 10, result.Strength)
	assert.Equal(suite.T(), 12, result.Dexterity)
	assert.Equal(suite.T(), 14, result.Constitution)
	assert.Equal(suite.T(), 13, result.Intelligence)
	assert.Equal(suite.T(), 15, result.Wisdom)
	assert.Equal(suite.T(), 8, result.Charisma)
}

// TestBuildAbilityScoresWithStrengthBonus tests STR racial bonus
func (suite *BuildAbilityScoresTestSuite) TestBuildAbilityScoresWithStrengthBonus() {
	playerRace := race.Race{
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{Name: "STR"},
				Bonus:        2,
			},
		},
	}

	result := abilities.BuildAbilityScores(suite.base, playerRace)

	assert.Equal(suite.T(), 12, result.Strength) // 10 + 2
	assert.Equal(suite.T(), 12, result.Dexterity)
	assert.Equal(suite.T(), 14, result.Constitution)
}

// TestBuildAbilityScoresWithDexterityBonus tests DEX racial bonus
func (suite *BuildAbilityScoresTestSuite) TestBuildAbilityScoresWithDexterityBonus() {
	playerRace := race.Race{
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{Name: "DEX"},
				Bonus:        2,
			},
		},
	}

	result := abilities.BuildAbilityScores(suite.base, playerRace)

	assert.Equal(suite.T(), 10, result.Strength)
	assert.Equal(suite.T(), 14, result.Dexterity) // 12 + 2
}

// TestBuildAbilityScoresWithConstitutionBonus tests CON racial bonus
func (suite *BuildAbilityScoresTestSuite) TestBuildAbilityScoresWithConstitutionBonus() {
	playerRace := race.Race{
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{Name: "CON"},
				Bonus:        2,
			},
		},
	}

	result := abilities.BuildAbilityScores(suite.base, playerRace)

	assert.Equal(suite.T(), 16, result.Constitution) // 14 + 2
}

// TestBuildAbilityScoresWithIntelligenceBonus tests INT racial bonus
func (suite *BuildAbilityScoresTestSuite) TestBuildAbilityScoresWithIntelligenceBonus() {
	playerRace := race.Race{
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{Name: "INT"},
				Bonus:        2,
			},
		},
	}

	result := abilities.BuildAbilityScores(suite.base, playerRace)

	assert.Equal(suite.T(), 15, result.Intelligence) // 13 + 2
}

// TestBuildAbilityScoresWithWisdomBonus tests WIS racial bonus
func (suite *BuildAbilityScoresTestSuite) TestBuildAbilityScoresWithWisdomBonus() {
	playerRace := race.Race{
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{Name: "WIS"},
				Bonus:        2,
			},
		},
	}

	result := abilities.BuildAbilityScores(suite.base, playerRace)

	assert.Equal(suite.T(), 17, result.Wisdom) // 15 + 2
}

// TestBuildAbilityScoresWithCharismaBonus tests CHA racial bonus
func (suite *BuildAbilityScoresTestSuite) TestBuildAbilityScoresWithCharismaBonus() {
	playerRace := race.Race{
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{Name: "CHA"},
				Bonus:        2,
			},
		},
	}

	result := abilities.BuildAbilityScores(suite.base, playerRace)

	assert.Equal(suite.T(), 10, result.Charisma) // 8 + 2
}

// TestBuildAbilityScoresMultipleBonuses tests multiple racial bonuses (like Half-Elf)
func (suite *BuildAbilityScoresTestSuite) TestBuildAbilityScoresMultipleBonuses() {
	playerRace := race.Race{
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{Name: "CHA"},
				Bonus:        2,
			},
			{
				AbilityScore: reference.Reference{Name: "STR"},
				Bonus:        1,
			},
			{
				AbilityScore: reference.Reference{Name: "DEX"},
				Bonus:        1,
			},
		},
	}

	result := abilities.BuildAbilityScores(suite.base, playerRace)

	assert.Equal(suite.T(), 11, result.Strength)  // 10 + 1
	assert.Equal(suite.T(), 13, result.Dexterity) // 12 + 1
	assert.Equal(suite.T(), 14, result.Constitution)
	assert.Equal(suite.T(), 13, result.Intelligence)
	assert.Equal(suite.T(), 15, result.Wisdom)
	assert.Equal(suite.T(), 10, result.Charisma) // 8 + 2
}

// TestBuildAbilityScoresElfRace tests typical Elf racial bonus (+2 DEX)
func (suite *BuildAbilityScoresTestSuite) TestBuildAbilityScoresElfRace() {
	playerRace := race.Race{
		Name: "Elf",
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{Name: "DEX"},
				Bonus:        2,
			},
		},
	}

	result := abilities.BuildAbilityScores(suite.base, playerRace)

	assert.Equal(suite.T(), 14, result.Dexterity) // 12 + 2
}

// TestBuildAbilityScoresDwarfRace tests typical Dwarf racial bonus (+2 CON)
func (suite *BuildAbilityScoresTestSuite) TestBuildAbilityScoresDwarfRace() {
	playerRace := race.Race{
		Name: "Dwarf",
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{Name: "CON"},
				Bonus:        2,
			},
		},
	}

	result := abilities.BuildAbilityScores(suite.base, playerRace)

	assert.Equal(suite.T(), 16, result.Constitution) // 14 + 2
}

// TestBuildAbilityScoresHumanRace tests Human racial bonus (+1 to all)
func (suite *BuildAbilityScoresTestSuite) TestBuildAbilityScoresHumanRace() {
	playerRace := race.Race{
		Name: "Human",
		AbilityBonuses: []race.AbilityBonus{
			{AbilityScore: reference.Reference{Name: "STR"}, Bonus: 1},
			{AbilityScore: reference.Reference{Name: "DEX"}, Bonus: 1},
			{AbilityScore: reference.Reference{Name: "CON"}, Bonus: 1},
			{AbilityScore: reference.Reference{Name: "INT"}, Bonus: 1},
			{AbilityScore: reference.Reference{Name: "WIS"}, Bonus: 1},
			{AbilityScore: reference.Reference{Name: "CHA"}, Bonus: 1},
		},
	}

	result := abilities.BuildAbilityScores(suite.base, playerRace)

	assert.Equal(suite.T(), 11, result.Strength)     // 10 + 1
	assert.Equal(suite.T(), 13, result.Dexterity)    // 12 + 1
	assert.Equal(suite.T(), 15, result.Constitution) // 14 + 1
	assert.Equal(suite.T(), 14, result.Intelligence) // 13 + 1
	assert.Equal(suite.T(), 16, result.Wisdom)       // 15 + 1
	assert.Equal(suite.T(), 9, result.Charisma)      // 8 + 1
}

// TestBuildAbilityScoresInvalidAbilityName tests unknown ability name
func (suite *BuildAbilityScoresTestSuite) TestBuildAbilityScoresInvalidAbilityName() {
	playerRace := race.Race{
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{Name: "INVALID"},
				Bonus:        2,
			},
			{
				AbilityScore: reference.Reference{Name: "STR"},
				Bonus:        1,
			},
		},
	}

	result := abilities.BuildAbilityScores(suite.base, playerRace)

	// Should stop processing at invalid ability (break OuterLoop)
	// Original scores should remain
	assert.Equal(suite.T(), 10, result.Strength)
	assert.Equal(suite.T(), 12, result.Dexterity)
}

// TestBuildAbilityScoresNegativeBonus tests negative racial bonus (unusual)
func (suite *BuildAbilityScoresTestSuite) TestBuildAbilityScoresNegativeBonus() {
	playerRace := race.Race{
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{Name: "STR"},
				Bonus:        -2,
			},
		},
	}

	result := abilities.BuildAbilityScores(suite.base, playerRace)

	assert.Equal(suite.T(), 8, result.Strength) // 10 - 2
}

// TestBuildAbilityScoresLargeBonuses tests large bonuses
func (suite *BuildAbilityScoresTestSuite) TestBuildAbilityScoresLargeBonuses() {
	playerRace := race.Race{
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{Name: "STR"},
				Bonus:        5,
			},
		},
	}

	result := abilities.BuildAbilityScores(suite.base, playerRace)

	assert.Equal(suite.T(), 15, result.Strength) // 10 + 5
}

func TestBuildAbilityScoresTestSuite(t *testing.T) {
	suite.Run(t, new(BuildAbilityScoresTestSuite))
}

// TestBuildAbilityScoresTableDriven uses table-driven tests
func TestBuildAbilityScoresTableDriven(t *testing.T) {
	testCases := []struct {
		name        string
		baseScores  template.AbilityScores
		racialBonus []race.AbilityBonus
		expectedSTR int
		expectedDEX int
		expectedCON int
	}{
		{
			name: "No bonuses",
			baseScores: template.AbilityScores{
				Strength: 10, Dexterity: 10, Constitution: 10,
			},
			racialBonus: []race.AbilityBonus{},
			expectedSTR: 10, expectedDEX: 10, expectedCON: 10,
		},
		{
			name: "STR +2",
			baseScores: template.AbilityScores{
				Strength: 15, Dexterity: 10, Constitution: 10,
			},
			racialBonus: []race.AbilityBonus{
				{AbilityScore: reference.Reference{Name: "STR"}, Bonus: 2},
			},
			expectedSTR: 17, expectedDEX: 10, expectedCON: 10,
		},
		{
			name: "Multiple bonuses",
			baseScores: template.AbilityScores{
				Strength: 14, Dexterity: 12, Constitution: 13,
			},
			racialBonus: []race.AbilityBonus{
				{AbilityScore: reference.Reference{Name: "STR"}, Bonus: 2},
				{AbilityScore: reference.Reference{Name: "CON"}, Bonus: 1},
			},
			expectedSTR: 16, expectedDEX: 12, expectedCON: 14,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			base := &template.Character{
				AbilityScores: tc.baseScores,
			}
			playerRace := race.Race{
				AbilityBonuses: tc.racialBonus,
			}

			result := abilities.BuildAbilityScores(base, playerRace)

			assert.Equal(t, tc.expectedSTR, result.Strength)
			assert.Equal(t, tc.expectedDEX, result.Dexterity)
			assert.Equal(t, tc.expectedCON, result.Constitution)
		})
	}
}

// BenchmarkBuildAbilityScoresNoBonus benchmarks with no racial bonuses
func BenchmarkBuildAbilityScoresNoBonus(b *testing.B) {
	base := &template.Character{
		AbilityScores: template.AbilityScores{
			Strength: 10, Dexterity: 12, Constitution: 14,
			Intelligence: 13, Wisdom: 15, Charisma: 8,
		},
	}
	playerRace := race.Race{AbilityBonuses: []race.AbilityBonus{}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = abilities.BuildAbilityScores(base, playerRace)
	}
}

// BenchmarkBuildAbilityScoresMultipleBonus benchmarks with multiple bonuses
func BenchmarkBuildAbilityScoresMultipleBonus(b *testing.B) {
	base := &template.Character{
		AbilityScores: template.AbilityScores{
			Strength: 10, Dexterity: 12, Constitution: 14,
			Intelligence: 13, Wisdom: 15, Charisma: 8,
		},
	}
	playerRace := race.Race{
		AbilityBonuses: []race.AbilityBonus{
			{AbilityScore: reference.Reference{Name: "CHA"}, Bonus: 2},
			{AbilityScore: reference.Reference{Name: "STR"}, Bonus: 1},
			{AbilityScore: reference.Reference{Name: "DEX"}, Bonus: 1},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = abilities.BuildAbilityScores(base, playerRace)
	}
}

// ExampleBuildAbilityScores demonstrates basic usage
func ExampleBuildAbilityScores() {
	base := &template.Character{
		AbilityScores: template.AbilityScores{
			Strength:     10,
			Dexterity:    14,
			Constitution: 12,
			Intelligence: 13,
			Wisdom:       15,
			Charisma:     8,
		},
	}

	elfRace := race.Race{
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{Name: "DEX"},
				Bonus:        2,
			},
		},
	}

	result := abilities.BuildAbilityScores(base, elfRace)
	_ = result.Dexterity // 16 (14 + 2)
}
