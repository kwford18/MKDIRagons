package core_test

import (
	"testing"

	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// AbilityTestSuite defines the test suite for Ability enum
type AbilityTestSuite struct {
	suite.Suite
}

// TestAbilityEnumValues tests that enum values are correctly assigned
func (suite *AbilityTestSuite) TestAbilityEnumValues() {
	assert.Equal(suite.T(), 0, int(core.Strength))
	assert.Equal(suite.T(), 1, int(core.Dexterity))
	assert.Equal(suite.T(), 2, int(core.Constitution))
	assert.Equal(suite.T(), 3, int(core.Intelligence))
	assert.Equal(suite.T(), 4, int(core.Wisdom))
	assert.Equal(suite.T(), 5, int(core.Charisma))
}

// TestAbilityString tests the String() method for all abilities
func (suite *AbilityTestSuite) TestAbilityString() {
	assert.Equal(suite.T(), "Strength", core.Strength.String())
	assert.Equal(suite.T(), "Dexterity", core.Dexterity.String())
	assert.Equal(suite.T(), "Constitution", core.Constitution.String())
	assert.Equal(suite.T(), "Intelligence", core.Intelligence.String())
	assert.Equal(suite.T(), "Wisdom", core.Wisdom.String())
	assert.Equal(suite.T(), "Charisma", core.Charisma.String())
}

// TestAbilityStringConsistency tests that string values are unique
func (suite *AbilityTestSuite) TestAbilityStringConsistency() {
	abilities := []core.Ability{
		core.Strength,
		core.Dexterity,
		core.Constitution,
		core.Intelligence,
		core.Wisdom,
		core.Charisma,
	}

	// Check all strings are unique
	seen := make(map[string]bool)
	for _, ability := range abilities {
		str := ability.String()
		assert.False(suite.T(), seen[str], "Duplicate string found: %s", str)
		seen[str] = true
	}

	// Should have exactly 6 unique strings
	assert.Len(suite.T(), seen, 6)
}

// TestAbilityEnumOrdering tests that abilities maintain expected order
func (suite *AbilityTestSuite) TestAbilityEnumOrdering() {
	assert.True(suite.T(), core.Strength < core.Dexterity)
	assert.True(suite.T(), core.Dexterity < core.Constitution)
	assert.True(suite.T(), core.Constitution < core.Intelligence)
	assert.True(suite.T(), core.Intelligence < core.Wisdom)
	assert.True(suite.T(), core.Wisdom < core.Charisma)
}

func TestAbilityTestSuite(t *testing.T) {
	suite.Run(t, new(AbilityTestSuite))
}

// TestAbilityStringTableDriven uses table-driven tests for String() method
func TestAbilityStringTableDriven(t *testing.T) {
	testCases := []struct {
		name     string
		ability  core.Ability
		expected string
	}{
		{"Strength", core.Strength, "Strength"},
		{"Dexterity", core.Dexterity, "Dexterity"},
		{"Constitution", core.Constitution, "Constitution"},
		{"Intelligence", core.Intelligence, "Intelligence"},
		{"Wisdom", core.Wisdom, "Wisdom"},
		{"Charisma", core.Charisma, "Charisma"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.ability.String()
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestAbilityAsMapKey tests that Ability can be used as a map key
func TestAbilityAsMapKey(t *testing.T) {
	abilityMap := map[core.Ability]int{
		core.Strength:     10,
		core.Dexterity:    12,
		core.Constitution: 14,
		core.Intelligence: 16,
		core.Wisdom:       18,
		core.Charisma:     20,
	}

	assert.Equal(t, 10, abilityMap[core.Strength])
	assert.Equal(t, 12, abilityMap[core.Dexterity])
	assert.Equal(t, 14, abilityMap[core.Constitution])
	assert.Equal(t, 16, abilityMap[core.Intelligence])
	assert.Equal(t, 18, abilityMap[core.Wisdom])
	assert.Equal(t, 20, abilityMap[core.Charisma])
}

// TestAbilityInSwitch tests that Ability works correctly in switch statements
func TestAbilityInSwitch(t *testing.T) {
	getAbilityName := func(a core.Ability) string {
		switch a {
		case core.Strength:
			return "STR"
		case core.Dexterity:
			return "DEX"
		case core.Constitution:
			return "CON"
		case core.Intelligence:
			return "INT"
		case core.Wisdom:
			return "WIS"
		case core.Charisma:
			return "CHA"
		default:
			return "UNKNOWN"
		}
	}

	assert.Equal(t, "STR", getAbilityName(core.Strength))
	assert.Equal(t, "DEX", getAbilityName(core.Dexterity))
	assert.Equal(t, "CON", getAbilityName(core.Constitution))
	assert.Equal(t, "INT", getAbilityName(core.Intelligence))
	assert.Equal(t, "WIS", getAbilityName(core.Wisdom))
	assert.Equal(t, "CHA", getAbilityName(core.Charisma))
}

// TestAbilityIteration tests iterating over all abilities
func TestAbilityIteration(t *testing.T) {
	abilities := []core.Ability{
		core.Strength,
		core.Dexterity,
		core.Constitution,
		core.Intelligence,
		core.Wisdom,
		core.Charisma,
	}

	count := 0
	for _, ability := range abilities {
		assert.NotEmpty(t, ability.String())
		count++
	}

	assert.Equal(t, 6, count)
}

// BenchmarkAbilityString benchmarks the String() method
func BenchmarkAbilityString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = core.Strength.String()
	}
}

// BenchmarkAbilitySwitch benchmarks switch statements with Ability
func BenchmarkAbilitySwitch(b *testing.B) {
	ability := core.Dexterity
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result int
		switch ability {
		case core.Strength:
			result = 1
		case core.Dexterity:
			result = 2
		case core.Constitution:
			result = 3
		case core.Intelligence:
			result = 4
		case core.Wisdom:
			result = 5
		case core.Charisma:
			result = 6
		}
		_ = result
	}
}

// ExampleAbility_String demonstrates the String() method
func ExampleAbility_String() {
	str := core.Wisdom.String()
	_ = str // "Wisdom"
}
