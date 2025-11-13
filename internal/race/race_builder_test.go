package race_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/internal/race"
	"github.com/kwford18/MKDIRagons/internal/reference"
	"github.com/kwford18/MKDIRagons/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ========== Mock Fetcher ==========

// MockFetcher is a mock implementation of core.Fetcher for unit tests
type MockFetcher struct {
	mock.Mock
}

func (m *MockFetcher) FetchJSON(property reference.Fetchable, input string) error {
	args := m.Called(property, input)
	return args.Error(0)
}

// ========== Unit Tests with Mock Fetcher ==========

// FetchRaceUnitTestSuite uses mocks for fast, isolated unit tests
type FetchRaceUnitTestSuite struct {
	suite.Suite
	mockFetcher *MockFetcher
	base        *templates.TemplateCharacter
	testRace    *race.Race
}

func (suite *FetchRaceUnitTestSuite) SetupTest() {
	suite.mockFetcher = new(MockFetcher)
	suite.base = &templates.TemplateCharacter{
		Race: "elf",
	}
	suite.testRace = &race.Race{}
}

func (suite *FetchRaceUnitTestSuite) TestFetchRaceSuccess() {
	// Setup mock to return success
	suite.mockFetcher.On("FetchJSON", suite.testRace, "elf").Return(nil).Run(func(args mock.Arguments) {
		// Simulate populating the race data
		r := args.Get(0).(*race.Race)
		r.Index = "elf"
		r.Name = "Elf"
		r.Speed = 30
		r.Size = "Medium"
	})

	err := race.FetchRaceWithFetcher(suite.mockFetcher, suite.base, suite.testRace)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "elf", suite.testRace.Index)
	assert.Equal(suite.T(), "Elf", suite.testRace.Name)
	assert.Equal(suite.T(), 30, suite.testRace.Speed)
	suite.mockFetcher.AssertExpectations(suite.T())
}

func (suite *FetchRaceUnitTestSuite) TestFetchRaceError() {
	// Setup mock to return an error
	expectedErr := errors.New("network error")
	suite.mockFetcher.On("FetchJSON", suite.testRace, "elf").Return(expectedErr)

	err := race.FetchRaceWithFetcher(suite.mockFetcher, suite.base, suite.testRace)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedErr, err)
	suite.mockFetcher.AssertExpectations(suite.T())
}

func (suite *FetchRaceUnitTestSuite) TestFetchRaceEmptyString() {
	suite.base.Race = ""
	suite.mockFetcher.On("FetchJSON", suite.testRace, "").Return(errors.New("empty race name"))

	err := race.FetchRaceWithFetcher(suite.mockFetcher, suite.base, suite.testRace)

	assert.Error(suite.T(), err)
	suite.mockFetcher.AssertExpectations(suite.T())
}

func (suite *FetchRaceUnitTestSuite) TestFetchRaceNilBase() {
	assert.Panics(suite.T(), func() {
		race.FetchRaceWithFetcher(suite.mockFetcher, nil, suite.testRace)
	})
}

func (suite *FetchRaceUnitTestSuite) TestFetchRaceNilRace() {
	assert.Panics(suite.T(), func() {
		race.FetchRaceWithFetcher(suite.mockFetcher, suite.base, nil)
	})
}

func (suite *FetchRaceUnitTestSuite) TestFetchRaceDifferentRaces() {
	testCases := []struct {
		raceName      string
		expectedName  string
		expectedSpeed int
	}{
		{"elf", "Elf", 30},
		{"dwarf", "Dwarf", 25},
		{"halfling", "Halfling", 25},
	}

	for _, tc := range testCases {
		suite.Run(tc.raceName, func() {
			mockFetcher := new(MockFetcher)
			base := &templates.TemplateCharacter{Race: tc.raceName}
			testRace := &race.Race{}

			mockFetcher.On("FetchJSON", testRace, tc.raceName).Return(nil).Run(func(args mock.Arguments) {
				r := args.Get(0).(*race.Race)
				r.Index = tc.raceName
				r.Name = tc.expectedName
				r.Speed = tc.expectedSpeed
			})

			err := race.FetchRaceWithFetcher(mockFetcher, base, testRace)

			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tc.raceName, testRace.Index)
			assert.Equal(suite.T(), tc.expectedSpeed, testRace.Speed)
			mockFetcher.AssertExpectations(suite.T())
		})
	}
}

func TestFetchRaceUnitTestSuite(t *testing.T) {
	suite.Run(t, new(FetchRaceUnitTestSuite))
}

// ========== Integration Tests with HTTP Server ==========

// FetchRaceIntegrationTestSuite uses a real HTTP server for integration tests
type FetchRaceIntegrationTestSuite struct {
	suite.Suite
	server  *httptest.Server
	fetcher *core.HTTPFetcher
}

func (suite *FetchRaceIntegrationTestSuite) SetupSuite() {
	// Create a mock HTTP server that mimics the D&D API
	suite.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/races/elf":
			mockRace := race.Race{
				Index:     "elf",
				Name:      "Elf",
				Speed:     30,
				Size:      "Medium",
				Alignment: "Chaotic Good",
				Age:       "Although elves reach physical maturity at about the same age as humans, the elven understanding of adulthood goes beyond physical growth to encompass worldly experience. An elf typically claims adulthood and an adult name around the age of 100 and can live to be 750 years old.",
				URL:       "/api/races/elf",
				AbilityBonuses: []race.AbilityBonus{
					{
						AbilityScore: reference.Reference{Index: "dex", Name: "DEX", URL: "/api/ability-scores/dex"},
						Bonus:        2,
					},
				},
				Languages: []reference.Reference{
					{Index: "common", Name: "Common", URL: "/api/languages/common"},
					{Index: "elvish", Name: "Elvish", URL: "/api/languages/elvish"},
				},
				Traits: []reference.Reference{
					{Index: "darkvision", Name: "Darkvision", URL: "/api/traits/darkvision"},
					{Index: "fey-ancestry", Name: "Fey Ancestry", URL: "/api/traits/fey-ancestry"},
					{Index: "trance", Name: "Trance", URL: "/api/traits/trance"},
				},
				Subraces: []reference.Reference{
					{Index: "high-elf", Name: "High Elf", URL: "/api/subraces/high-elf"},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockRace)

		case "/races/dwarf":
			mockRace := race.Race{
				Index: "dwarf",
				Name:  "Dwarf",
				Speed: 25,
				Size:  "Medium",
				URL:   "/api/races/dwarf",
				AbilityBonuses: []race.AbilityBonus{
					{
						AbilityScore: reference.Reference{Index: "con", Name: "CON"},
						Bonus:        2,
					},
				},
				Languages: []reference.Reference{
					{Index: "common", Name: "Common"},
					{Index: "dwarvish", Name: "Dwarvish"},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockRace)

		case "/races/half-elf":
			mockRace := race.Race{
				Index:           "half-elf",
				Name:            "Half-Elf",
				Speed:           30,
				Size:            "Medium",
				SizeDescription: "Half-elves are about the same size as humans, ranging from 5 to 6 feet tall.",
				AbilityBonuses: []race.AbilityBonus{
					{
						AbilityScore: reference.Reference{Index: "cha", Name: "CHA"},
						Bonus:        2,
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockRace)

		case "/races/error":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))

		case "/races/invalid-json":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"invalid json`))

		case "/races/not-found":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Race not found"))

		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}
	}))

	// Create a fetcher that uses our test server
	suite.fetcher = &core.HTTPFetcher{
		Client:  suite.server.Client(),
		BaseURL: suite.server.URL + "/",
	}
}

func (suite *FetchRaceIntegrationTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *FetchRaceIntegrationTestSuite) TestFetchRaceElfSuccess() {
	base := &templates.TemplateCharacter{Race: "elf"}
	testRace := &race.Race{}

	err := race.FetchRaceWithFetcher(suite.fetcher, base, testRace)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "elf", testRace.Index)
	assert.Equal(suite.T(), "Elf", testRace.Name)
	assert.Equal(suite.T(), 30, testRace.Speed)
	assert.Equal(suite.T(), "Medium", testRace.Size)
	assert.Len(suite.T(), testRace.AbilityBonuses, 1)
	assert.Equal(suite.T(), 2, testRace.AbilityBonuses[0].Bonus)
	assert.Equal(suite.T(), "dex", testRace.AbilityBonuses[0].AbilityScore.Index)
	assert.Len(suite.T(), testRace.Languages, 2)
	assert.Len(suite.T(), testRace.Traits, 3)
	assert.Len(suite.T(), testRace.Subraces, 1)
	assert.NotEmpty(suite.T(), testRace.Age)
}

func (suite *FetchRaceIntegrationTestSuite) TestFetchRaceDwarfSuccess() {
	base := &templates.TemplateCharacter{Race: "dwarf"}
	testRace := &race.Race{}

	err := race.FetchRaceWithFetcher(suite.fetcher, base, testRace)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "dwarf", testRace.Index)
	assert.Equal(suite.T(), "Dwarf", testRace.Name)
	assert.Equal(suite.T(), 25, testRace.Speed)
	assert.Equal(suite.T(), "con", testRace.AbilityBonuses[0].AbilityScore.Index)
	assert.Len(suite.T(), testRace.Languages, 2)
}

func (suite *FetchRaceIntegrationTestSuite) TestFetchRaceWithComplexData() {
	base := &templates.TemplateCharacter{Race: "half-elf"}
	testRace := &race.Race{}

	err := race.FetchRaceWithFetcher(suite.fetcher, base, testRace)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "half-elf", testRace.Index)
	assert.Equal(suite.T(), "Half-Elf", testRace.Name)
	assert.NotEmpty(suite.T(), testRace.SizeDescription)
	assert.Len(suite.T(), testRace.AbilityBonuses, 1)
	assert.Equal(suite.T(), "cha", testRace.AbilityBonuses[0].AbilityScore.Index)
}

func (suite *FetchRaceIntegrationTestSuite) TestFetchRaceServerError() {
	base := &templates.TemplateCharacter{Race: "error"}
	testRace := &race.Race{}

	err := race.FetchRaceWithFetcher(suite.fetcher, base, testRace)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "500")
}

func (suite *FetchRaceIntegrationTestSuite) TestFetchRaceInvalidJSON() {
	base := &templates.TemplateCharacter{Race: "invalid-json"}
	testRace := &race.Race{}

	err := race.FetchRaceWithFetcher(suite.fetcher, base, testRace)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "decode")
}

func (suite *FetchRaceIntegrationTestSuite) TestFetchRaceNotFound() {
	base := &templates.TemplateCharacter{Race: "not-found"}
	testRace := &race.Race{}

	err := race.FetchRaceWithFetcher(suite.fetcher, base, testRace)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "404")
}

func (suite *FetchRaceIntegrationTestSuite) TestFetchRaceMultipleSequential() {
	// Test multiple sequential fetches to ensure no state issues
	races := []string{"elf", "dwarf", "half-elf"}

	for _, raceName := range races {
		base := &templates.TemplateCharacter{Race: raceName}
		testRace := &race.Race{}

		err := race.FetchRaceWithFetcher(suite.fetcher, base, testRace)

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), raceName, testRace.Index)
	}
}

func TestFetchRaceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(FetchRaceIntegrationTestSuite))
}

// ========== Tests for Default FetchRace (uses real API) ==========

func TestFetchRaceWithRealAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test with real API in short mode")
	}

	base := &templates.TemplateCharacter{Race: "elf"}
	testRace := &race.Race{}

	err := race.FetchRace(base, testRace)

	// Only check if call succeeded, don't assert specific data in case API changes
	if err == nil {
		assert.NotEmpty(t, testRace.Index)
		assert.NotEmpty(t, testRace.Name)
	} else {
		t.Logf("Real API call failed (this is okay in tests): %v", err)
	}
}

// ========== Benchmark Tests ==========

func BenchmarkFetchRaceWithMock(b *testing.B) {
	mockFetcher := new(MockFetcher)
	mockFetcher.On("FetchJSON", mock.Anything, "elf").Return(nil)

	base := &templates.TemplateCharacter{Race: "elf"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testRace := &race.Race{}
		race.FetchRaceWithFetcher(mockFetcher, base, testRace)
	}
}

func BenchmarkFetchRaceWithHTTPServer(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockRace := race.Race{Index: "elf", Name: "Elf", Speed: 30}
		json.NewEncoder(w).Encode(mockRace)
	}))
	defer server.Close()

	fetcher := &core.HTTPFetcher{
		Client:  server.Client(),
		BaseURL: server.URL + "/",
	}
	base := &templates.TemplateCharacter{Race: "elf"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testRace := &race.Race{}
		race.FetchRaceWithFetcher(fetcher, base, testRace)
	}
}

// ========== Example Tests ==========

func ExampleFetchRace() {
	base := &templates.TemplateCharacter{
		Race: "elf",
	}
	testRace := &race.Race{}

	err := race.FetchRace(base, testRace)
	if err != nil {
		// Handle error
		return
	}

	// Use the fetched race
	testRace.Print()
	// Output: Race: Elf
}

func ExampleFetchRaceWithFetcher() {
	// Create a custom fetcher (e.g., for testing)
	fetcher := core.NewFetcher()

	base := &templates.TemplateCharacter{
		Race: "dwarf",
	}
	testRace := &race.Race{}

	err := race.FetchRaceWithFetcher(fetcher, base, testRace)
	if err != nil {
		// Handle error
		return
	}

	testRace.Print()
}
