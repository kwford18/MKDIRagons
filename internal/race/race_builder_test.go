package race_test

import (
	"errors"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/internal/race"
	"github.com/kwford18/MKDIRagons/internal/reference"
	"github.com/kwford18/MKDIRagons/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ============================================================================
// MOCK FETCHERS
// ============================================================================

type MockFetcher struct {
	mock.Mock
}

func (m *MockFetcher) FetchJSON(property reference.Fetchable, input string) error {
	args := m.Called(property, input)
	return args.Error(0)
}

type MockFetcherWithFixtures struct {
	mock.Mock
	t *testing.T
}

func NewMockFetcherWithFixtures(t *testing.T) *MockFetcherWithFixtures {
	return &MockFetcherWithFixtures{t: t}
}

func (m *MockFetcherWithFixtures) FetchJSON(property reference.Fetchable, input string) error {
	args := m.Called(property, input)

	if args.Error(0) == nil && input != "" {
		fixtureFile := input + ".json"
		// REPLACED: Local helper call with core.LoadFixtureInto
		core.LoadFixtureInto(m.t, fixtureFile, property)
	}

	return args.Error(0)
}

// ============================================================================
// TEST SUITE
// ============================================================================

type RaceBuilderTestSuite struct {
	suite.Suite
	mockFetcher         *MockFetcher
	fixtureBasedFetcher *MockFetcherWithFixtures
	baseCharacter       *template.Character
	raceData            *race.Race
}

func (suite *RaceBuilderTestSuite) SetupTest() {
	suite.mockFetcher = new(MockFetcher)
	suite.fixtureBasedFetcher = NewMockFetcherWithFixtures(suite.T())

	suite.baseCharacter = &template.Character{
		Name:  "Test Human",
		Level: 5,
		Race:  "human",
		Class: "fighter",
		AbilityScores: template.AbilityScores{
			Strength:     16,
			Dexterity:    14,
			Constitution: 15,
			Intelligence: 10,
			Wisdom:       12,
			Charisma:     8,
		},
	}

	suite.raceData = &race.Race{}
}

func TestRaceBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(RaceBuilderTestSuite))
}

// ============================================================================
// UNIT TESTS - Behavior Testing
// ============================================================================

func (suite *RaceBuilderTestSuite) TestFetchRaceWithFetcher_Success() {
	suite.mockFetcher.On("FetchJSON", suite.raceData, "human").Return(nil)

	err := race.FetchRaceWithFetcher(suite.mockFetcher, suite.baseCharacter, suite.raceData)

	assert.NoError(suite.T(), err)
	suite.mockFetcher.AssertExpectations(suite.T())
	suite.mockFetcher.AssertCalled(suite.T(), "FetchJSON", suite.raceData, "human")
}

func (suite *RaceBuilderTestSuite) TestFetchRaceWithFetcher_NetworkError() {
	expectedError := errors.New("network timeout")
	suite.mockFetcher.On("FetchJSON", suite.raceData, "human").Return(expectedError)

	err := race.FetchRaceWithFetcher(suite.mockFetcher, suite.baseCharacter, suite.raceData)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockFetcher.AssertExpectations(suite.T())
}

func (suite *RaceBuilderTestSuite) TestFetchRaceWithFetcher_NotFoundError() {
	notFoundError := errors.New("404 not found")
	suite.mockFetcher.On("FetchJSON", suite.raceData, "human").Return(notFoundError)

	err := race.FetchRaceWithFetcher(suite.mockFetcher, suite.baseCharacter, suite.raceData)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), notFoundError, err)
}

func (suite *RaceBuilderTestSuite) TestFetchRaceWithFetcher_NilFetcher() {
	assert.Panics(suite.T(), func() {
		_ = race.FetchRaceWithFetcher(nil, suite.baseCharacter, suite.raceData)
	})
}

func (suite *RaceBuilderTestSuite) TestFetchRaceWithFetcher_NilCharacter() {
	assert.Panics(suite.T(), func() {
		_ = race.FetchRaceWithFetcher(suite.mockFetcher, nil, suite.raceData)
	})
}

func (suite *RaceBuilderTestSuite) TestFetchRaceWithFetcher_EmptyRaceName() {
	emptyCharacter := &template.Character{
		Name:          "Test",
		Level:         1,
		Race:          "",
		Class:         "fighter",
		AbilityScores: template.AbilityScores{Strength: 10},
	}

	suite.mockFetcher.On("FetchJSON", suite.raceData, "").Return(nil)

	err := race.FetchRaceWithFetcher(suite.mockFetcher, emptyCharacter, suite.raceData)

	assert.NoError(suite.T(), err)
	suite.mockFetcher.AssertCalled(suite.T(), "FetchJSON", suite.raceData, "")
}

func (suite *RaceBuilderTestSuite) TestFetchRaceWithFetcher_MultipleRaces() {
	races := []string{"human", "elf", "dwarf", "halfling", "dragonborn", "gnome", "half-elf", "half-orc", "tiefling"}

	for _, raceName := range races {
		character := &template.Character{
			Name:          "Test",
			Level:         1,
			Race:          raceName,
			Class:         "fighter",
			AbilityScores: template.AbilityScores{Strength: 10},
		}

		mockFetcher := new(MockFetcher)
		raceData := &race.Race{}

		mockFetcher.On("FetchJSON", raceData, raceName).Return(nil)

		err := race.FetchRaceWithFetcher(mockFetcher, character, raceData)

		assert.NoError(suite.T(), err, "Failed for race: %s", raceName)
		mockFetcher.AssertExpectations(suite.T())
	}
}

// ============================================================================
// UNIT TESTS - With Fixture Data
// ============================================================================

func (suite *RaceBuilderTestSuite) TestFetchRaceWithFetcher_HumanWithFixture() {
	suite.fixtureBasedFetcher.On("FetchJSON", suite.raceData, "human").Return(nil)

	err := race.FetchRaceWithFetcher(suite.fixtureBasedFetcher, suite.baseCharacter, suite.raceData)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Human", suite.raceData.Name)
	assert.Equal(suite.T(), "human", suite.raceData.Index)
	assert.Equal(suite.T(), 30, suite.raceData.Speed)
	assert.Equal(suite.T(), "Medium", suite.raceData.Size)
	assert.Len(suite.T(), suite.raceData.AbilityBonuses, 6) // All abilities +1
	assert.NotEmpty(suite.T(), suite.raceData.Languages)

	suite.fixtureBasedFetcher.AssertExpectations(suite.T())
}

func (suite *RaceBuilderTestSuite) TestFetchRaceWithFetcher_ElfWithFixture() {
	elfCharacter := &template.Character{
		Name:  "Test Elf",
		Level: 3,
		Race:  "elf",
		Class: "wizard",
		AbilityScores: template.AbilityScores{
			Intelligence: 16,
			Dexterity:    14,
			Wisdom:       12,
		},
	}

	elfRace := &race.Race{}

	suite.fixtureBasedFetcher.On("FetchJSON", elfRace, "elf").Return(nil)

	err := race.FetchRaceWithFetcher(suite.fixtureBasedFetcher, elfCharacter, elfRace)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Elf", elfRace.Name)
	assert.Equal(suite.T(), "elf", elfRace.Index)
	assert.Equal(suite.T(), 30, elfRace.Speed)
	assert.Len(suite.T(), elfRace.AbilityBonuses, 1) // DEX +2
	assert.Equal(suite.T(), "dex", elfRace.AbilityBonuses[0].AbilityScore.Index)
	assert.Equal(suite.T(), 2, elfRace.AbilityBonuses[0].Bonus)
	assert.NotEmpty(suite.T(), elfRace.Traits)
	assert.NotEmpty(suite.T(), elfRace.Subraces)

	suite.fixtureBasedFetcher.AssertExpectations(suite.T())
}

func (suite *RaceBuilderTestSuite) TestFetchRaceWithFetcher_DwarfWithFixture() {
	dwarfCharacter := &template.Character{
		Name:  "Test Dwarf",
		Level: 5,
		Race:  "dwarf",
		Class: "cleric",
		AbilityScores: template.AbilityScores{
			Wisdom:       16,
			Constitution: 14,
			Strength:     12,
		},
	}

	dwarfRace := &race.Race{}

	suite.fixtureBasedFetcher.On("FetchJSON", dwarfRace, "dwarf").Return(nil)

	err := race.FetchRaceWithFetcher(suite.fixtureBasedFetcher, dwarfCharacter, dwarfRace)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Dwarf", dwarfRace.Name)
	assert.Equal(suite.T(), "dwarf", dwarfRace.Index)
	assert.Equal(suite.T(), 25, dwarfRace.Speed)       // Dwarves are slower
	assert.Len(suite.T(), dwarfRace.AbilityBonuses, 1) // CON +2
	assert.Equal(suite.T(), "con", dwarfRace.AbilityBonuses[0].AbilityScore.Index)
	assert.Equal(suite.T(), 2, dwarfRace.AbilityBonuses[0].Bonus)

	suite.fixtureBasedFetcher.AssertExpectations(suite.T())
}

func (suite *RaceBuilderTestSuite) TestFetchRaceWithFetcher_VerifyAbilityBonuses() {
	suite.fixtureBasedFetcher.On("FetchJSON", suite.raceData, "human").Return(nil)

	err := race.FetchRaceWithFetcher(suite.fixtureBasedFetcher, suite.baseCharacter, suite.raceData)

	assert.NoError(suite.T(), err)

	// Verify all ability scores get +1 for human
	abilityIndexes := []string{"str", "dex", "con", "int", "wis", "cha"}
	for _, abilityIndex := range abilityIndexes {
		found := false
		for _, bonus := range suite.raceData.AbilityBonuses {
			if bonus.AbilityScore.Index == abilityIndex {
				assert.Equal(suite.T(), 1, bonus.Bonus)
				found = true
				break
			}
		}
		assert.True(suite.T(), found, "Expected to find ability bonus for %s", abilityIndex)
	}
}

func (suite *RaceBuilderTestSuite) TestFetchRaceWithFetcher_VerifyTraits() {
	elfRace := &race.Race{}
	elfCharacter := &template.Character{
		Name:          "Test Elf",
		Level:         1,
		Race:          "elf",
		Class:         "ranger",
		AbilityScores: template.AbilityScores{Dexterity: 16},
	}

	suite.fixtureBasedFetcher.On("FetchJSON", elfRace, "elf").Return(nil)

	err := race.FetchRaceWithFetcher(suite.fixtureBasedFetcher, elfCharacter, elfRace)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), elfRace.Traits)

	traitNames := make([]string, len(elfRace.Traits))
	for i, trait := range elfRace.Traits {
		traitNames[i] = trait.Name
	}

	assert.Contains(suite.T(), traitNames, "Darkvision")
	assert.Contains(suite.T(), traitNames, "Keen Senses")
	assert.Contains(suite.T(), traitNames, "Fey Ancestry")
}

// ============================================================================
// INTEGRATION TESTS - Real API (Optional)
// ============================================================================

func (suite *RaceBuilderTestSuite) TestFetchRace_RealAPI_Human() {
	if testing.Short() {
		suite.T().Skip("Skipping integration test in short mode")
	}

	character := &template.Character{
		Name:          "Integration Test Human",
		Level:         1,
		Race:          "human",
		Class:         "fighter",
		AbilityScores: template.AbilityScores{Strength: 16},
	}

	raceData := &race.Race{}

	err := race.FetchRace(character, raceData)

	if err == nil {
		assert.Equal(suite.T(), "Human", raceData.Name)
		assert.Equal(suite.T(), 30, raceData.Speed)
		suite.T().Log("✓ Real API call succeeded")
	} else {
		suite.T().Logf("⚠ API not available: %v", err)
	}
}

// ============================================================================
// EDGE CASE TESTS
// ============================================================================

func (suite *RaceBuilderTestSuite) TestFetchRaceWithFetcher_SequentialFetches() {
	suite.mockFetcher.On("FetchJSON", suite.raceData, "human").Return(nil).Once()
	err1 := race.FetchRaceWithFetcher(suite.mockFetcher, suite.baseCharacter, suite.raceData)
	assert.NoError(suite.T(), err1)

	elfCharacter := &template.Character{
		Name:          "Elf",
		Level:         1,
		Race:          "elf",
		Class:         "wizard",
		AbilityScores: template.AbilityScores{Intelligence: 16},
	}
	elfRace := &race.Race{}

	suite.mockFetcher.On("FetchJSON", elfRace, "elf").Return(nil).Once()
	err2 := race.FetchRaceWithFetcher(suite.mockFetcher, elfCharacter, elfRace)
	assert.NoError(suite.T(), err2)

	suite.mockFetcher.AssertNumberOfCalls(suite.T(), "FetchJSON", 2)
}

func (suite *RaceBuilderTestSuite) TestFetchRaceWithFetcher_RaceNameCaseSensitivity() {
	testCases := []string{"human", "Human", "HUMAN"}

	for _, raceName := range testCases {
		character := &template.Character{
			Name:          "Test",
			Level:         1,
			Race:          raceName,
			Class:         "fighter",
			AbilityScores: template.AbilityScores{Strength: 10},
		}

		mockFetcher := new(MockFetcher)
		raceData := &race.Race{}

		mockFetcher.On("FetchJSON", raceData, raceName).Return(nil)

		err := race.FetchRaceWithFetcher(mockFetcher, character, raceData)

		assert.NoError(suite.T(), err, "Failed for race name: %s", raceName)
		mockFetcher.AssertCalled(suite.T(), "FetchJSON", raceData, raceName)
	}
}
