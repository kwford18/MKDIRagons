package character_test

import (
	"errors"
	"github.com/kwford18/MKDIRagons/internal/character"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/class"
	"github.com/kwford18/MKDIRagons/internal/inventory"
	"github.com/kwford18/MKDIRagons/internal/race"
	"github.com/kwford18/MKDIRagons/internal/reference"
	"github.com/kwford18/MKDIRagons/template"
	"github.com/stretchr/testify/assert"
)

// MockFetcher simulates the core.Fetcher interface.
type MockFetcher struct {
	ShouldFail bool
	FailOn     string // "race", "class", "inventory", or "spells"
}

// FetchJSON mocks the network call by populating the property based on its type.
func (m *MockFetcher) FetchJSON(property reference.Fetchable, input string) error {
	if m.ShouldFail {
		// Simple logic to fail specific calls based on what we are fetching
		switch property.(type) {
		case *race.Race:
			if m.FailOn == "race" {
				return errors.New("failed to fetch race")
			}
		case *class.Class:
			if m.FailOn == "class" {
				return errors.New("failed to fetch class")
			}
		case *inventory.Item, *inventory.Armor:
			if m.FailOn == "inventory" {
				return errors.New("failed to fetch inventory item")
			}
		}
	}

	// Populate property with dummy data based on type
	switch d := property.(type) {
	case *race.Race:
		d.Name = "TestRace"
		d.Speed = 30
	case *class.Class:
		d.Name = "TestClass"
		d.HitDie = 10 // Important for Stats calculation
	// The builder iterates and fetches individual items/armor
	case *inventory.Item:
		d.BaseEquipment.Name = "Test Item"
	case *inventory.Armor:
		d.BaseEquipment.Name = "Leather Armor"
		d.ArmorClass = inventory.ArmorClass{
			Base:     11,
			DexBonus: true,
		}
	case *inventory.Inventory:
		// Fallback if the builder fetches the whole inventory container (unlikely but possible)
		d.Items = []inventory.Item{{
			BaseEquipment: inventory.BaseEquipment{
				Name: "Test Item",
			},
		}}
		d.Armor = []inventory.Armor{{
			BaseEquipment: inventory.BaseEquipment{
				Name: "Leather Armor",
			},
			ArmorClass: inventory.ArmorClass{
				Base:     11,
				DexBonus: true,
			},
		}}
	default:
		// Handle other types if necessary
	}

	return nil
}

func TestBuildCharacterWithFetcher_Success(t *testing.T) {
	// Setup
	base := &template.Character{
		Name:          "TestHero",
		Level:         5,
		Proficiencies: []string{"Stealth"},
		// Corrected: Use the nested Inventory struct as defined in template/character.go
		Inventory: template.Inventory{
			Armor: []string{"Leather Armor"},
			Items: []string{"Test Item"},
		},
		AbilityScores: template.AbilityScores{
			Dexterity:    14,
			Constitution: 12,
		},
	}

	fetcher := &MockFetcher{ShouldFail: false}

	// Execute
	char, err := character.BuildCharacterWithFetcher(fetcher, base, false)

	// Verify
	assert.NoError(t, err)
	assert.NotNil(t, char)

	// Check basic fields
	assert.Equal(t, "TestHero", char.Name)
	assert.Equal(t, 5, char.Level)
	assert.Equal(t, "TestRace", char.Race.Name)
	assert.Equal(t, "TestClass", char.Class.Name)

	// Check Derived Stats (HitDie 10, Level 5, Con +1)
	// Avg HP for d10 is 6.
	// HP = (6 + 1) * 5 = 35
	assert.Equal(t, 35, char.Stats.HP)

	// Check AC (Leather Armor 11 + Dex 2 = 13)
	assert.Equal(t, 13, char.Stats.AC)

	// Check Inventory
	// Note: We check that at least one item was added.
	// The mock might populate "Test Item" for both calls depending on how FetchInventoryWithFetcher distinguishes types,
	// or "Leather Armor" for the armor call.
	assert.NotEmpty(t, char.Inventory.Items, "Inventory items should not be empty")

	// Check specifically for the item name if available in Items list
	// (Adjust index if Armor ends up in Items list or vice versa depending on logic)
	foundItem := false
	for _, item := range char.Inventory.Items {
		if item.Name == "Test Item" {
			foundItem = true
			break
		}
	}
	assert.True(t, foundItem, "Expected to find 'Test Item' in inventory")
}

func TestBuildCharacterWithFetcher_FetchError(t *testing.T) {
	// Setup
	base := &template.Character{
		Name:  "TestHero",
		Level: 1,
	}

	// Configure mock to fail on Race fetch
	fetcher := &MockFetcher{
		ShouldFail: true,
		FailOn:     "race",
	}

	// Execute
	char, err := character.BuildCharacterWithFetcher(fetcher, base, false)

	// Verify
	assert.Error(t, err)
	assert.Nil(t, char)
	assert.Contains(t, err.Error(), "failed to fetch race")
}

func TestBuildCharacterWithFetcher_ConcurrentErrorHandling(t *testing.T) {
	// Ensure that even if multiple fail, or one fails late, we get an error return.
	base := &template.Character{Name: "TestHero", Level: 1}

	fetcher := &MockFetcher{
		ShouldFail: true,
		FailOn:     "class",
	}

	char, err := character.BuildCharacterWithFetcher(fetcher, base, false)

	assert.Error(t, err)
	assert.Nil(t, char)
	assert.Contains(t, err.Error(), "failed to fetch class")
}
