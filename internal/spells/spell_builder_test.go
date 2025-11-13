package spells_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/internal/reference"
	"github.com/kwford18/MKDIRagons/internal/spells"
	"github.com/kwford18/MKDIRagons/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// ========== Mock Fetcher ==========

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

// ========== Unit Tests with Mock Fetcher ==========

type FetchSpellsUnitTestSuite struct {
	suite.Suite
	mockFetcher *MockFetcher
	base        *templates.TemplateCharacter
	spellbook   [][]spells.Spell
}

func (suite *FetchSpellsUnitTestSuite) SetupTest() {
	suite.mockFetcher = new(MockFetcher)
	suite.base = &templates.TemplateCharacter{
		Spells: templates.TemplateSpells{
			Level: [][]string{
				{"fire-bolt", "mage-hand"},  // Level 0 (cantrips)
				{"magic-missile", "shield"}, // Level 1
				{"misty-step"},              // Level 2
			},
		},
	}
	suite.spellbook = spells.InitSpellbook(suite.base)
}

func (suite *FetchSpellsUnitTestSuite) TestInitSpellbook() {
	assert.Len(suite.T(), suite.spellbook, 3)
	assert.Len(suite.T(), suite.spellbook[0], 2) // 2 cantrips
	assert.Len(suite.T(), suite.spellbook[1], 2) // 2 level 1 spells
	assert.Len(suite.T(), suite.spellbook[2], 1) // 1 level 2 spell
}

func (suite *FetchSpellsUnitTestSuite) TestFetchSpellsSuccess() {
	// Mock cantrips
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*spells.Spell"), "fire-bolt").Return(nil).Run(func(args mock.Arguments) {
		spell := args.Get(0).(*spells.Spell)
		spell.Index = "fire-bolt"
		spell.Name = "Fire Bolt"
		spell.Level = 0
	})
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*spells.Spell"), "mage-hand").Return(nil).Run(func(args mock.Arguments) {
		spell := args.Get(0).(*spells.Spell)
		spell.Index = "mage-hand"
		spell.Name = "Mage Hand"
		spell.Level = 0
	})

	// Mock level 1 spells
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*spells.Spell"), "magic-missile").Return(nil).Run(func(args mock.Arguments) {
		spell := args.Get(0).(*spells.Spell)
		spell.Index = "magic-missile"
		spell.Name = "Magic Missile"
		spell.Level = 1
	})
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*spells.Spell"), "shield").Return(nil).Run(func(args mock.Arguments) {
		spell := args.Get(0).(*spells.Spell)
		spell.Index = "shield"
		spell.Name = "Shield"
		spell.Level = 1
	})

	// Mock level 2 spells
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*spells.Spell"), "misty-step").Return(nil).Run(func(args mock.Arguments) {
		spell := args.Get(0).(*spells.Spell)
		spell.Index = "misty-step"
		spell.Name = "Misty Step"
		spell.Level = 2
	})

	err := spells.FetchSpellsWithFetcher(suite.mockFetcher, suite.base, suite.spellbook)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "fire-bolt", suite.spellbook[0][0].Index)
	assert.Equal(suite.T(), "mage-hand", suite.spellbook[0][1].Index)
	assert.Equal(suite.T(), "magic-missile", suite.spellbook[1][0].Index)
	assert.Equal(suite.T(), "shield", suite.spellbook[1][1].Index)
	assert.Equal(suite.T(), "misty-step", suite.spellbook[2][0].Index)
	suite.mockFetcher.AssertExpectations(suite.T())
}

func (suite *FetchSpellsUnitTestSuite) TestFetchSpellsError() {
	// One spell fetch fails
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*spells.Spell"), "fire-bolt").Return(errors.New("spell fetch failed"))
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*spells.Spell"), mock.Anything).Return(nil).Maybe()

	err := spells.FetchSpellsWithFetcher(suite.mockFetcher, suite.base, suite.spellbook)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "spell fetch failed")
}

func (suite *FetchSpellsUnitTestSuite) TestFetchSpellsEmptySpellbook() {
	suite.base.Spells = templates.TemplateSpells{
		Level: [][]string{},
	}
	suite.spellbook = spells.InitSpellbook(suite.base)

	err := spells.FetchSpellsWithFetcher(suite.mockFetcher, suite.base, suite.spellbook)

	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), suite.spellbook)
}

func (suite *FetchSpellsUnitTestSuite) TestFetchSpellsOnlyCantrips() {
	suite.base.Spells = templates.TemplateSpells{
		Level: [][]string{
			{"fire-bolt", "mage-hand"},
		},
	}
	suite.spellbook = spells.InitSpellbook(suite.base)

	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*spells.Spell"), "fire-bolt").Return(nil).Run(func(args mock.Arguments) {
		spell := args.Get(0).(*spells.Spell)
		spell.Index = "fire-bolt"
	})
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*spells.Spell"), "mage-hand").Return(nil).Run(func(args mock.Arguments) {
		spell := args.Get(0).(*spells.Spell)
		spell.Index = "mage-hand"
	})

	err := spells.FetchSpellsWithFetcher(suite.mockFetcher, suite.base, suite.spellbook)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), suite.spellbook, 1)
	assert.Len(suite.T(), suite.spellbook[0], 2)
}

func (suite *FetchSpellsUnitTestSuite) TestFetchSpellsNilBase() {
	err := spells.FetchSpellsWithFetcher(suite.mockFetcher, nil, suite.spellbook)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "base TemplateCharacter")
}

func (suite *FetchSpellsUnitTestSuite) TestFetchSpellsNilSpellbook() {
	err := spells.FetchSpellsWithFetcher(suite.mockFetcher, suite.base, nil)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "spellbook")
}

func TestFetchSpellsUnitTestSuite(t *testing.T) {
	suite.Run(t, new(FetchSpellsUnitTestSuite))
}

// ========== Integration Tests with HTTP Server ==========

type FetchSpellsIntegrationTestSuite struct {
	suite.Suite
	server  *httptest.Server
	fetcher *core.HTTPFetcher
}

func (suite *FetchSpellsIntegrationTestSuite) SetupSuite() {
	suite.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		// Cantrips (Level 0)
		case "/spells/fire-bolt":
			err := json.NewEncoder(w).Encode(spells.Spell{
				Index:       "fire-bolt",
				Name:        "Fire Bolt",
				Level:       0,
				School:      reference.Reference{Index: "evocation", Name: "Evocation"},
				CastingTime: "1 action",
				Range:       "120 feet",
			})
			if err != nil {
				return
			}
		case "/spells/mage-hand":
			err := json.NewEncoder(w).Encode(spells.Spell{
				Index:       "mage-hand",
				Name:        "Mage Hand",
				Level:       0,
				School:      reference.Reference{Index: "conjuration", Name: "Conjuration"},
				CastingTime: "1 action",
			})
			if err != nil {
				return
			}
		case "/spells/prestidigitation":
			err := json.NewEncoder(w).Encode(spells.Spell{
				Index: "prestidigitation",
				Name:  "Prestidigitation",
				Level: 0,
			})
			if err != nil {
				return
			}

		// Level 1 Spells
		case "/spells/magic-missile":
			err := json.NewEncoder(w).Encode(spells.Spell{
				Index:       "magic-missile",
				Name:        "Magic Missile",
				Level:       1,
				School:      reference.Reference{Index: "evocation", Name: "Evocation"},
				CastingTime: "1 action",
				Range:       "120 feet",
			})
			if err != nil {
				return
			}
		case "/spells/shield":
			err := json.NewEncoder(w).Encode(spells.Spell{
				Index:       "shield",
				Name:        "Shield",
				Level:       1,
				School:      reference.Reference{Index: "abjuration", Name: "Abjuration"},
				CastingTime: "1 reaction",
			})
			if err != nil {
				return
			}
		case "/spells/detect-magic":
			err := json.NewEncoder(w).Encode(spells.Spell{
				Index: "detect-magic",
				Name:  "Detect Magic",
				Level: 1,
			})
			if err != nil {
				return
			}

		// Level 2 Spells
		case "/spells/misty-step":
			err := json.NewEncoder(w).Encode(spells.Spell{
				Index:       "misty-step",
				Name:        "Misty Step",
				Level:       2,
				School:      reference.Reference{Index: "conjuration", Name: "Conjuration"},
				CastingTime: "1 bonus action",
			})
			if err != nil {
				return
			}
		case "/spells/scorching-ray":
			err := json.NewEncoder(w).Encode(spells.Spell{
				Index: "scorching-ray",
				Name:  "Scorching Ray",
				Level: 2,
			})
			if err != nil {
				return
			}

		// Level 3 Spells
		case "/spells/fireball":
			err := json.NewEncoder(w).Encode(spells.Spell{
				Index:       "fireball",
				Name:        "Fireball",
				Level:       3,
				School:      reference.Reference{Index: "evocation", Name: "Evocation"},
				CastingTime: "1 action",
				Range:       "150 feet",
			})
			if err != nil {
				return
			}

		// Error endpoints
		case "/spells/error":
			w.WriteHeader(http.StatusInternalServerError)
		case "/spells/invalid":
			_, err := w.Write([]byte(`{"invalid json`))
			if err != nil {
				return
			}

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	suite.fetcher = &core.HTTPFetcher{
		Client:  suite.server.Client(),
		BaseURL: suite.server.URL + "/",
	}
}

func (suite *FetchSpellsIntegrationTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *FetchSpellsIntegrationTestSuite) TestFetchSpellsComplete() {
	base := &templates.TemplateCharacter{
		Spells: templates.TemplateSpells{
			Level: [][]string{
				{"fire-bolt", "mage-hand"},
				{"magic-missile", "shield"},
				{"misty-step"},
			},
		},
	}
	spellbook := spells.InitSpellbook(base)

	err := spells.FetchSpellsWithFetcher(suite.fetcher, base, spellbook)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), spellbook, 3)

	// Verify cantrips
	assert.Equal(suite.T(), "fire-bolt", spellbook[0][0].Index)
	assert.Equal(suite.T(), "Fire Bolt", spellbook[0][0].Name)
	assert.Equal(suite.T(), 0, spellbook[0][0].Level)
	assert.Equal(suite.T(), "mage-hand", spellbook[0][1].Index)

	// Verify level 1 spells
	assert.Equal(suite.T(), "magic-missile", spellbook[1][0].Index)
	assert.Equal(suite.T(), 1, spellbook[1][0].Level)
	assert.Equal(suite.T(), "shield", spellbook[1][1].Index)

	// Verify level 2 spells
	assert.Equal(suite.T(), "misty-step", spellbook[2][0].Index)
	assert.Equal(suite.T(), 2, spellbook[2][0].Level)
}

func (suite *FetchSpellsIntegrationTestSuite) TestFetchSpellsLargeSpellbook() {
	base := &templates.TemplateCharacter{
		Spells: templates.TemplateSpells{
			Level: [][]string{
				{"fire-bolt", "mage-hand", "prestidigitation"},
				{"magic-missile", "shield", "detect-magic"},
				{"misty-step", "scorching-ray"},
				{"fireball"},
			},
		},
	}
	spellbook := spells.InitSpellbook(base)

	err := spells.FetchSpellsWithFetcher(suite.fetcher, base, spellbook)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), spellbook, 4)
	assert.Len(suite.T(), spellbook[0], 3) // 3 cantrips
	assert.Len(suite.T(), spellbook[1], 3) // 3 level 1 spells
	assert.Len(suite.T(), spellbook[2], 2) // 2 level 2 spells
	assert.Len(suite.T(), spellbook[3], 1) // 1 level 3 spell
}

func (suite *FetchSpellsIntegrationTestSuite) TestFetchSpellsOnlyCantrips() {
	base := &templates.TemplateCharacter{
		Spells: templates.TemplateSpells{
			Level: [][]string{
				{"fire-bolt", "mage-hand"},
			},
		},
	}
	spellbook := spells.InitSpellbook(base)

	err := spells.FetchSpellsWithFetcher(suite.fetcher, base, spellbook)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), spellbook, 1)
	assert.Len(suite.T(), spellbook[0], 2)
}

func (suite *FetchSpellsIntegrationTestSuite) TestFetchSpellsServerError() {
	base := &templates.TemplateCharacter{
		Spells: templates.TemplateSpells{
			Level: [][]string{
				{"error"},
			},
		},
	}
	spellbook := spells.InitSpellbook(base)

	err := spells.FetchSpellsWithFetcher(suite.fetcher, base, spellbook)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "500")
}

func (suite *FetchSpellsIntegrationTestSuite) TestFetchSpellsInvalidJSON() {
	base := &templates.TemplateCharacter{
		Spells: templates.TemplateSpells{
			Level: [][]string{
				{"invalid"},
			},
		},
	}
	spellbook := spells.InitSpellbook(base)

	err := spells.FetchSpellsWithFetcher(suite.fetcher, base, spellbook)

	assert.Error(suite.T(), err)
}

func (suite *FetchSpellsIntegrationTestSuite) TestFetchSpellsNotFound() {
	base := &templates.TemplateCharacter{
		Spells: templates.TemplateSpells{
			Level: [][]string{
				{"nonexistent-spell"},
			},
		},
	}
	spellbook := spells.InitSpellbook(base)

	err := spells.FetchSpellsWithFetcher(suite.fetcher, base, spellbook)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "404")
}

func (suite *FetchSpellsIntegrationTestSuite) TestFetchSpellsConcurrency() {
	// Test that concurrent fetches don't cause race conditions
	base := &templates.TemplateCharacter{
		Spells: templates.TemplateSpells{
			Level: [][]string{
				{"fire-bolt", "mage-hand", "prestidigitation"},
				{"magic-missile", "shield", "detect-magic"},
				{"misty-step", "scorching-ray"},
			},
		},
	}

	// Run multiple times to increase chance of catching race conditions
	for i := 0; i < 10; i++ {
		spellbook := spells.InitSpellbook(base)
		err := spells.FetchSpellsWithFetcher(suite.fetcher, base, spellbook)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), spellbook, 3)
		assert.Len(suite.T(), spellbook[0], 3)
		assert.Len(suite.T(), spellbook[1], 3)
		assert.Len(suite.T(), spellbook[2], 2)
	}
}

func TestFetchSpellsIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(FetchSpellsIntegrationTestSuite))
}

// ========== Benchmark Tests ==========

func BenchmarkFetchSpellsSmall(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(spells.Spell{Index: "fire-bolt", Name: "Fire Bolt", Level: 0})
		if err != nil {
			return
		}
	}))
	defer server.Close()

	fetcher := &core.HTTPFetcher{
		Client:  server.Client(),
		BaseURL: server.URL + "/",
	}

	base := &templates.TemplateCharacter{
		Spells: templates.TemplateSpells{
			Level: [][]string{
				{"fire-bolt"},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spellbook := spells.InitSpellbook(base)
		err := spells.FetchSpellsWithFetcher(fetcher, base, spellbook)
		if err != nil {
			return
		}
	}
}

func BenchmarkFetchSpellsLarge(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(spells.Spell{Index: "spell", Name: "Spell", Level: 0})
		if err != nil {
			return
		}
	}))
	defer server.Close()

	fetcher := &core.HTTPFetcher{
		Client:  server.Client(),
		BaseURL: server.URL + "/",
	}

	base := &templates.TemplateCharacter{
		Spells: templates.TemplateSpells{
			Level: [][]string{
				{"s1", "s2", "s3", "s4"},
				{"s5", "s6", "s7"},
				{"s8", "s9"},
				{"s10"},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spellbook := spells.InitSpellbook(base)
		err := spells.FetchSpellsWithFetcher(fetcher, base, spellbook)
		if err != nil {
			return
		}
	}
}

// ========== Example Tests ==========

func ExampleInitSpellbook() {
	base := &templates.TemplateCharacter{
		Spells: templates.TemplateSpells{
			Level: [][]string{
				{"fire-bolt", "mage-hand"},
				{"magic-missile"},
			},
		},
	}

	spellbook := spells.InitSpellbook(base)

	// Spellbook is now a 2D array ready to be filled
	_ = len(spellbook)    // 2 spell levels
	_ = len(spellbook[0]) // 2 cantrips
	_ = len(spellbook[1]) // 1 level 1 spell
}

func ExampleFetchSpells() {
	base := &templates.TemplateCharacter{
		Spells: templates.TemplateSpells{
			Level: [][]string{
				{"fire-bolt", "mage-hand"},
				{"magic-missile"},
			},
		},
	}
	spellbook := spells.InitSpellbook(base)

	err := spells.FetchSpells(base, spellbook)
	if err != nil {
		// Handle error
		return
	}

	// Use the fetched spells
	for level, spellsAtLevel := range spellbook {
		for _, spell := range spellsAtLevel {
			_ = spell.Name
			_ = level
		}
	}
}
