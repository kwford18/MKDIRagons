package abilities_test

import (
	"testing"

	"github.com/kwford18/MKDIRagons/internal/abilities"
	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AbilityScoreTestSuite defines the test suite for AbilityScore
type AbilityScoreTestSuite struct {
	suite.Suite
	scores *abilities.AbilityScore
}

// SetupTest runs before each test
func (suite *AbilityScoreTestSuite) SetupTest() {
	suite.scores = &abilities.AbilityScore{
		Strength:     16,
		Dexterity:    14,
		Constitution: 15,
		Intelligence: 12,
		Wisdom:       10,
		Charisma:     8,
	}
}

// TestAbilityScoreFields tests basic field assignment
func (suite *AbilityScoreTestSuite) TestAbilityScoreFields() {
	assert.Equal(suite.T(), 16, suite.scores.Strength)
	assert.Equal(suite.T(), 14, suite.scores.Dexterity)
	assert.Equal(suite.T(), 15, suite.scores.Constitution)
	assert.Equal(suite.T(), 12, suite.scores.Intelligence)
	assert.Equal(suite.T(), 10, suite.scores.Wisdom)
	assert.Equal(suite.T(), 8, suite.scores.Charisma)
}

// TestModifierStrength tests Strength modifier calculation
func (suite *AbilityScoreTestSuite) TestModifierStrength() {
	modifier := suite.scores.Modifier(core.Strength)
	assert.Equal(suite.T(), 3, modifier) // (16-10)/2 = 3
}

// TestModifierDexterity tests Dexterity modifier calculation
func (suite *AbilityScoreTestSuite) TestModifierDexterity() {
	modifier := suite.scores.Modifier(core.Dexterity)
	assert.Equal(suite.T(), 2, modifier) // (14-10)/2 = 2
}

// TestModifierConstitution tests Constitution modifier calculation
func (suite *AbilityScoreTestSuite) TestModifierConstitution() {
	modifier := suite.scores.Modifier(core.Constitution)
	assert.Equal(suite.T(), 2, modifier) // (15-10)/2 = 2
}

// TestModifierIntelligence tests Intelligence modifier calculation
func (suite *AbilityScoreTestSuite) TestModifierIntelligence() {
	modifier := suite.scores.Modifier(core.Intelligence)
	assert.Equal(suite.T(), 1, modifier) // (12-10)/2 = 1
}

// TestModifierWisdom tests Wisdom modifier calculation
func (suite *AbilityScoreTestSuite) TestModifierWisdom() {
	modifier := suite.scores.Modifier(core.Wisdom)
	assert.Equal(suite.T(), 0, modifier) // (10-10)/2 = 0
}

// TestModifierCharisma tests Charisma modifier calculation
func (suite *AbilityScoreTestSuite) TestModifierCharisma() {
	modifier := suite.scores.Modifier(core.Charisma)
	assert.Equal(suite.T(), -1, modifier) // (8-10)/2 = -1
}

// TestGetEndpoint tests the GetEndpoint method
func (suite *AbilityScoreTestSuite) TestGetEndpoint() {
	endpoint := suite.scores.GetEndpoint()
	assert.Equal(suite.T(), "ability-scores/", endpoint)
}

// TestPrint tests that Print doesn't panic
func (suite *AbilityScoreTestSuite) TestPrint() {
	assert.NotPanics(suite.T(), func() {
		suite.scores.Print()
	})
}

func TestAbilityScoreTestSuite(t *testing.T) {
	suite.Run(t, new(AbilityScoreTestSuite))
}

// TestModifierCalculationEdgeCases tests edge cases for modifier calculation
func TestModifierCalculationEdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		score    int
		expected int
	}{
		{"Score of 1", 1, -5},   // (1-10)/2 = -4.5 -> -4
		{"Score of 3", 3, -4},   // (3-10)/2 = -3.5 -> -4
		{"Score of 8", 8, -1},   // (8-10)/2 = -1
		{"Score of 9", 9, -1},   // (9-10)/2 = -0.5 -> -1
		{"Score of 10", 10, 0},  // (10-10)/2 = 0
		{"Score of 11", 11, 0},  // (11-10)/2 = 0.5 -> 0
		{"Score of 12", 12, 1},  // (12-10)/2 = 1
		{"Score of 18", 18, 4},  // (18-10)/2 = 4
		{"Score of 20", 20, 5},  // (20-10)/2 = 5
		{"Score of 30", 30, 10}, // (30-10)/2 = 10
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scores := &abilities.AbilityScore{
				Strength: tc.score,
			}
			modifier := scores.Modifier(core.Strength)
			assert.Equal(t, tc.expected, modifier)
		})
	}
}

// TestModifierAllAbilities tests modifier calculation for all abilities
func TestModifierAllAbilities(t *testing.T) {
	scores := &abilities.AbilityScore{
		Strength:     18,
		Dexterity:    16,
		Constitution: 14,
		Intelligence: 12,
		Wisdom:       10,
		Charisma:     8,
	}

	testCases := []struct {
		ability  core.Ability
		expected int
	}{
		{core.Strength, 4},     // (18-10)/2 = 4
		{core.Dexterity, 3},    // (16-10)/2 = 3
		{core.Constitution, 2}, // (14-10)/2 = 2
		{core.Intelligence, 1}, // (12-10)/2 = 1
		{core.Wisdom, 0},       // (10-10)/2 = 0
		{core.Charisma, -1},    // (8-10)/2 = -1
	}

	for _, tc := range testCases {
		t.Run(tc.ability.String(), func(t *testing.T) {
			modifier := scores.Modifier(tc.ability)
			assert.Equal(t, tc.expected, modifier)
		})
	}
}

// TestEmptyAbilityScore tests zero-initialized ability scores
func TestEmptyAbilityScore(t *testing.T) {
	scores := &abilities.AbilityScore{}

	assert.Equal(t, 0, scores.Strength)
	assert.Equal(t, 0, scores.Dexterity)
	assert.Equal(t, 0, scores.Constitution)
	assert.Equal(t, 0, scores.Intelligence)
	assert.Equal(t, 0, scores.Wisdom)
	assert.Equal(t, 0, scores.Charisma)

	// All modifiers should be -5 for score of 0
	assert.Equal(t, -5, scores.Modifier(core.Strength))
	assert.Equal(t, -5, scores.Modifier(core.Dexterity))
}

// TestStandardArray tests the standard array ability scores
func TestStandardArray(t *testing.T) {
	// Standard array: 15, 14, 13, 12, 10, 8
	scores := &abilities.AbilityScore{
		Strength:     15,
		Dexterity:    14,
		Constitution: 13,
		Intelligence: 12,
		Wisdom:       10,
		Charisma:     8,
	}

	assert.Equal(t, 2, scores.Modifier(core.Strength))     // +2
	assert.Equal(t, 2, scores.Modifier(core.Dexterity))    // +2
	assert.Equal(t, 1, scores.Modifier(core.Constitution)) // +1
	assert.Equal(t, 1, scores.Modifier(core.Intelligence)) // +1
	assert.Equal(t, 0, scores.Modifier(core.Wisdom))       // +0
	assert.Equal(t, -1, scores.Modifier(core.Charisma))    // -1
}

// TestPointBuyMaximums tests point-buy maximum scores
func TestPointBuyMaximums(t *testing.T) {
	// Point buy allows max 15 before racial bonuses
	scores := &abilities.AbilityScore{
		Strength:     15,
		Dexterity:    15,
		Constitution: 15,
		Intelligence: 15,
		Wisdom:       15,
		Charisma:     15,
	}

	// All should have +2 modifier
	for _, ability := range []core.Ability{
		core.Strength, core.Dexterity, core.Constitution,
		core.Intelligence, core.Wisdom, core.Charisma,
	} {
		assert.Equal(t, 2, scores.Modifier(ability))
	}
}

// TestHighLevelAbilityScores tests scores with ASIs applied
func TestHighLevelAbilityScores(t *testing.T) {
	// Character with multiple ASIs (max is 20)
	scores := &abilities.AbilityScore{
		Strength:     20,
		Dexterity:    20,
		Constitution: 18,
		Intelligence: 16,
		Wisdom:       14,
		Charisma:     12,
	}

	assert.Equal(t, 5, scores.Modifier(core.Strength))
	assert.Equal(t, 5, scores.Modifier(core.Dexterity))
	assert.Equal(t, 4, scores.Modifier(core.Constitution))
	assert.Equal(t, 3, scores.Modifier(core.Intelligence))
	assert.Equal(t, 2, scores.Modifier(core.Wisdom))
	assert.Equal(t, 1, scores.Modifier(core.Charisma))
}

// TestAbilityScoreModification tests that ability scores can be modified
func TestAbilityScoreModification(t *testing.T) {
	scores := &abilities.AbilityScore{
		Strength: 10,
	}

	assert.Equal(t, 0, scores.Modifier(core.Strength))

	// Apply bonus (like racial bonus or ASI)
	scores.Strength += 2
	assert.Equal(t, 12, scores.Strength)
	assert.Equal(t, 1, scores.Modifier(core.Strength))

	scores.Strength += 4
	assert.Equal(t, 16, scores.Strength)
	assert.Equal(t, 3, scores.Modifier(core.Strength))
}

// TestNegativeAbilityScores tests negative ability scores (unusual but possible)
func TestNegativeAbilityScores(t *testing.T) {
	scores := &abilities.AbilityScore{
		Strength: -2,
	}

	// (-2-10)/2 = -12/2 = -6
	assert.Equal(t, -6, scores.Modifier(core.Strength))
}

// BenchmarkModifierCalculation benchmarks the Modifier method
func BenchmarkModifierCalculation(b *testing.B) {
	scores := &abilities.AbilityScore{
		Strength: 16,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = scores.Modifier(core.Strength)
	}
}

// BenchmarkModifierAllAbilities benchmarks calculating all modifiers
func BenchmarkModifierAllAbilities(b *testing.B) {
	scores := &abilities.AbilityScore{
		Strength:     16,
		Dexterity:    14,
		Constitution: 15,
		Intelligence: 12,
		Wisdom:       10,
		Charisma:     8,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = scores.Modifier(core.Strength)
		_ = scores.Modifier(core.Dexterity)
		_ = scores.Modifier(core.Constitution)
		_ = scores.Modifier(core.Intelligence)
		_ = scores.Modifier(core.Wisdom)
		_ = scores.Modifier(core.Charisma)
	}
}

// ExampleAbilityScore demonstrates basic AbilityScore usage
func ExampleAbilityScore() {
	scores := &abilities.AbilityScore{
		Strength:     16,
		Dexterity:    14,
		Constitution: 15,
		Intelligence: 12,
		Wisdom:       10,
		Charisma:     8,
	}

	strMod := scores.Modifier(core.Strength)
	_ = strMod // 3
}

// ExampleAbilityScore_Modifier demonstrates the Modifier method
func ExampleAbilityScore_Modifier() {
	scores := &abilities.AbilityScore{
		Wisdom: 16,
	}

	wisMod := scores.Modifier(core.Wisdom)
	_ = wisMod // 3 (because (16-10)/2 = 3)
}
