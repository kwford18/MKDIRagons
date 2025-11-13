package inventory_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/internal/inventory"
	"github.com/kwford18/MKDIRagons/internal/reference"
	"github.com/kwford18/MKDIRagons/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

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
				},
				ArmorClass: inventory.ArmorClass{Base: 11, DexBonus: true},
			})
			if err != nil {
				return
			}
		case "/equipment/shield":
			err := json.NewEncoder(w).Encode(inventory.Armor{
				BaseEquipment: inventory.BaseEquipment{
					Index: "shield",
					Name:  "Shield",
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
				},
				ArmorClass: inventory.ArmorClass{Base: 16},
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
				},
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
				},
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

// ======= Tests =======

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

	// Use ElementsMatch to avoid order issues
	assert.ElementsMatch(suite.T(), []string{"leather-armor", "shield"}, getArmorIndexes(inv))
	assert.ElementsMatch(suite.T(), []string{"longsword", "shortbow"}, getWeaponIndexes(inv))
	assert.ElementsMatch(suite.T(), []string{"rope", "torch"}, getItemIndexes(inv))
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

	assert.ElementsMatch(suite.T(), []string{"leather-armor", "shield", "chain-mail"}, getArmorIndexes(inv))
	assert.ElementsMatch(suite.T(), []string{"longsword", "shortbow", "dagger"}, getWeaponIndexes(inv))
	assert.ElementsMatch(suite.T(), []string{"rope", "torch", "backpack"}, getItemIndexes(inv))
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
	assert.ElementsMatch(suite.T(), []string{"longsword", "dagger"}, getWeaponIndexes(inv))
	assert.Empty(suite.T(), inv.Items)
}

// ======= Helpers =======

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

func TestFetchInventoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(FetchInventoryIntegrationTestSuite))
}
