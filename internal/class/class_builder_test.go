package class_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/class"
	"github.com/kwford18/MKDIRagons/internal/reference"
	"github.com/kwford18/MKDIRagons/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// loadFixture loads a JSON fixture file
func loadFixture(t *testing.T, filename string, target interface{}) {
	t.Helper()

	fixtureDir := filepath.Join("testdata", "fixtures")
	filePath := filepath.Join(fixtureDir, filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read fixture file %s: %v", filePath, err)
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		t.Fatalf("Failed to unmarshal fixture file %s: %v", filePath, err)
	}
}

// ============================================================================
// MOCK FETCHERS
// ============================================================================

// MockFetcher - simple mock for behavior testing
type MockFetcher struct {
	mock.Mock
}

func (m *MockFetcher) FetchJSON(property reference.Fetchable, input string) error {
	args := m.Called(property, input)
	return args.Error(0)
}

// MockFetcherWithFixtures - mock that loads real fixture data
type MockFetcherWithFixtures struct {
	mock.Mock
	t *testing.T
}

func NewMockFetcherWithFixtures(t *testing.T) *MockFetcherWithFixtures {
	return &MockFetcherWithFixtures{t: t}
}

func (m *MockFetcherWithFixtures) FetchJSON(property reference.Fetchable, input string) error {
	args := m.Called(property, input)

	// If mock expects success, load the fixture
	if args.Error(0) == nil && input != "" {
		fixtureFile := input + ".json"
		loadFixture(m.t, fixtureFile, property)
	}

	return args.Error(0)
}

// ============================================================================
// TEST SUITE
// ============================================================================

type ClassBuilderTestSuite struct {
	suite.Suite
	mockFetcher         *MockFetcher
	fixtureBasedFetcher *MockFetcherWithFixtures
	baseCharacter       *template.Character
	classData           *class.Class
}

func (suite *ClassBuilderTestSuite) SetupTest() {
	suite.mockFetcher = new(MockFetcher)
	suite.fixtureBasedFetcher = NewMockFetcherWithFixtures(suite.T())

	suite.baseCharacter = &template.Character{
		Name:  "Test Fighter",
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

	suite.classData = &class.Class{}
}

func TestClassBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(ClassBuilderTestSuite))
}

// ============================================================================
// UNIT TESTS - Behavior Testing (No Real Data)
// ============================================================================

func (suite *ClassBuilderTestSuite) TestFetchClassWithFetcher_Success() {
	suite.mockFetcher.On("FetchJSON", suite.classData, "fighter").Return(nil)

	err := class.FetchClassWithFetcher(suite.mockFetcher, suite.baseCharacter, suite.classData)

	assert.NoError(suite.T(), err)
	suite.mockFetcher.AssertExpectations(suite.T())
	suite.mockFetcher.AssertCalled(suite.T(), "FetchJSON", suite.classData, "fighter")
}

func (suite *ClassBuilderTestSuite) TestFetchClassWithFetcher_NetworkError() {
	expectedError := errors.New("network timeout")
	suite.mockFetcher.On("FetchJSON", suite.classData, "fighter").Return(expectedError)

	err := class.FetchClassWithFetcher(suite.mockFetcher, suite.baseCharacter, suite.classData)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
	suite.mockFetcher.AssertExpectations(suite.T())
}

func (suite *ClassBuilderTestSuite) TestFetchClassWithFetcher_NotFoundError() {
	notFoundError := errors.New("404 not found")
	suite.mockFetcher.On("FetchJSON", suite.classData, "fighter").Return(notFoundError)

	err := class.FetchClassWithFetcher(suite.mockFetcher, suite.baseCharacter, suite.classData)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), notFoundError, err)
}

func (suite *ClassBuilderTestSuite) TestFetchClassWithFetcher_NilFetcher() {
	assert.Panics(suite.T(), func() {
		_ = class.FetchClassWithFetcher(nil, suite.baseCharacter, suite.classData)
	})
}

func (suite *ClassBuilderTestSuite) TestFetchClassWithFetcher_NilCharacter() {
	assert.Panics(suite.T(), func() {
		_ = class.FetchClassWithFetcher(suite.mockFetcher, nil, suite.classData)
	})
}

func (suite *ClassBuilderTestSuite) TestFetchClassWithFetcher_EmptyClassName() {
	emptyCharacter := &template.Character{
		Name:          "Test",
		Level:         1,
		Race:          "human",
		Class:         "",
		AbilityScores: template.AbilityScores{Strength: 10},
	}

	suite.mockFetcher.On("FetchJSON", suite.classData, "").Return(nil)

	err := class.FetchClassWithFetcher(suite.mockFetcher, emptyCharacter, suite.classData)

	assert.NoError(suite.T(), err)
	suite.mockFetcher.AssertCalled(suite.T(), "FetchJSON", suite.classData, "")
}

func (suite *ClassBuilderTestSuite) TestFetchClassWithFetcher_MultipleClassNames() {
	classNames := []string{"barbarian", "bard", "cleric", "druid", "monk", "paladin", "ranger", "rogue", "sorcerer", "warlock"}

	for _, className := range classNames {
		character := &template.Character{
			Name:          "Test",
			Level:         1,
			Race:          "human",
			Class:         className,
			AbilityScores: template.AbilityScores{Strength: 10},
		}

		mockFetcher := new(MockFetcher)
		classData := &class.Class{}

		mockFetcher.On("FetchJSON", classData, className).Return(nil)

		err := class.FetchClassWithFetcher(mockFetcher, character, classData)

		assert.NoError(suite.T(), err, "Failed for class: %s", className)
		mockFetcher.AssertExpectations(suite.T())
	}
}

// ============================================================================
// UNIT TESTS - With Fixture Data (Realistic Testing)
// ============================================================================

func (suite *ClassBuilderTestSuite) TestFetchClassWithFetcher_FighterWithFixture() {
	suite.fixtureBasedFetcher.On("FetchJSON", suite.classData, "fighter").Return(nil)

	err := class.FetchClassWithFetcher(suite.fixtureBasedFetcher, suite.baseCharacter, suite.classData)

	assert.NoError(suite.T(), err)

	// Verify the data was properly loaded from fixture
	assert.Equal(suite.T(), "Fighter", suite.classData.Name)
	assert.Equal(suite.T(), "fighter", suite.classData.Index)
	assert.Equal(suite.T(), 10, suite.classData.HitDie)
	assert.NotEmpty(suite.T(), suite.classData.Proficiencies)
	assert.NotEmpty(suite.T(), suite.classData.SavingThrows)
	assert.Len(suite.T(), suite.classData.SavingThrows, 2) // STR and CON

	suite.fixtureBasedFetcher.AssertExpectations(suite.T())
}

func (suite *ClassBuilderTestSuite) TestFetchClassWithFetcher_WizardWithFixture() {
	wizardCharacter := &template.Character{
		Name:  "Test Wizard",
		Level: 3,
		Race:  "elf",
		Class: "wizard",
		AbilityScores: template.AbilityScores{
			Intelligence: 16,
			Wisdom:       14,
			Dexterity:    12,
		},
	}

	wizardClass := &class.Class{}

	suite.fixtureBasedFetcher.On("FetchJSON", wizardClass, "wizard").Return(nil)

	err := class.FetchClassWithFetcher(suite.fixtureBasedFetcher, wizardCharacter, wizardClass)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Wizard", wizardClass.Name)
	assert.Equal(suite.T(), "wizard", wizardClass.Index)
	assert.Equal(suite.T(), 6, wizardClass.HitDie)
	assert.NotNil(suite.T(), wizardClass.Spellcasting)
	assert.Equal(suite.T(), "int", wizardClass.Spellcasting.SpellcastingAbility.Index)

	suite.fixtureBasedFetcher.AssertExpectations(suite.T())
}

func (suite *ClassBuilderTestSuite) TestFetchClassWithFetcher_VerifyProficienciesLoaded() {
	suite.fixtureBasedFetcher.On("FetchJSON", suite.classData, "fighter").Return(nil)

	err := class.FetchClassWithFetcher(suite.fixtureBasedFetcher, suite.baseCharacter, suite.classData)

	assert.NoError(suite.T(), err)

	// Verify specific proficiencies from fixture
	assert.NotEmpty(suite.T(), suite.classData.Proficiencies)

	proficiencyNames := make([]string, len(suite.classData.Proficiencies))
	for i, prof := range suite.classData.Proficiencies {
		proficiencyNames[i] = prof.Name
	}

	assert.Contains(suite.T(), proficiencyNames, "All armor")
	assert.Contains(suite.T(), proficiencyNames, "Shields")
}

func (suite *ClassBuilderTestSuite) TestFetchClassWithFetcher_VerifySavingThrowsLoaded() {
	suite.fixtureBasedFetcher.On("FetchJSON", suite.classData, "fighter").Return(nil)

	err := class.FetchClassWithFetcher(suite.fixtureBasedFetcher, suite.baseCharacter, suite.classData)

	assert.NoError(suite.T(), err)

	// Verify saving throws from fixture
	assert.Len(suite.T(), suite.classData.SavingThrows, 2)

	savingThrowIndexes := []string{
		suite.classData.SavingThrows[0].Index,
		suite.classData.SavingThrows[1].Index,
	}

	assert.Contains(suite.T(), savingThrowIndexes, "str")
	assert.Contains(suite.T(), savingThrowIndexes, "con")
}

// ============================================================================
// INTEGRATION TESTS - Real API Calls (Optional, Run Separately)
// ============================================================================

func (suite *ClassBuilderTestSuite) TestFetchClass_RealAPI_Fighter() {
	if testing.Short() {
		suite.T().Skip("Skipping integration test in short mode")
	}

	character := &template.Character{
		Name:  "Integration Test Fighter",
		Level: 5,
		Race:  "human",
		Class: "fighter",
		AbilityScores: template.AbilityScores{
			Strength: 16,
		},
	}

	classData := &class.Class{}

	err := class.FetchClass(character, classData)

	// If API is available, verify data; otherwise just ensure no panic
	if err == nil {
		assert.Equal(suite.T(), "Fighter", classData.Name)
		assert.Equal(suite.T(), 10, classData.HitDie)
		suite.T().Log("✓ Real API call succeeded")
	} else {
		suite.T().Logf("⚠ API not available: %v", err)
	}
}

func (suite *ClassBuilderTestSuite) TestFetchClass_RealAPI_InvalidClass() {
	if testing.Short() {
		suite.T().Skip("Skipping integration test in short mode")
	}

	character := &template.Character{
		Name:  "Invalid Class Test",
		Level: 1,
		Race:  "human",
		Class: "nonexistent-class",
		AbilityScores: template.AbilityScores{
			Strength: 10,
		},
	}

	classData := &class.Class{}

	err := class.FetchClass(character, classData)

	// Should error for invalid class
	assert.Error(suite.T(), err)
	suite.T().Logf("✓ Invalid class properly rejected: %v", err)
}

// ============================================================================
// EDGE CASE TESTS
// ============================================================================

func (suite *ClassBuilderTestSuite) TestFetchClassWithFetcher_SequentialFetches() {
	// First fetch - fighter
	suite.mockFetcher.On("FetchJSON", suite.classData, "fighter").Return(nil).Once()
	err1 := class.FetchClassWithFetcher(suite.mockFetcher, suite.baseCharacter, suite.classData)
	assert.NoError(suite.T(), err1)

	// Second fetch - wizard with different data
	wizardCharacter := &template.Character{
		Name:          "Wizard",
		Level:         1,
		Race:          "elf",
		Class:         "wizard",
		AbilityScores: template.AbilityScores{Intelligence: 16},
	}
	wizardClass := &class.Class{}

	suite.mockFetcher.On("FetchJSON", wizardClass, "wizard").Return(nil).Once()
	err2 := class.FetchClassWithFetcher(suite.mockFetcher, wizardCharacter, wizardClass)
	assert.NoError(suite.T(), err2)

	suite.mockFetcher.AssertNumberOfCalls(suite.T(), "FetchJSON", 2)
}

func (suite *ClassBuilderTestSuite) TestFetchClassWithFetcher_ClassNameCaseSensitivity() {
	testCases := []string{"fighter", "Fighter", "FIGHTER"}

	for _, className := range testCases {
		character := &template.Character{
			Name:          "Test",
			Level:         1,
			Race:          "human",
			Class:         className,
			AbilityScores: template.AbilityScores{Strength: 10},
		}

		mockFetcher := new(MockFetcher)
		classData := &class.Class{}

		mockFetcher.On("FetchJSON", classData, className).Return(nil)

		err := class.FetchClassWithFetcher(mockFetcher, character, classData)

		assert.NoError(suite.T(), err, "Failed for class name: %s", className)
		mockFetcher.AssertCalled(suite.T(), "FetchJSON", classData, className)
	}
}
