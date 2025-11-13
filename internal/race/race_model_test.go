package race_test

import (
	"github.com/kwford18/MKDIRagons/internal/race"
	"github.com/kwford18/MKDIRagons/internal/reference"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

// RaceModelTestSuite defines the test suite for Race model
type RaceModelTestSuite struct {
	suite.Suite
	race *race.Race
}

// SetupTest runs before each test
func (suite *RaceModelTestSuite) SetupTest() {
	suite.race = &race.Race{
		Index:     "elf",
		Name:      "Elf",
		Speed:     30,
		Size:      "Medium",
		Alignment: "Chaotic Good",
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{
					Index: "dex",
					Name:  "DEX",
					URL:   "/api/ability-scores/dex",
				},
				Bonus: 2,
			},
		},
		StartingProficiencies: []reference.Reference{
			{Index: "skill-perception", Name: "Skill: Perception", URL: "/api/proficiencies/skill-perception"},
		},
		Languages: []reference.Reference{
			{Index: "common", Name: "Common", URL: "/api/languages/common"},
			{Index: "elvish", Name: "Elvish", URL: "/api/languages/elvish"},
		},
		Traits: []reference.Reference{
			{Index: "darkvision", Name: "Darkvision", URL: "/api/traits/darkvision"},
			{Index: "fey-ancestry", Name: "Fey Ancestry", URL: "/api/traits/fey-ancestry"},
		},
		URL: "/api/races/elf",
	}
}

// TestRaceBasicFields tests basic field assignment and retrieval
func (suite *RaceModelTestSuite) TestRaceBasicFields() {
	assert.Equal(suite.T(), "elf", suite.race.Index)
	assert.Equal(suite.T(), "Elf", suite.race.Name)
	assert.Equal(suite.T(), 30, suite.race.Speed)
	assert.Equal(suite.T(), "Medium", suite.race.Size)
	assert.Equal(suite.T(), "Chaotic Good", suite.race.Alignment)
	assert.Equal(suite.T(), "/api/races/elf", suite.race.URL)
}

// TestRaceAbilityBonuses tests ability bonuses
func (suite *RaceModelTestSuite) TestRaceAbilityBonuses() {
	assert.Len(suite.T(), suite.race.AbilityBonuses, 1)
	assert.Equal(suite.T(), "dex", suite.race.AbilityBonuses[0].AbilityScore.Index)
	assert.Equal(suite.T(), "DEX", suite.race.AbilityBonuses[0].AbilityScore.Name)
	assert.Equal(suite.T(), 2, suite.race.AbilityBonuses[0].Bonus)
}

// TestRaceStartingProficiencies tests starting proficiencies
func (suite *RaceModelTestSuite) TestRaceStartingProficiencies() {
	assert.Len(suite.T(), suite.race.StartingProficiencies, 1)
	assert.Equal(suite.T(), "skill-perception", suite.race.StartingProficiencies[0].Index)
	assert.Equal(suite.T(), "Skill: Perception", suite.race.StartingProficiencies[0].Name)
}

// TestRaceLanguages tests languages
func (suite *RaceModelTestSuite) TestRaceLanguages() {
	assert.Len(suite.T(), suite.race.Languages, 2)
	assert.Equal(suite.T(), "common", suite.race.Languages[0].Index)
	assert.Equal(suite.T(), "elvish", suite.race.Languages[1].Index)
}

// TestRaceTraits tests traits
func (suite *RaceModelTestSuite) TestRaceTraits() {
	assert.Len(suite.T(), suite.race.Traits, 2)
	assert.Equal(suite.T(), "darkvision", suite.race.Traits[0].Index)
	assert.Equal(suite.T(), "fey-ancestry", suite.race.Traits[1].Index)
}

// TestGetEndpoint tests the GetEndpoint method
func (suite *RaceModelTestSuite) TestGetEndpoint() {
	endpoint := suite.race.GetEndpoint()
	assert.Equal(suite.T(), "races/", endpoint)
}

// TestPrint tests the Print method doesn't panic
func (suite *RaceModelTestSuite) TestPrint() {
	assert.NotPanics(suite.T(), func() {
		suite.race.Print()
	})
}

// Run the test suite
func TestRaceModelTestSuite(t *testing.T) {
	suite.Run(t, new(RaceModelTestSuite))
}

// TestAbilityBonus tests AbilityBonus struct
func TestAbilityBonus(t *testing.T) {
	ab := race.AbilityBonus{
		AbilityScore: reference.Reference{
			Index: "str",
			Name:  "STR",
			URL:   "/api/ability-scores/str",
		},
		Bonus: 2,
	}

	assert.Equal(t, "str", ab.AbilityScore.Index)
	assert.Equal(t, "STR", ab.AbilityScore.Name)
	assert.Equal(t, 2, ab.Bonus)
}

// TestMultipleAbilityBonuses tests race with multiple ability bonuses
func TestMultipleAbilityBonuses(t *testing.T) {
	testRace := &race.Race{
		Index: "half-elf",
		Name:  "Half-Elf",
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{Index: "cha", Name: "CHA"},
				Bonus:        2,
			},
			{
				AbilityScore: reference.Reference{Index: "str", Name: "STR"},
				Bonus:        1,
			},
			{
				AbilityScore: reference.Reference{Index: "dex", Name: "DEX"},
				Bonus:        1,
			},
		},
	}

	assert.Len(t, testRace.AbilityBonuses, 3)
	assert.Equal(t, 2, testRace.AbilityBonuses[0].Bonus)
	assert.Equal(t, 1, testRace.AbilityBonuses[1].Bonus)
	assert.Equal(t, 1, testRace.AbilityBonuses[2].Bonus)
}

// TestEmptyRace tests empty race initialization
func TestEmptyRace(t *testing.T) {
	testRace := &race.Race{}

	assert.Empty(t, testRace.Index)
	assert.Empty(t, testRace.Name)
	assert.Equal(t, 0, testRace.Speed)
	assert.Empty(t, testRace.Size)
	assert.Nil(t, testRace.AbilityBonuses)
	assert.Nil(t, testRace.StartingProficiencies)
	assert.Nil(t, testRace.Languages)
	assert.Nil(t, testRace.Traits)
	assert.Nil(t, testRace.Subraces)
	assert.Equal(t, "races/", testRace.GetEndpoint())
}

// TestRaceWithSubraces tests race with subraces
func TestRaceWithSubraces(t *testing.T) {
	testRace := &race.Race{
		Index: "elf",
		Name:  "Elf",
		Subraces: []reference.Reference{
			{Index: "high-elf", Name: "High Elf", URL: "/api/subraces/high-elf"},
			{Index: "wood-elf", Name: "Wood Elf", URL: "/api/subraces/wood-elf"},
			{Index: "dark-elf", Name: "Dark Elf (Drow)", URL: "/api/subraces/dark-elf"},
		},
	}

	assert.Len(t, testRace.Subraces, 3)
	assert.Equal(t, "high-elf", testRace.Subraces[0].Index)
	assert.Equal(t, "wood-elf", testRace.Subraces[1].Index)
	assert.Equal(t, "dark-elf", testRace.Subraces[2].Index)
}

// TestComplexRace tests a fully populated race
func TestComplexRace(t *testing.T) {
	testRace := &race.Race{
		Index:           "dwarf",
		Name:            "Dwarf",
		Speed:           25,
		Size:            "Medium",
		SizeDescription: "Dwarves stand between 4 and 5 feet tall and average about 150 pounds.",
		Age:             "Dwarves mature at the same rate as humans, but they're considered young until they reach the age of 50.",
		Alignment:       "Most dwarves are lawful, believing firmly in the benefits of a well-ordered society.",
		LanguageDesc:    "You can speak, read, and write Common and Dwarvish.",
		AbilityBonuses: []race.AbilityBonus{
			{
				AbilityScore: reference.Reference{Index: "con", Name: "CON"},
				Bonus:        2,
			},
		},
		StartingProficiencies: []reference.Reference{
			{Index: "battleaxes", Name: "Battleaxes"},
			{Index: "handaxes", Name: "Handaxes"},
		},
		Languages: []reference.Reference{
			{Index: "common", Name: "Common"},
			{Index: "dwarvish", Name: "Dwarvish"},
		},
		Traits: []reference.Reference{
			{Index: "darkvision", Name: "Darkvision"},
			{Index: "dwarven-resilience", Name: "Dwarven Resilience"},
		},
		Subraces: []reference.Reference{
			{Index: "hill-dwarf", Name: "Hill Dwarf"},
			{Index: "mountain-dwarf", Name: "Mountain Dwarf"},
		},
		URL: "/api/races/dwarf",
	}

	assert.Equal(t, "dwarf", testRace.Index)
	assert.Equal(t, "Dwarf", testRace.Name)
	assert.Equal(t, 25, testRace.Speed)
	assert.NotEmpty(t, testRace.Age)
	assert.NotEmpty(t, testRace.Alignment)
	assert.NotEmpty(t, testRace.SizeDescription)
	assert.NotEmpty(t, testRace.LanguageDesc)
	assert.Len(t, testRace.AbilityBonuses, 1)
	assert.Len(t, testRace.StartingProficiencies, 2)
	assert.Len(t, testRace.Languages, 2)
	assert.Len(t, testRace.Traits, 2)
	assert.Len(t, testRace.Subraces, 2)
}

// TestRaceWithNoAbilityBonuses tests race without ability bonuses
func TestRaceWithNoAbilityBonuses(t *testing.T) {
	testRace := &race.Race{
		Index: "custom-race",
		Name:  "Custom Race",
		Speed: 30,
	}

	assert.Empty(t, testRace.AbilityBonuses)
	assert.Nil(t, testRace.AbilityBonuses)
}

// TestRaceDescriptionFields tests description fields
func TestRaceDescriptionFields(t *testing.T) {
	testRace := &race.Race{
		Index:           "human",
		Name:            "Human",
		Age:             "Humans reach adulthood in their late teens and live less than a century.",
		Alignment:       "Humans tend toward no particular alignment.",
		SizeDescription: "Humans vary widely in height and build, from barely 5 feet to well over 6 feet tall.",
		LanguageDesc:    "You can speak, read, and write Common and one extra language of your choice.",
	}

	assert.NotEmpty(t, testRace.Age)
	assert.NotEmpty(t, testRace.Alignment)
	assert.NotEmpty(t, testRace.SizeDescription)
	assert.NotEmpty(t, testRace.LanguageDesc)
	assert.Contains(t, testRace.Age, "adulthood")
	assert.Contains(t, testRace.SizeDescription, "height")
	assert.Contains(t, testRace.LanguageDesc, "Common")
}

// TestRaceWithVariousSpeedValues tests different speed values
func TestRaceWithVariousSpeedValues(t *testing.T) {
	testCases := []struct {
		name  string
		race  string
		speed int
	}{
		{"Dwarf (slow)", "dwarf", 25},
		{"Human (normal)", "human", 30},
		{"Wood Elf (fast)", "wood-elf", 35},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testRace := &race.Race{
				Index: tc.race,
				Speed: tc.speed,
			}
			assert.Equal(t, tc.speed, testRace.Speed)
		})
	}
}

// TestRaceSizes tests different size categories
func TestRaceSizes(t *testing.T) {
	testCases := []struct {
		name string
		size string
	}{
		{"Small race", "Small"},
		{"Medium race", "Medium"},
		{"Large race", "Large"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testRace := &race.Race{
				Index: "test-race",
				Size:  tc.size,
			}
			assert.Equal(t, tc.size, testRace.Size)
		})
	}
}

// TestRaceWithNoTraits tests race without traits
func TestRaceWithNoTraits(t *testing.T) {
	testRace := &race.Race{
		Index:  "basic-race",
		Name:   "Basic Race",
		Traits: []reference.Reference{},
	}

	assert.Empty(t, testRace.Traits)
	assert.NotNil(t, testRace.Traits)
}

// TestRaceWithManyLanguages tests race with multiple languages
func TestRaceWithManyLanguages(t *testing.T) {
	testRace := &race.Race{
		Index: "polyglot-race",
		Name:  "Polyglot Race",
		Languages: []reference.Reference{
			{Index: "common", Name: "Common"},
			{Index: "elvish", Name: "Elvish"},
			{Index: "dwarvish", Name: "Dwarvish"},
			{Index: "draconic", Name: "Draconic"},
		},
	}

	assert.Len(t, testRace.Languages, 4)
	var languageIndices []string
	for _, lang := range testRace.Languages {
		languageIndices = append(languageIndices, lang.Index)
	}
	assert.Contains(t, languageIndices, "common")
	assert.Contains(t, languageIndices, "elvish")
	assert.Contains(t, languageIndices, "dwarvish")
	assert.Contains(t, languageIndices, "draconic")
}
