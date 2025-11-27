package inventory_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/internal/inventory"
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
// UNIT TEST SUITE - With Fixtures
// ============================================================================

type InventoryBuilderTestSuite struct {
	suite.Suite
	mockFetcher         *MockFetcher
	fixtureBasedFetcher *MockFetcherWithFixtures
	baseCharacter       *template.Character
	inventory           *inventory.Inventory
}

func (suite *InventoryBuilderTestSuite) SetupTest() {
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
		Inventory: template.Inventory{
			Weapons: []string{"longsword"},
			Armor:   []string{"plate-armor"},
			Items:   []string{"abacus"},
		},
	}

	suite.inventory = &inventory.Inventory{}
}

func TestInventoryBuilderTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryBuilderTestSuite))
}

func TestFetchInventoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(FetchInventoryIntegrationTestSuite))
}

// ============================================================================
// UNIT TESTS - Behavior Testing (No Real Data)
// ============================================================================

func (suite *InventoryBuilderTestSuite) TestFetchInventoryWithFetcher_Success() {
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Weapon"), "longsword").Return(nil)
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), "plate-armor").Return(nil)
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Item"), "abacus").Return(nil)

	err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, suite.baseCharacter, suite.inventory)

	assert.NoError(suite.T(), err)
	suite.mockFetcher.AssertExpectations(suite.T())
}

func (suite *InventoryBuilderTestSuite) TestFetchInventoryWithFetcher_NetworkError() {
	expectedError := errors.New("network timeout")
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Weapon"), "longsword").Return(expectedError)
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), "plate-armor").Return(nil)
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Item"), "abacus").Return(nil)

	err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, suite.baseCharacter, suite.inventory)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
}

func (suite *InventoryBuilderTestSuite) TestFetchInventoryWithFetcher_NilCharacter() {
	assert.Panics(suite.T(), func() {
		_ = inventory.FetchInventoryWithFetcher(suite.mockFetcher, nil, suite.inventory)
	})
}

func (suite *InventoryBuilderTestSuite) TestFetchInventoryWithFetcher_NilInventory() {
	assert.Panics(suite.T(), func() {
		_ = inventory.FetchInventoryWithFetcher(suite.mockFetcher, suite.baseCharacter, nil)
	})
}

func (suite *InventoryBuilderTestSuite) TestFetchInventoryWithFetcher_EmptyInventory() {
	emptyCharacter := &template.Character{
		Name:          "Test",
		Level:         1,
		Race:          "human",
		Class:         "fighter",
		AbilityScores: template.AbilityScores{Strength: 10},
		Inventory: template.Inventory{
			Weapons: []string{},
			Armor:   []string{},
			Items:   []string{},
		},
	}

	err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, emptyCharacter, suite.inventory)

	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), suite.inventory.Weapons)
	assert.Empty(suite.T(), suite.inventory.Armor)
	assert.Empty(suite.T(), suite.inventory.Items)
}

func (suite *InventoryBuilderTestSuite) TestFetchInventoryWithFetcher_MultipleItems() {
	multiItemCharacter := &template.Character{
		Name:          "Test",
		Level:         1,
		Race:          "human",
		Class:         "fighter",
		AbilityScores: template.AbilityScores{Strength: 10},
		Inventory: template.Inventory{
			Weapons: []string{"longsword", "dagger"},
			Armor:   []string{"plate-armor", "leather-armor"},
			Items:   []string{"abacus"},
		},
	}

	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Weapon"), "longsword").Return(nil)
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Weapon"), "dagger").Return(nil)
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), "plate-armor").Return(nil)
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), "leather-armor").Return(nil)
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Item"), "abacus").Return(nil)

	err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, multiItemCharacter, suite.inventory)

	assert.NoError(suite.T(), err)
	suite.mockFetcher.AssertExpectations(suite.T())
}

// ============================================================================
// UNIT TESTS - With Fixture Data (Realistic Testing)
// ============================================================================

func (suite *InventoryBuilderTestSuite) TestFetchInventoryWithFetcher_WithFixtures() {
	suite.fixtureBasedFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Weapon"), "longsword").Return(nil)
	suite.fixtureBasedFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), "plate-armor").Return(nil)
	suite.fixtureBasedFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Item"), "abacus").Return(nil)

	err := inventory.FetchInventoryWithFetcher(suite.fixtureBasedFetcher, suite.baseCharacter, suite.inventory)

	// Give goroutines time to complete
	time.Sleep(100 * time.Millisecond)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), suite.inventory.Weapons, 1)
	assert.Len(suite.T(), suite.inventory.Armor, 1)
	assert.Len(suite.T(), suite.inventory.Items, 1)

	// Verify weapon data
	weapon := suite.inventory.Weapons[0]
	assert.Equal(suite.T(), "Longsword", weapon.Name)
	assert.Equal(suite.T(), "longsword", weapon.Index)
	assert.Equal(suite.T(), "Martial", weapon.WeaponCategory)
	assert.Equal(suite.T(), "1d8", weapon.Damage.DamageDice)
	assert.NotNil(suite.T(), weapon.TwoHandedDamage)
	assert.Equal(suite.T(), "1d10", weapon.TwoHandedDamage.DamageDice)

	// Verify armor data
	armor := suite.inventory.Armor[0]
	assert.Equal(suite.T(), "Plate Armor", armor.Name)
	assert.Equal(suite.T(), "Heavy", armor.ArmorCategory)
	assert.Equal(suite.T(), 18, armor.ArmorClass.Base)
	assert.False(suite.T(), armor.ArmorClass.DexBonus)
	assert.Equal(suite.T(), 15, armor.StrMinimum)
	assert.True(suite.T(), armor.StealthDisadvantage)

	// Verify item data
	item := suite.inventory.Items[0]
	assert.Equal(suite.T(), "Abacus", item.Name)
	assert.Equal(suite.T(), 2, item.Cost.Quantity)
	assert.Equal(suite.T(), "gp", item.Cost.Unit)

	suite.fixtureBasedFetcher.AssertExpectations(suite.T())
}

func (suite *InventoryBuilderTestSuite) TestFetchInventoryWithFetcher_MultipleWeaponsWithFixtures() {
	multiWeaponCharacter := &template.Character{
		Name:          "Rogue",
		Level:         3,
		Race:          "halfling",
		Class:         "rogue",
		AbilityScores: template.AbilityScores{Dexterity: 18},
		Inventory: template.Inventory{
			Weapons: []string{"dagger", "longsword"},
			Armor:   []string{"leather-armor"},
			Items:   []string{},
		},
	}

	suite.fixtureBasedFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Weapon"), "dagger").Return(nil)
	suite.fixtureBasedFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Weapon"), "longsword").Return(nil)
	suite.fixtureBasedFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), "leather-armor").Return(nil)

	err := inventory.FetchInventoryWithFetcher(suite.fixtureBasedFetcher, multiWeaponCharacter, suite.inventory)

	time.Sleep(100 * time.Millisecond)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), suite.inventory.Weapons, 2)
	assert.Len(suite.T(), suite.inventory.Armor, 1)

	// Find dagger and longsword (order may vary due to concurrency)
	var dagger, longsword *inventory.Weapon
	for i := range suite.inventory.Weapons {
		if suite.inventory.Weapons[i].Index == "dagger" {
			dagger = &suite.inventory.Weapons[i]
		} else if suite.inventory.Weapons[i].Index == "longsword" {
			longsword = &suite.inventory.Weapons[i]
		}
	}

	assert.NotNil(suite.T(), dagger)
	assert.NotNil(suite.T(), longsword)

	// Verify dagger
	assert.Equal(suite.T(), "Dagger", dagger.Name)
	assert.Equal(suite.T(), "Simple", dagger.WeaponCategory)
	assert.Equal(suite.T(), "1d4", dagger.Damage.DamageDice)
	assert.Nil(suite.T(), dagger.TwoHandedDamage)

	// Verify longsword
	assert.Equal(suite.T(), "Longsword", longsword.Name)
	assert.Equal(suite.T(), "Martial", longsword.WeaponCategory)
	assert.NotNil(suite.T(), longsword.TwoHandedDamage)

	suite.fixtureBasedFetcher.AssertExpectations(suite.T())
}

func (suite *InventoryBuilderTestSuite) TestFetchInventoryWithFetcher_VerifyArmorProperties() {
	armorCharacter := &template.Character{
		Name:          "Paladin",
		Level:         5,
		Race:          "human",
		Class:         "paladin",
		AbilityScores: template.AbilityScores{Strength: 16},
		Inventory: template.Inventory{
			Weapons: []string{},
			Armor:   []string{"plate-armor", "leather-armor"},
			Items:   []string{},
		},
	}

	suite.fixtureBasedFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), "plate-armor").Return(nil)
	suite.fixtureBasedFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), "leather-armor").Return(nil)

	err := inventory.FetchInventoryWithFetcher(suite.fixtureBasedFetcher, armorCharacter, suite.inventory)

	time.Sleep(100 * time.Millisecond)

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), suite.inventory.Armor, 2)

	// Find plate and leather
	var plate, leather *inventory.Armor
	for i := range suite.inventory.Armor {
		if suite.inventory.Armor[i].Index == "plate-armor" {
			plate = &suite.inventory.Armor[i]
		} else if suite.inventory.Armor[i].Index == "leather-armor" {
			leather = &suite.inventory.Armor[i]
		}
	}

	assert.NotNil(suite.T(), plate)
	assert.NotNil(suite.T(), leather)

	// Verify plate (heavy armor)
	assert.Equal(suite.T(), "Heavy", plate.ArmorCategory)
	assert.Equal(suite.T(), 18, plate.ArmorClass.Base)
	assert.False(suite.T(), plate.ArmorClass.DexBonus)
	assert.True(suite.T(), plate.StealthDisadvantage)
	assert.Equal(suite.T(), 15, plate.StrMinimum)

	// Verify leather (light armor)
	assert.Equal(suite.T(), "Light", leather.ArmorCategory)
	assert.Equal(suite.T(), 11, leather.ArmorClass.Base)
	assert.True(suite.T(), leather.ArmorClass.DexBonus)
	assert.False(suite.T(), leather.StealthDisadvantage)
	assert.Equal(suite.T(), 0, leather.StrMinimum)

	suite.fixtureBasedFetcher.AssertExpectations(suite.T())
}

// ============================================================================
// CONCURRENCY TESTS
// ============================================================================

func (suite *InventoryBuilderTestSuite) TestFetchInventoryWithFetcher_ConcurrentFetches() {
	largeInventory := &template.Character{
		Name:          "Adventurer",
		Level:         10,
		Race:          "human",
		Class:         "fighter",
		AbilityScores: template.AbilityScores{Strength: 16},
		Inventory: template.Inventory{
			Weapons: []string{"longsword", "dagger"},
			Armor:   []string{"plate-armor", "leather-armor"},
			Items:   []string{"abacus"},
		},
	}

	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Weapon"), mock.Anything).Return(nil)
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), mock.Anything).Return(nil)
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Item"), mock.Anything).Return(nil)

	err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, largeInventory, suite.inventory)

	assert.NoError(suite.T(), err)
	suite.mockFetcher.AssertNumberOfCalls(suite.T(), "FetchJSON", 5)
}

func (suite *InventoryBuilderTestSuite) TestFetchInventoryWithFetcher_ErrorInOneConcurrentFetch() {
	character := &template.Character{
		Name:          "Test",
		Level:         1,
		Race:          "human",
		Class:         "fighter",
		AbilityScores: template.AbilityScores{Strength: 10},
		Inventory: template.Inventory{
			Weapons: []string{"longsword", "dagger"},
			Armor:   []string{"plate-armor"},
			Items:   []string{},
		},
	}

	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Weapon"), "longsword").Return(nil)
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Weapon"), "dagger").Return(errors.New("404 not found"))
	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), "plate-armor").Return(nil)

	err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, character, suite.inventory)

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "404")
}

// ============================================================================
// EDGE CASE TESTS
// ============================================================================

func (suite *InventoryBuilderTestSuite) TestFetchInventoryWithFetcher_OnlyWeapons() {
	weaponsOnlyChar := &template.Character{
		Name:          "Test",
		Level:         1,
		Race:          "human",
		Class:         "fighter",
		AbilityScores: template.AbilityScores{Strength: 10},
		Inventory: template.Inventory{
			Weapons: []string{"longsword"},
			Armor:   []string{},
			Items:   []string{},
		},
	}

	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Weapon"), "longsword").Return(nil)

	err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, weaponsOnlyChar, suite.inventory)

	assert.NoError(suite.T(), err)
	suite.mockFetcher.AssertNumberOfCalls(suite.T(), "FetchJSON", 1)
}

func (suite *InventoryBuilderTestSuite) TestFetchInventoryWithFetcher_OnlyArmor() {
	armorOnlyChar := &template.Character{
		Name:          "Test",
		Level:         1,
		Race:          "human",
		Class:         "fighter",
		AbilityScores: template.AbilityScores{Strength: 10},
		Inventory: template.Inventory{
			Weapons: []string{},
			Armor:   []string{"plate-armor"},
			Items:   []string{},
		},
	}

	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Armor"), "plate-armor").Return(nil)

	err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, armorOnlyChar, suite.inventory)

	assert.NoError(suite.T(), err)
	suite.mockFetcher.AssertNumberOfCalls(suite.T(), "FetchJSON", 1)
}

func (suite *InventoryBuilderTestSuite) TestFetchInventoryWithFetcher_OnlyItems() {
	itemsOnlyChar := &template.Character{
		Name:          "Test",
		Level:         1,
		Race:          "human",
		Class:         "wizard",
		AbilityScores: template.AbilityScores{Intelligence: 16},
		Inventory: template.Inventory{
			Weapons: []string{},
			Armor:   []string{},
			Items:   []string{"abacus"},
		},
	}

	suite.mockFetcher.On("FetchJSON", mock.AnythingOfType("*inventory.Item"), "abacus").Return(nil)

	err := inventory.FetchInventoryWithFetcher(suite.mockFetcher, itemsOnlyChar, suite.inventory)

	assert.NoError(suite.T(), err)
	suite.mockFetcher.AssertNumberOfCalls(suite.T(), "FetchJSON", 1)
}

// ============================================================================
// INTEGRATION TEST SUITE - With HTTP Server
// ============================================================================

type FetchInventoryIntegrationTestSuite struct {
	suite.Suite
	server  *httptest.Server
	fetcher *core.HTTPFetcher
}

func (suite *FetchInventoryIntegrationTestSuite) SetupSuite() {
	suite.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		// Armor
		case "/equipment/leather-armor":
			err := json.NewEncoder(w).Encode(inventory.Armor{
				BaseEquipment: inventory.BaseEquipment{
					Index: "leather-armor",
					Name:  "Leather Armor",
					Cost:  inventory.Cost{Quantity: 10, Unit: "gp"},
				},
				ArmorCategory: "Light",
				ArmorClass:    inventory.ArmorClass{Base: 11, DexBonus: true},
			})
			if err != nil {
				return
			}
		case "/equipment/shield":
			err := json.NewEncoder(w).Encode(inventory.Armor{
				BaseEquipment: inventory.BaseEquipment{
					Index: "shield",
					Name:  "Shield",
					Cost:  inventory.Cost{Quantity: 10, Unit: "gp"},
				},
				ArmorClass: inventory.ArmorClass{Base: 2},
			})
			if err != nil {
				return
			}
		case "/equipment/chain-mail":
			err := json.NewEncoder(w).Encode(inventory.Armor{
				BaseEquipment: inventory.BaseEquipment{
					Index: "chain-mail",
					Name:  "Chain Mail",
					Cost:  inventory.Cost{Quantity: 75, Unit: "gp"},
				},
				ArmorCategory: "Heavy",
				ArmorClass:    inventory.ArmorClass{Base: 16},
			})
			if err != nil {
				return
			}
		case "/equipment/plate-armor":
			err := json.NewEncoder(w).Encode(inventory.Armor{
				BaseEquipment: inventory.BaseEquipment{
					Index: "plate-armor",
					Name:  "Plate Armor",
					Cost:  inventory.Cost{Quantity: 1500, Unit: "gp"},
				},
				ArmorCategory:       "Heavy",
				ArmorClass:          inventory.ArmorClass{Base: 18, DexBonus: false},
				StrMinimum:          15,
				StealthDisadvantage: true,
			})
			if err != nil {
				return
			}

		// Weapons
		case "/equipment/longsword":
			err := json.NewEncoder(w).Encode(inventory.Weapon{
				BaseEquipment: inventory.BaseEquipment{
					Index: "longsword",
					Name:  "Longsword",
					Cost:  inventory.Cost{Quantity: 15, Unit: "gp"},
				},
				WeaponCategory: "Martial",
				Damage: inventory.Damage{
					DamageDice: "1d8",
					DamageType: reference.Reference{Index: "slashing", Name: "Slashing"},
				},
			})
			if err != nil {
				return
			}
		case "/equipment/shortbow":
			err := json.NewEncoder(w).Encode(inventory.Weapon{
				BaseEquipment: inventory.BaseEquipment{
					Index: "shortbow",
					Name:  "Shortbow",
					Cost:  inventory.Cost{Quantity: 25, Unit: "gp"},
				},
				WeaponCategory: "Simple",
				Damage: inventory.Damage{
					DamageDice: "1d6",
					DamageType: reference.Reference{Index: "piercing", Name: "Piercing"},
				},
			})
			if err != nil {
				return
			}
		case "/equipment/dagger":
			err := json.NewEncoder(w).Encode(inventory.Weapon{
				BaseEquipment: inventory.BaseEquipment{
					Index: "dagger",
					Name:  "Dagger",
					Cost:  inventory.Cost{Quantity: 2, Unit: "gp"},
				},
				WeaponCategory: "Simple",
				Damage: inventory.Damage{
					DamageDice: "1d4",
					DamageType: reference.Reference{Index: "piercing", Name: "Piercing"},
				},
			})
			if err != nil {
				return
			}

		// Items
		case "/equipment/rope":
			err := json.NewEncoder(w).Encode(inventory.Item{
				BaseEquipment: inventory.BaseEquipment{
					Index: "rope",
					Name:  "Rope, Hempen (50 feet)",
					Cost:  inventory.Cost{Quantity: 1, Unit: "gp"},
				},
			})
			if err != nil {
				return
			}
		case "/equipment/torch":
			err := json.NewEncoder(w).Encode(inventory.Item{
				BaseEquipment: inventory.BaseEquipment{
					Index: "torch",
					Name:  "Torch",
					Cost:  inventory.Cost{Quantity: 1, Unit: "cp"},
				},
			})
			if err != nil {
				return
			}
		case "/equipment/backpack":
			err := json.NewEncoder(w).Encode(inventory.Item{
				BaseEquipment: inventory.BaseEquipment{
					Index: "backpack",
					Name:  "Backpack",
					Cost:  inventory.Cost{Quantity: 2, Unit: "gp"},
				},
			})
			if err != nil {
				return
			}
		case "/equipment/abacus":
			err := json.NewEncoder(w).Encode(inventory.Item{
				BaseEquipment: inventory.BaseEquipment{
					Index: "abacus",
					Name:  "Abacus",
					Cost:  inventory.Cost{Quantity: 2, Unit: "gp"},
				},
			})
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
	base := &template.Character{
		Name:          "Fighter",
		Level:         5,
		Race:          "human",
		Class:         "fighter",
		AbilityScores: template.AbilityScores{Strength: 16},
		Inventory: template.Inventory{
			Armor:   []string{"leather-armor", "shield"},
			Weapons: []string{"longsword", "shortbow"},
			Items:   []string{"rope", "torch"},
		},
	}
	inv := &inventory.Inventory{}

	err := inventory.FetchInventoryWithFetcher(suite.fetcher, base, inv)

	time.Sleep(100 * time.Millisecond) // Wait for goroutines

	assert.NoError(suite.T(), err)
	assert.ElementsMatch(suite.T(), []string{"leather-armor", "shield"}, getArmorIndexes(inv))
	assert.ElementsMatch(suite.T(), []string{"longsword", "shortbow"}, getWeaponIndexes(inv))
	assert.ElementsMatch(suite.T(), []string{"rope", "torch"}, getItemIndexes(inv))
}

func (suite *FetchInventoryIntegrationTestSuite) TestFetchInventoryLargeInventory() {
	base := &template.Character{
		Name:          "Adventurer",
		Level:         10,
		Race:          "human",
		Class:         "fighter",
		AbilityScores: template.AbilityScores{Strength: 16},
		Inventory: template.Inventory{
			Armor:   []string{"leather-armor", "shield", "chain-mail"},
			Weapons: []string{"longsword", "shortbow", "dagger"},
			Items:   []string{"rope", "torch", "backpack"},
		},
	}
	inv := &inventory.Inventory{}

	err := inventory.FetchInventoryWithFetcher(suite.fetcher, base, inv)

	time.Sleep(150 * time.Millisecond)

	assert.NoError(suite.T(), err)
	assert.ElementsMatch(suite.T(), []string{"leather-armor", "shield", "chain-mail"}, getArmorIndexes(inv))
	assert.ElementsMatch(suite.T(), []string{"longsword", "shortbow", "dagger"}, getWeaponIndexes(inv))
	assert.ElementsMatch(suite.T(), []string{"rope", "torch", "backpack"}, getItemIndexes(inv))
}

func (suite *FetchInventoryIntegrationTestSuite) TestFetchInventoryOnlyWeapons() {
	base := &template.Character{
		Name:          "Warrior",
		Level:         3,
		Race:          "human",
		Class:         "fighter",
		AbilityScores: template.AbilityScores{Strength: 14},
		Inventory: template.Inventory{
			Armor:   []string{},
			Weapons: []string{"longsword", "dagger"},
			Items:   []string{},
		},
	}
	inv := &inventory.Inventory{}

	err := inventory.FetchInventoryWithFetcher(suite.fetcher, base, inv)
	time.Sleep(100 * time.Millisecond)

	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), inv.Armor)
	assert.ElementsMatch(suite.T(), []string{"longsword", "dagger"}, getWeaponIndexes(inv))
	assert.Empty(suite.T(), inv.Items)
}
func (suite *FetchInventoryIntegrationTestSuite) TestFetchInventoryVerifyData() {
	base := &template.Character{
		Name:          "Paladin",
		Level:         5,
		Race:          "human",
		Class:         "paladin",
		AbilityScores: template.AbilityScores{Strength: 16},
		Inventory: template.Inventory{
			Armor:   []string{"plate-armor"},
			Weapons: []string{"longsword"},
			Items:   []string{"abacus"},
		},
	}
	inv := &inventory.Inventory{}
	err := inventory.FetchInventoryWithFetcher(suite.fetcher, base, inv)

	time.Sleep(100 * time.Millisecond)

	assert.NoError(suite.T(), err)

	// Verify armor details
	assert.Len(suite.T(), inv.Armor, 1)
	assert.Equal(suite.T(), "Plate Armor", inv.Armor[0].Name)
	assert.Equal(suite.T(), "Heavy", inv.Armor[0].ArmorCategory)
	assert.Equal(suite.T(), 18, inv.Armor[0].ArmorClass.Base)
	assert.True(suite.T(), inv.Armor[0].StealthDisadvantage)

	// Verify weapon details
	assert.Len(suite.T(), inv.Weapons, 1)
	assert.Equal(suite.T(), "Longsword", inv.Weapons[0].Name)
	assert.Equal(suite.T(), "Martial", inv.Weapons[0].WeaponCategory)
	assert.Equal(suite.T(), "1d8", inv.Weapons[0].Damage.DamageDice)

	// Verify item details
	assert.Len(suite.T(), inv.Items, 1)
	assert.Equal(suite.T(), "Abacus", inv.Items[0].Name)
}

// ============================================================================
// REAL API INTEGRATION TEST (Optional)
// ============================================================================
func (suite *FetchInventoryIntegrationTestSuite) TestFetchInventory_RealAPI() {
	if testing.Short() {
		suite.T().Skip("Skipping integration test in short mode")
	}
	base := &template.Character{
		Name:          "Real API Test",
		Level:         1,
		Race:          "human",
		Class:         "fighter",
		AbilityScores: template.AbilityScores{Strength: 16},
		Inventory: template.Inventory{
			Weapons: []string{"longsword"},
			Armor:   []string{"leather-armor"},
			Items:   []string{"rope-hempen-50-feet"},
		},
	}

	inv := &inventory.Inventory{}

	err := inventory.FetchInventory(base, inv)

	time.Sleep(500 * time.Millisecond) // Network operations need more time

	if err == nil {
		assert.NotEmpty(suite.T(), inv.Weapons)
		assert.NotEmpty(suite.T(), inv.Armor)
		suite.T().Log("✓ Real API call succeeded")
	} else {
		suite.T().Logf("⚠ API not available: %v", err)
	}
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================
func getArmorIndexes(inv *inventory.Inventory) []string {
	indexes := make([]string, len(inv.Armor))
	for i, a := range inv.Armor {
		indexes[i] = a.Index
	}
	return indexes
}
func getWeaponIndexes(inv *inventory.Inventory) []string {
	indexes := make([]string, len(inv.Weapons))
	for i, w := range inv.Weapons {
		indexes[i] = w.Index
	}
	return indexes
}
func getItemIndexes(inv *inventory.Inventory) []string {
	indexes := make([]string, len(inv.Items))
	for i, it := range inv.Items {
		indexes[i] = it.Index
	}
	return indexes
}
