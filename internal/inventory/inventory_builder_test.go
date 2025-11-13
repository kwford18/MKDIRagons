package inventory_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/internal/inventory"
	"github.com/kwford18/MKDIRagons/internal/reference"
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

type FetchInventoryUnitTestSuite struct {
	suite.Suite
	mockFetcher *MockFetcher
	base        *templates.TemplateCharacter
	inv         *inventory.Inventory
}

func (suite *FetchInventoryUnitTestSuite) SetupTest() {
	suite.mockFetcher = new(MockFetcher)
	suite.base = &templates.TemplateCharacter{
		Inventory: templates.TemplateInventory{
			Armor:   []string{"leather-armor", "shield"},
			Weapons: []string{"longsword", "shortbow"},
			Items:   []string{"rope", "torch"},
		},
	}
	suite.inv = &inventory.Inventory{}
}

func (suite *FetchInventoryUnitTestSuite) TestFetchInventorySuccess() {
	// Mock all armor fetches
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), "leather-armor").Return(nil).Run(func(args mock.Arguments) {
		armor := args.Get(0).(*inventory.Armor)
		armor.Index = "leather-armor"
		armor.Name = "Leather Armor"
	})
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), "shield").Return(nil).Run(func(args mock.Arguments) {
		armor := args.Get(0).(*inventory.Armor)
		armor.Index = "shield"
		armor.Name = "Shield"
	})

	// Mock all weapon fetches
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Weapon"), "longsword").Return(nil).Run(func(args mock.Arguments) {
		weapon := args.Get(0).(*inventory.Weapon)
		weapon.Index = "longsword"
		weapon.Name = "Longsword"
	})
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Weapon"), "shortbow").Return(nil).Run(func(args mock.Arguments) {
		weapon := args.Get(0).(*inventory.Weapon)
		weapon.Index = "shortbow"
		weapon.Name = "Shortbow"
	})

	// Mock all item fetches
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Item"), "rope").Return(nil).Run(func(args mock.Arguments) {
		item := args.Get(0).(*inventory.Item)
		item.Index = "rope"
		item.Name = "Rope"
	})
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Item"), "torch").Return(nil).Run(func(args mock.Arguments) {
		item := args.Get(0).(*inventory.Item)
		item.Index = "torch"
		item.Name = "Torch"
	})

	err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, suite.base, suite.inv)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), suite.inv.Armor, 2)
	assert.Len(suite.T(), suite.inv.Weapons, 2)
	assert.Len(suite.T(), suite.inv.Items, 2)
	suite.mockFetcher.AssertExpectations(suite.T())
}

func (suite *FetchInventoryUnitTestSuite) TestFetchInventoryArmorError() {
	// One armor fetch fails
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), "leather-armor").Return(errors.New("armor fetch failed"))
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), "shield").Return(nil)
	suite.mockFetcher.On("FetchJSON", mock.Anything, mock.Anything).Return(nil).Maybe()

	err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, suite.base, suite.inv)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "armor fetch failed")
}

func (suite *FetchInventoryUnitTestSuite) TestFetchInventoryWeaponError() {
	// Setup armor to succeed
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), mock.Anything).Return(nil).Maybe()

	// One weapon fetch fails
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Weapon"), "longsword").Return(errors.New("weapon fetch failed"))
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Weapon"), "shortbow").Return(nil).Maybe()
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Item"), mock.Anything).Return(nil).Maybe()

	err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, suite.base, suite.inv)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "weapon fetch failed")
}

func (suite *FetchInventoryUnitTestSuite) TestFetchInventoryEmptyInventory() {
	suite.base.Inventory = templates.TemplateInventory{
		Armor:   []string{},
		Weapons: []string{},
		Items:   []string{},
	}

	err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, suite.base, suite.inv)

	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), suite.inv.Armor)
	assert.Empty(suite.T(), suite.inv.Weapons)
	assert.Empty(suite.T(), suite.inv.Items)
}

func (suite *FetchInventoryUnitTestSuite) TestFetchInventoryOnlyArmor() {
	suite.base.Inventory = templates.TemplateInventory{
		Armor:   []string{"leather-armor"},
		Weapons: []string{},
		Items:   []string{},
	}

	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), "leather-armor").Return(nil).Run(func(args mock.Arguments) {
		armor := args.Get(0).(*inventory.Armor)
		armor.Index = "leather-armor"
	})

	err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, suite.base, suite.inv)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), suite.inv.Armor, 1)
	assert.Empty(suite.T(), suite.inv.Weapons)
	assert.Empty(suite.T(), suite.inv.Items)
}

func (suite *FetchInventoryUnitTestSuite) TestFetchInventoryNilBase() {
	assert.Panics(suite.T(), func() {
		err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, nil, suite.inv)
		if err != nil {
			return
		}
	})
}

func (suite *FetchInventoryUnitTestSuite) TestFetchInventoryNilInventory() {
	assert.Panics(suite.T(), func() {
		err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, suite.base, nil)
		if err != nil {
			return
		}
	})
}

func TestFetchInventoryUnitTestSuite(t *testing.T) {
	suite.Run(t, new(FetchInventoryUnitTestSuite))
}

// ========== Integration Tests with HTTP Server ==========

type FetchInventoryIntegrationTestSuite struct {
	suite.Suite
	server  *httptest.Server
	fetcher *core.HTTPFetcher
}

func (suite *FetchInventoryIntegrationTestSuite) SetupSuite() {
	suite.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		// Armor endpoints
		case "/armor/leather-armor":
			err := json.NewEncoder(w).Encode(inventory.Armor{
				Index:      "leather-armor",
				Name:       "Leather Armor",
				ArmorClass: inventory.ArmorClass{Base: 11, DexBonus: true},
			})
			if err != nil {
				return
			}
		case "/armor/shield":
			err := json.NewEncoder(w).Encode(inventory.Armor{
				Index:      "shield",
				Name:       "Shield",
				ArmorClass: inventory.ArmorClass{Base: 2},
			})
			if err != nil {
				return
			}
		case "/armor/chain-mail":
			err := json.NewEncoder(w).Encode(inventory.Armor{
				Index:      "chain-mail",
				Name:       "Chain Mail",
				ArmorClass: inventory.ArmorClass{Base: 16},
			})
			if err != nil {
				return
			}

		// Weapon endpoints
		case "/weapons/longsword":
			err := json.NewEncoder(w).Encode(inventory.Weapon{
				Index: "longsword",
				Name:  "Longsword",
				Damage: inventory.Damage{
					DamageDice: "1d8",
					DamageType: reference.Reference{Index: "slashing", Name: "Slashing"},
				},
			})
			if err != nil {
				return
			}
		case "/weapons/shortbow":
			err := json.NewEncoder(w).Encode(inventory.Weapon{
				Index: "shortbow",
				Name:  "Shortbow",
				Damage: inventory.Damage{
					DamageDice: "1d6",
					DamageType: reference.Reference{Index: "piercing", Name: "Piercing"},
				},
			})
			if err != nil {
				return
			}
		case "/weapons/dagger":
			err := json.NewEncoder(w).Encode(inventory.Weapon{
				Index: "dagger",
				Name:  "Dagger",
			})
			if err != nil {
				return
			}

		// Item endpoints
		case "/equipment/rope":
			err := json.NewEncoder(w).Encode(inventory.Item{
				Index: "rope",
				Name:  "Rope, Hempen (50 feet)",
			})
			if err != nil {
				return
			}
		case "/equipment/torch":
			err := json.NewEncoder(w).Encode(inventory.Item{
				Index: "torch",
				Name:  "Torch",
			})
			if err != nil {
				return
			}
		case "/equipment/backpack":
			err := json.NewEncoder(w).Encode(inventory.Item{
				Index: "backpack",
				Name:  "Backpack",
			})
			if err != nil {
				return
			}

		// Error endpoints
		case "/armor/error", "/weapons/error", "/equipment/error":
			w.WriteHeader(http.StatusInternalServerError)
		case "/armor/invalid", "/weapons/invalid", "/equipment/invalid":
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

func (suite *FetchInventoryIntegrationTestSuite) TearDownSuite() {
	suite.server.Close()
}

func (suite *FetchInventoryIntegrationTestSuite) TestFetchInventoryComplete() {
	base := &templates.TemplateCharacter{
		Inventory: templates.TemplateInventory{
			Armor:   []string{"leather-armor", "shield"},
			Weapons: []string{"longsword", "shortbow"},
			Items:   []string{"rope", "torch"},
		},
	}
	inv := &inventory.Inventory{}

	err := inventory.FetchInventoryWithFetcher(suite.fetcher, base, inv)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), inv.Armor, 2)
	assert.Len(suite.T(), inv.Weapons, 2)
	assert.Len(suite.T(), inv.Items, 2)

	// Verify armor
	assert.Equal(suite.T(), "leather-armor", inv.Armor[0].Index)
	assert.Equal(suite.T(), "shield", inv.Armor[1].Index)

	// Verify weapons
	assert.Equal(suite.T(), "longsword", inv.Weapons[0].Index)
	assert.Equal(suite.T(), "shortbow", inv.Weapons[1].Index)

	// Verify items
	assert.Equal(suite.T(), "rope", inv.Items[0].Index)
	assert.Equal(suite.T(), "torch", inv.Items[1].Index)
}

func (suite *FetchInventoryIntegrationTestSuite) TestFetchInventoryLargeInventory() {
	base := &templates.TemplateCharacter{
		Inventory: templates.TemplateInventory{
			Armor:   []string{"leather-armor", "shield", "chain-mail"},
			Weapons: []string{"longsword", "shortbow", "dagger"},
			Items:   []string{"rope", "torch", "backpack"},
		},
	}
	inv := &inventory.Inventory{}

	err := inventory.FetchInventoryWithFetcher(suite.fetcher, base, inv)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), inv.Armor, 3)
	assert.Len(suite.T(), inv.Weapons, 3)
	assert.Len(suite.T(), inv.Items, 3)
}

func (suite *FetchInventoryIntegrationTestSuite) TestFetchInventoryOnlyWeapons() {
	base := &templates.TemplateCharacter{
		Inventory: templates.TemplateInventory{
			Armor:   []string{},
			Weapons: []string{"longsword", "dagger"},
			Items:   []string{},
		},
	}
	inv := &inventory.Inventory{}

	err := inventory.FetchInventoryWithFetcher(suite.fetcher, base, inv)

	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), inv.Armor)
	assert.Len(suite.T(), inv.Weapons, 2)
	assert.Empty(suite.T(), inv.Items)
}

func (suite *FetchInventoryIntegrationTestSuite) TestFetchInventoryServerError() {
	base := &templates.TemplateCharacter{
		Inventory: templates.TemplateInventory{
			Armor:   []string{"error"},
			Weapons: []string{},
			Items:   []string{},
		},
	}
	inv := &inventory.Inventory{}

	err := inventory.FetchInventoryWithFetcher(suite.fetcher, base, inv)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "500")
}

func (suite *FetchInventoryIntegrationTestSuite) TestFetchInventoryInvalidJSON() {
	base := &templates.TemplateCharacter{
		Inventory: templates.TemplateInventory{
			Armor:   []string{},
			Weapons: []string{"invalid"},
			Items:   []string{},
		},
	}
	inv := &inventory.Inventory{}

	err := inventory.FetchInventoryWithFetcher(suite.fetcher, base, inv)

	assert.Error(suite.T(), err)
}

func (suite *FetchInventoryIntegrationTestSuite) TestFetchInventoryNotFound() {
	base := &templates.TemplateCharacter{
		Inventory: templates.TemplateInventory{
			Armor:   []string{},
			Weapons: []string{},
			Items:   []string{"nonexistent-item"},
		},
	}
	inv := &inventory.Inventory{}

	err := inventory.FetchInventoryWithFetcher(suite.fetcher, base, inv)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "404")
}

func (suite *FetchInventoryIntegrationTestSuite) TestFetchInventoryConcurrency() {
	// Test that concurrent fetches don't cause race conditions
	base := &templates.TemplateCharacter{
		Inventory: templates.TemplateInventory{
			Armor:   []string{"leather-armor", "shield", "chain-mail"},
			Weapons: []string{"longsword", "shortbow", "dagger"},
			Items:   []string{"rope", "torch", "backpack"},
		},
	}

	// Run multiple times to increase chance of catching race conditions
	for i := 0; i < 10; i++ {
		inv := &inventory.Inventory{}
		err := inventory.FetchInventoryWithFetcher(suite.fetcher, base, inv)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), inv.Armor, 3)
		assert.Len(suite.T(), inv.Weapons, 3)
		assert.Len(suite.T(), inv.Items, 3)
	}
}

func TestFetchInventoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(FetchInventoryIntegrationTestSuite))
}

// ========== Benchmark Tests ==========

func BenchmarkFetchInventorySmall(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(inventory.Armor{Index: "leather-armor", Name: "Leather Armor"})
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
		Inventory: templates.TemplateInventory{
			Armor: []string{"leather-armor"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inv := &inventory.Inventory{}
		err := inventory.FetchInventoryWithFetcher(fetcher, base, inv)
		if err != nil {
			return
		}
	}
}

func BenchmarkFetchInventoryLarge(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(inventory.Item{Index: "item", Name: "Item"})
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
		Inventory: templates.TemplateInventory{
			Armor:   []string{"a1", "a2", "a3"},
			Weapons: []string{"w1", "w2", "w3"},
			Items:   []string{"i1", "i2", "i3", "i4", "i5"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inv := &inventory.Inventory{}
		err := inventory.FetchInventoryWithFetcher(fetcher, base, inv)
		if err != nil {
			return
		}
	}
}

// ========== Example Tests ==========

func ExampleFetchInventory() {
	base := &templates.TemplateCharacter{
		Inventory: templates.TemplateInventory{
			Armor:   []string{"leather-armor"},
			Weapons: []string{"longsword"},
			Items:   []string{"rope"},
		},
	}
	inv := &inventory.Inventory{}

	err := inventory.FetchInventory(base, inv)
	if err != nil {
		// Handle error
		return
	}

	// Use the fetched inventory
	for _, armor := range inv.Armor {
		_ = armor.Name
	}
}
