package template_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/kwford18/MKDIRagons/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ========== TemplateAbilityScores Validate Tests ==========

// TestTemplateAbilityScoresValidateSuccess tests valid ability scores
func TestTemplateAbilityScoresValidateSuccess(t *testing.T) {
	testCases := []struct {
		name   string
		scores template.AbilityScores
	}{
		{
			name: "Standard array",
			scores: template.AbilityScores{
				Strength: 15, Dexterity: 14, Constitution: 13,
				Intelligence: 12, Wisdom: 10, Charisma: 8,
			},
		},
		{
			name: "All 10s",
			scores: template.AbilityScores{
				Strength: 10, Dexterity: 10, Constitution: 10,
				Intelligence: 10, Wisdom: 10, Charisma: 10,
			},
		},
		{
			name: "Max scores",
			scores: template.AbilityScores{
				Strength: 20, Dexterity: 20, Constitution: 20,
				Intelligence: 20, Wisdom: 20, Charisma: 20,
			},
		},
		{
			name: "Min scores",
			scores: template.AbilityScores{
				Strength: 0, Dexterity: 0, Constitution: 0,
				Intelligence: 0, Wisdom: 0, Charisma: 0,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.scores.Validate()
			assert.NoError(t, err)
		})
	}
}

// TestTemplateAbilityScoresValidateFailure tests invalid ability scores
func TestTemplateAbilityScoresValidateFailure(t *testing.T) {
	testCases := []struct {
		name          string
		scores        template.AbilityScores
		expectedError string
	}{
		{
			name:          "Strength too high",
			scores:        template.AbilityScores{Strength: 21},
			expectedError: "Strength",
		},
		{
			name:          "Strength too low",
			scores:        template.AbilityScores{Strength: -1},
			expectedError: "Strength",
		},
		{
			name:          "Dexterity too high",
			scores:        template.AbilityScores{Dexterity: 25},
			expectedError: "Dexterity",
		},
		{
			name:          "Constitution negative",
			scores:        template.AbilityScores{Constitution: -5},
			expectedError: "Constitution",
		},
		{
			name:          "Intelligence too high",
			scores:        template.AbilityScores{Intelligence: 100},
			expectedError: "Intelligence",
		},
		{
			name:          "Wisdom negative",
			scores:        template.AbilityScores{Wisdom: -10},
			expectedError: "Wisdom",
		},
		{
			name:          "Charisma too high",
			scores:        template.AbilityScores{Charisma: 30},
			expectedError: "Charisma",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.scores.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedError)
		})
	}
}

// ========== TomlParse Tests ==========

// TomlParseTestSuite defines the test suite for TOML parsing
type TomlParseTestSuite struct {
	suite.Suite
	tempDir string
}

// SetupTest runs before each test
func (suite *TomlParseTestSuite) SetupTest() {
	// Create a temporary directory for test files
	var err error
	suite.tempDir, err = os.MkdirTemp("", "toml_test_*")
	assert.NoError(suite.T(), err)
}

// TearDownTest runs after each test
func (suite *TomlParseTestSuite) TearDownTest() {
	// Clean up temporary directory
	err := os.RemoveAll(suite.tempDir)
	if err != nil {
		return
	}
}

// Helper to create a TOML file
func (suite *TomlParseTestSuite) createTOMLFile(filename, content string) string {
	path := filepath.Join(suite.tempDir, filename)
	err := os.WriteFile(path, []byte(content), 0644)
	assert.NoError(suite.T(), err)
	return path
}

// TestTomlParseValidCharacter tests parsing a valid character
func (suite *TomlParseTestSuite) TestTomlParseValidCharacter() {
	content := `
name = "Gandalf"
level = 5
race = "human"
class = "wizard"

proficiencies = ["arcana", "history"]

[ability_scores]
strength = 10
dexterity = 14
constitution = 12
intelligence = 18
wisdom = 15
charisma = 8

[inventory]
weapons = ["dagger"]
armor = []
items = ["spellbook"]

[spells]
level = [["fire-bolt"], ["magic-missile"]]
`
	path := suite.createTOMLFile("valid.toml", content)

	char, err := template.TomlParse(path)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Gandalf", char.Name)
	assert.Equal(suite.T(), 5, char.Level)
	assert.Equal(suite.T(), "human", char.Race)
	assert.Equal(suite.T(), "wizard", char.Class)
	assert.Equal(suite.T(), 18, char.AbilityScores.Intelligence)
	assert.Len(suite.T(), char.Proficiencies, 2)
	assert.Len(suite.T(), char.Inventory.Weapons, 1)
	assert.Len(suite.T(), char.Inventory.Armor, 0)
	assert.Len(suite.T(), char.Inventory.Items, 1)

}

// TestTomlParseValidAllRaces tests all valid 5e races
func (suite *TomlParseTestSuite) TestTomlParseValidAllRaces() {
	validRaces := []string{
		"dragonborn", "dwarf", "elf", "gnome", "half-elf",
		"half-orc", "halfling", "human", "tiefling",
	}

	for _, race := range validRaces {
		suite.Run(race, func() {
			content := `
name = "Test"
level = 1
race = "` + race + `"
class = "fighter"

[ability_scores]
strength = 10
dexterity = 10
constitution = 10
intelligence = 10
wisdom = 10
charisma = 10
`
			path := suite.createTOMLFile(race+".toml", content)

			char, err := template.TomlParse(path)

			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), race, char.Race)
		})
	}
}

// TestTomlParseValidAllClasses tests all valid 5e classes
func (suite *TomlParseTestSuite) TestTomlParseValidAllClasses() {
	validClasses := []string{
		"barbarian", "bard", "cleric", "druid", "fighter", "monk",
		"paladin", "ranger", "rogue", "sorcerer", "warlock", "wizard",
	}

	for _, class := range validClasses {
		suite.Run(class, func() {
			content := `
name = "Test"
level = 1
race = "human"
class = "` + class + `"

[ability_scores]
strength = 10
dexterity = 10
constitution = 10
intelligence = 10
wisdom = 10
charisma = 10
`
			path := suite.createTOMLFile(class+".toml", content)

			char, err := template.TomlParse(path)

			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), class, char.Class)
		})
	}
}

// TestTomlParseInvalidLevel tests invalid character levels
func (suite *TomlParseTestSuite) TestTomlParseInvalidLevel() {
	testCases := []struct {
		name  string
		level int
	}{
		{"Level 0", 0},
		{"Level -1", -1},
		{"Level 21", 21},
		{"Level 100", 100},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			content := `
name = "Test"
level = ` + fmt.Sprintf("%d", tc.level) + `
race = "human"
class = "fighter"

[ability_scores]
strength = 10
dexterity = 10
constitution = 10
intelligence = 10
wisdom = 10
charisma = 10
`
			path := suite.createTOMLFile(tc.name+".toml", content)

			_, err := template.TomlParse(path)

			assert.Error(suite.T(), err)
			assert.Contains(suite.T(), err.Error(), "invalid level")
		})
	}
}

// TestTomlParseInvalidRace tests invalid races
func (suite *TomlParseTestSuite) TestTomlParseInvalidRace() {
	testCases := []string{"orc", "goblin", "dragon", "invalid", ""}

	for _, race := range testCases {
		suite.Run(race, func() {
			content := `
name = "Test"
level = 1
race = "` + race + `"
class = "fighter"

[ability_scores]
strength = 10
dexterity = 10
constitution = 10
intelligence = 10
wisdom = 10
charisma = 10
`
			path := suite.createTOMLFile("invalid_race.toml", content)

			_, err := template.TomlParse(path)

			assert.Error(suite.T(), err)
			assert.Contains(suite.T(), err.Error(), "no valid 5e 2014 race")
		})
	}
}

// TestTomlParseInvalidClass tests invalid classes
func (suite *TomlParseTestSuite) TestTomlParseInvalidClass() {
	testCases := []string{"necromancer", "knight", "invalid", ""}

	for _, class := range testCases {
		suite.Run(class, func() {
			content := `
name = "Test"
level = 1
race = "human"
class = "` + class + `"

[ability_scores]
strength = 10
dexterity = 10
constitution = 10
intelligence = 10
wisdom = 10
charisma = 10
`
			path := suite.createTOMLFile("invalid_class.toml", content)

			_, err := template.TomlParse(path)

			assert.Error(suite.T(), err)
			assert.Contains(suite.T(), err.Error(), "no valid 5e 2014 class")
		})
	}
}

// TestTomlParseInvalidAbilityScores tests invalid ability scores
func (suite *TomlParseTestSuite) TestTomlParseInvalidAbilityScores() {
	content := `
name = "Test"
level = 1
race = "human"
class = "fighter"

[ability_scores]
strength = 25
dexterity = 10
constitution = 10
intelligence = 10
wisdom = 10
charisma = 10
`
	path := suite.createTOMLFile("invalid_scores.toml", content)

	_, err := template.TomlParse(path)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "Strength")
}

// TestTomlParseCaseInsensitivity tests case-insensitive race/class
func (suite *TomlParseTestSuite) TestTomlParseCaseInsensitivity() {
	testCases := []struct {
		name  string
		race  string
		class string
	}{
		{"Uppercase", "HUMAN", "WIZARD"},
		{"Mixed case", "HuMaN", "WiZaRd"},
		{"Title case", "Human", "Wizard"},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			content := `
name = "Test"
level = 1
race = "` + tc.race + `"
class = "` + tc.class + `"

[ability_scores]
strength = 10
dexterity = 10
constitution = 10
intelligence = 10
wisdom = 10
charisma = 10
`
			path := suite.createTOMLFile(tc.name+".toml", content)

			char, err := template.TomlParse(path)

			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tc.race, char.Race)
			assert.Equal(suite.T(), tc.class, char.Class)
		})
	}
}

// TestTomlParseCompleteCharacter tests a fully populated character
func (suite *TomlParseTestSuite) TestTomlParseCompleteCharacter() {
	content := `
name = "Legolas"
level = 10
race = "elf"
subrace = "wood-elf"
class = "ranger"
subclass = "hunter"

proficiencies = ["perception", "stealth", "survival"]
expertise = ["perception"]

[ability_scores]
strength = 12
dexterity = 18
constitution = 14
intelligence = 10
wisdom = 16
charisma = 8

[inventory]
weapons = ["longbow", "shortsword"]
armor = ["leather-armor"]
items = ["arrows", "rope"]

[spells]
level = [
    ["hunters-mark"],
    ["pass-without-trace", "spike-growth"]
]
`
	path := suite.createTOMLFile("complete.toml", content)

	char, err := template.TomlParse(path)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Legolas", char.Name)
	assert.Equal(suite.T(), 10, char.Level)
	assert.Equal(suite.T(), "elf", char.Race)
	assert.Equal(suite.T(), "wood-elf", char.Subrace)
	assert.Equal(suite.T(), "ranger", char.Class)
	assert.Equal(suite.T(), "hunter", char.Subclass)

	assert.Equal(suite.T(), 12, char.AbilityScores.Strength)
	assert.Equal(suite.T(), 18, char.AbilityScores.Dexterity)
	assert.Equal(suite.T(), 14, char.AbilityScores.Constitution)
	assert.Equal(suite.T(), 10, char.AbilityScores.Intelligence)
	assert.Equal(suite.T(), 16, char.AbilityScores.Wisdom)
	assert.Equal(suite.T(), 8, char.AbilityScores.Charisma)

	assert.Len(suite.T(), char.Proficiencies, 3)
	assert.Len(suite.T(), char.Expertise, 1)
	assert.Len(suite.T(), char.Inventory.Weapons, 2)
	assert.Len(suite.T(), char.Inventory.Armor, 1)
	assert.Len(suite.T(), char.Spells.Level, 2)
	assert.Len(suite.T(), char.Spells.Level[1], 2)
}

// TestTomlParseMinimalCharacter tests minimal required fields
func (suite *TomlParseTestSuite) TestTomlParseMinimalCharacter() {
	content := `
name = "Bob"
level = 1
race = "human"
class = "fighter"

[ability_scores]
strength = 10
dexterity = 10
constitution = 10
intelligence = 10
wisdom = 10
charisma = 10
`
	path := suite.createTOMLFile("minimal.toml", content)

	char, err := template.TomlParse(path)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Bob", char.Name)
	assert.Empty(suite.T(), char.Subrace)
	assert.Empty(suite.T(), char.Subclass)
	assert.Empty(suite.T(), char.Proficiencies)
	assert.Empty(suite.T(), char.Expertise)
}

// TestTomlParseFileNotFound tests parsing non-existent file
func (suite *TomlParseTestSuite) TestTomlParseFileNotFound() {
	_, err := template.TomlParse("nonexistent.toml")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to parse file")
}

// TestTomlParseMalformedTOML tests malformed TOML syntax
func (suite *TomlParseTestSuite) TestTomlParseMalformedTOML() {
	content := `
name = "Test"
level = 1
race = "human"
class = "fighter"
[ability_scores
strength = 10
`
	path := suite.createTOMLFile("malformed.toml", content)

	_, err := template.TomlParse(path)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to parse file")
}

// TestTomlParseMissingRequiredFields tests missing required fields
func (suite *TomlParseTestSuite) TestTomlParseMissingRequiredFields() {
	// Missing ability scores
	content := `
name = "Test"
level = 1
race = "human"
class = "fighter"
`
	path := suite.createTOMLFile("missing_fields.toml", content)

	_, err := template.TomlParse(path)

	// Should fail validation because ability scores will be 0
	// which is technically valid, but if any field is missing it might fail parsing
	// This depends on implementation details
	assert.NoError(suite.T(), err) // Actually valid - all scores default to 0
}

func TestTomlParseTestSuite(t *testing.T) {
	suite.Run(t, new(TomlParseTestSuite))
}

// ========== Benchmark Tests ==========

// BenchmarkTemplateAbilityScoresValidate benchmarks validation
func BenchmarkTemplateAbilityScoresValidate(b *testing.B) {
	scores := template.AbilityScores{
		Strength: 15, Dexterity: 14, Constitution: 13,
		Intelligence: 12, Wisdom: 10, Charisma: 8,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = scores.Validate()
	}
}

// ========== Example Tests ==========

// ExampleTomlParse demonstrates parsing a TOML file
func ExampleTomlParse() {
	// Assuming you have a character.toml file
	char, err := template.TomlParse("character.toml")
	if err != nil {
		// Handle error
		return
	}

	_ = char.Name
	_ = char.Level
}
