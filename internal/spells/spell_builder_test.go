package spells_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/internal/reference"
	"github.com/kwford18/MKDIRagons/internal/spells"
	"github.com/kwford18/MKDIRagons/template"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// loadFixture loads a JSON fixture file.
// Uses t.Errorf to report errors safely from within goroutines.
func loadFixture(t *testing.T, filename string, target interface{}) {
	t.Helper()

	// Adjust path as needed based on where you run 'go test'
	fixtureDir := filepath.Join("testdata", "fixtures")
	filePath := filepath.Join(fixtureDir, filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read fixture file %s: %v", filePath, err)
		return
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		t.Errorf("Failed to unmarshal fixture file %s: %v", filePath, err)
	}
}

// ============================================================================
// MOCK FETCHERS
// ============================================================================

// MockFetcher - simple mock for behavior testing (verifying calls/errors)
type MockFetcher struct {
	mock.Mock
	mu sync.Mutex
}

func (m *MockFetcher) FetchJSON(property reference.Fetchable, input string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(property, input)
	return args.Error(0)
}

// MockFetcherWithFixtures - mock that loads real fixture data from disk
// bypassing the HTTP layer entirely for fast unit tests.
type MockFetcherWithFixtures struct {
	mock.Mock
	t  *testing.T
	mu sync.Mutex
}

func NewMockFetcherWithFixtures(t *testing.T) *MockFetcherWithFixtures {
	return &MockFetcherWithFixtures{t: t}
}

func (m *MockFetcherWithFixtures) FetchJSON(property reference.Fetchable, input string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	args := m.Called(property, input)

	// If the test expects success (nil error), load the actual fixture data
	if args.Error(0) == nil && input != "" {
		fixtureFile := input + ".json"
		loadFixture(m.t, fixtureFile, property)
	}

	return args.Error(0)
}

// ============================================================================
// UNIT TEST SUITE (Logic + Fixtures)
// ============================================================================

type FetchSpellsUnitTestSuite struct {
	suite.Suite
	mockFetcher         *MockFetcher
	fixtureBasedFetcher *MockFetcherWithFixtures
	base                *template.Character
	spellbook           [][]spells.Spell
}

func (suite *FetchSpellsUnitTestSuite) SetupTest() {
	suite.mockFetcher = new(MockFetcher)
	suite.fixtureBasedFetcher = NewMockFetcherWithFixtures(suite.T())

	suite.base = &template.Character{
		Spells: template.Spells{
			Level: [][]string{
				{"fire-bolt", "mage-hand"},  // Level 0
				{"magic-missile", "shield"}, // Level 1
				{"misty-step"},              // Level 2
			},
		},
	}
	suite.spellbook = spells.InitSpellbook(suite.base)
}

func TestFetchSpellsUnitTestSuite(t *testing.T) {
	suite.Run(t, new(FetchSpellsUnitTestSuite))
}

// --- Behavior Tests (MockFetcher) ---

func (suite *FetchSpellsUnitTestSuite) TestInitSpellbook() {
	assert.Len(suite.T(), suite.spellbook, 3)
	assert.Len(suite.T(), suite.spellbook[0], 2)
	assert.Len(suite.T(), suite.spellbook[1], 2)
	assert.Len(suite.T(), suite.spellbook[2], 1)
}

func (suite *FetchSpellsUnitTestSuite) TestFetchSpells_PartialFailure() {
	// Simulate one spell failing and another succeeding
	suite.mockFetcher.On("FetchJSON", mock.Anything, "fire-bolt").Return(errors.New("network error"))
	suite.mockFetcher.On("FetchJSON", mock.Anything, mock.Anything).Return(nil).Maybe()

	err := spells.FetchSpellsWithFetcher(suite.mockFetcher, suite.base, suite.spellbook)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "network error")
}

func (suite *FetchSpellsUnitTestSuite) TestFetchSpells_NilInputs() {
	assert.Panics(suite.T(), func() {
		_ = spells.FetchSpellsWithFetcher(suite.mockFetcher, nil, suite.spellbook)
	})
	assert.Panics(suite.T(), func() {
		_ = spells.FetchSpellsWithFetcher(suite.mockFetcher, suite.base, nil)
	})
}

// --- Data Tests (MockFetcherWithFixtures) ---

func (suite *FetchSpellsUnitTestSuite) TestFetchSpells_FireballWithFixture() {
	// Setup character with just Fireball (Level 3)
	base := &template.Character{
		Spells: template.Spells{
			Level: [][]string{
				{}, {}, {}, {"fireball"},
			},
		},
	}
	spellbook := spells.InitSpellbook(base)

	// Expectation: FetchJSON called for "fireball", returns nil (Success)
	// The MockFetcherWithFixtures will automatically load "fireball.json" into the struct
	suite.fixtureBasedFetcher.On("FetchJSON", mock.AnythingOfType("*spells.Spell"), "fireball").Return(nil)

	err := spells.FetchSpellsWithFetcher(suite.fixtureBasedFetcher, base, spellbook)

	require.NoError(suite.T(), err)

	// Verify Data
	fireball := spellbook[3][0]
	assert.Equal(suite.T(), "Fireball", fireball.Name)
	assert.Equal(suite.T(), 3, fireball.Level)

	// Verify Deep Nested Data (Damage, DC, AoE)
	require.NotNil(suite.T(), fireball.Damage)
	assert.Equal(suite.T(), "fire", fireball.Damage.DamageType.Index)
	assert.Equal(suite.T(), "8d6", fireball.Damage.DamageAtSlotLevel["3"])

	require.NotNil(suite.T(), fireball.AreaOfEffect)
	assert.Equal(suite.T(), 20, fireball.AreaOfEffect.Size)
}

func (suite *FetchSpellsUnitTestSuite) TestFetchSpells_WishWithFixture() {
	// Setup character with Wish (Level 9)
	base := &template.Character{
		Spells: template.Spells{
			Level: [][]string{
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {"wish"},
			},
		},
	}
	spellbook := spells.InitSpellbook(base)

	suite.fixtureBasedFetcher.On("FetchJSON", mock.AnythingOfType("*spells.Spell"), "wish").Return(nil)

	err := spells.FetchSpellsWithFetcher(suite.fixtureBasedFetcher, base, spellbook)

	require.NoError(suite.T(), err)

	wish := spellbook[9][0]
	assert.Equal(suite.T(), "Wish", wish.Name)

	// Verify nil pointers for optional fields not present in Wish fixture
	assert.Nil(suite.T(), wish.Damage, "Wish should have nil damage")
	assert.Nil(suite.T(), wish.DC, "Wish should have nil DC")
}

// ============================================================================
// INTEGRATION TEST SUITE (HTTP Layer + Fixtures + Real API)
// ============================================================================

type FetchSpellsIntegrationTestSuite struct {
	suite.Suite
	server  *httptest.Server
	fetcher *core.HTTPFetcher
}

func (suite *FetchSpellsIntegrationTestSuite) SetupSuite() {
	// GENERIC FIXTURE SERVER
	// Instead of hardcoding JSON here, we serve the actual fixture files.
	// This ensures our Integration Tests and Unit Tests use the EXACT same data.
	suite.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract spell name from URL (e.g., "/spells/fire-bolt" -> "fire-bolt")
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 2 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		spellName := pathParts[len(pathParts)-1]

		// 2. Map to fixture file
		fixturePath := filepath.Join("testdata", "fixtures", spellName+".json")

		// 3. Serve file
		data, err := os.ReadFile(fixturePath)
		if err != nil {
			// If fixture doesn't exist, return 404
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}))

	// Point the real HTTPFetcher to our local test server
	suite.fetcher = &core.HTTPFetcher{
		Client:  suite.server.Client(),
		BaseURL: suite.server.URL + "/spells/",
	}
}

func (suite *FetchSpellsIntegrationTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *FetchSpellsIntegrationTestSuite) TestIntegration_FetchRealHTTP_Fireball() {
	// This test uses the REAL HTTPFetcher (from fetcher.go)
	// It hits the httptest server, which reads from testdata/fixtures/fireball.json

	base := &template.Character{
		Spells: template.Spells{
			Level: [][]string{
				{}, {}, {}, {"fireball"},
			},
		},
	}
	spellbook := spells.InitSpellbook(base)

	err := spells.FetchSpellsWithFetcher(suite.fetcher, base, spellbook)

	require.NoError(suite.T(), err)

	fireball := spellbook[3][0]
	assert.Equal(suite.T(), "Fireball", fireball.Name)
	assert.Equal(suite.T(), "150 feet", fireball.Range)
}

func (suite *FetchSpellsIntegrationTestSuite) TestIntegration_FetchRealHTTP_NotFound() {
	base := &template.Character{
		Spells: template.Spells{
			Level: [][]string{
				{"nonexistent-spell"},
			},
		},
	}
	spellbook := spells.InitSpellbook(base)

	err := spells.FetchSpellsWithFetcher(suite.fetcher, base, spellbook)

	// Should error because the Generic Fixture Server won't find "nonexistent-spell.json"
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "404")
}

// ============================================================================
// REAL API INTEGRATION TESTS (External Network)
// ============================================================================

func (suite *FetchSpellsIntegrationTestSuite) TestIntegration_RealAPI_Fireball() {
	// This test hits the ACTUAL 5e API (https://www.dnd5eapi.co)
	// It is skipped if `go test -short` is run.
	if testing.Short() {
		suite.T().Skip("Skipping real API test in short mode")
	}

	// Create a new fetcher that points to the production URL (DefaultBaseURL)
	realFetcher := core.NewFetcher()

	base := &template.Character{
		Spells: template.Spells{
			Level: [][]string{
				{}, {}, {}, {"fireball"},
			},
		},
	}
	spellbook := spells.InitSpellbook(base)

	// Perform fetch
	err := spells.FetchSpellsWithFetcher(realFetcher, base, spellbook)
	require.NoError(suite.T(), err, "Real API fetch failed - check network connection")

	// Verify Data
	fireball := spellbook[3][0]
	assert.Equal(suite.T(), "Fireball", fireball.Name)
	assert.Equal(suite.T(), "fireball", fireball.Index)
	// Real API should always have these fields populated for Fireball
	require.NotNil(suite.T(), fireball.Damage)
	assert.Equal(suite.T(), "fire", fireball.Damage.DamageType.Index)
}

func (suite *FetchSpellsIntegrationTestSuite) TestIntegration_RealAPI_Multiple() {
	if testing.Short() {
		suite.T().Skip("Skipping real API test in short mode")
	}

	realFetcher := core.NewFetcher()

	base := &template.Character{
		Spells: template.Spells{
			Level: [][]string{
				{"mage-hand"}, // Level 0
				{"shield"},    // Level 1
			},
		},
	}
	spellbook := spells.InitSpellbook(base)

	err := spells.FetchSpellsWithFetcher(realFetcher, base, spellbook)
	require.NoError(suite.T(), err)

	mageHand := spellbook[0][0]
	assert.Equal(suite.T(), "Mage Hand", mageHand.Name)

	shield := spellbook[1][0]
	assert.Equal(suite.T(), "Shield", shield.Name)
}

func TestFetchSpellsIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(FetchSpellsIntegrationTestSuite))
}
