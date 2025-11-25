package inventory_test

import (
	"testing"

	"github.com/kwford18/MKDIRagons/internal/inventory"
	"github.com/kwford18/MKDIRagons/internal/reference"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// InventoryTestSuite groups all inventory-related tests
type InventoryTestSuite struct {
	suite.Suite
	sampleItem   inventory.Item
	sampleArmor  inventory.Armor
	sampleWeapon inventory.Weapon
}

// SetupTest runs before each test
func (suite *InventoryTestSuite) SetupTest() {
	suite.sampleItem = inventory.Item{
		BaseEquipment: inventory.BaseEquipment{
			Index: "abacus",
			Name:  "Abacus",
			Desc:  []string{"A calculating tool"},
			EquipmentCategory: reference.Reference{
				Index: "adventuring-gear",
				Name:  "Adventuring Gear",
				URL:   "/api/equipment-categories/adventuring-gear",
			},
			Cost: inventory.Cost{
				Quantity: 2,
				Unit:     "gp",
			},
			Weight: 2,
			URL:    "/api/equipment/abacus",
		},
		GearCategory: &reference.Reference{
			Index: "standard-gear",
			Name:  "Standard Gear",
		},
	}

	suite.sampleArmor = inventory.Armor{
		BaseEquipment: inventory.BaseEquipment{
			Index: "padded",
			Name:  "Padded Armor",
			Desc:  []string{"Padded armor consists of quilted layers"},
			EquipmentCategory: reference.Reference{
				Index: "armor",
				Name:  "Armor",
				URL:   "/api/equipment-categories/armor",
			},
			Cost: inventory.Cost{
				Quantity: 5,
				Unit:     "gp",
			},
			Weight: 8,
			URL:    "/api/equipment/padded",
		},
		ArmorCategory: "Light",
		ArmorClass: inventory.ArmorClass{
			Base:     11,
			DexBonus: true,
		},
		StrMinimum:          0,
		StealthDisadvantage: true,
	}

	suite.sampleWeapon = inventory.Weapon{
		BaseEquipment: inventory.BaseEquipment{
			Index: "longsword",
			Name:  "Longsword",
			Desc:  []string{"A versatile blade"},
			EquipmentCategory: reference.Reference{
				Index: "weapon",
				Name:  "Weapon",
				URL:   "/api/equipment-categories/weapon",
			},
			Cost: inventory.Cost{
				Quantity: 15,
				Unit:     "gp",
			},
			Weight: 3,
			URL:    "/api/equipment/longsword",
			Properties: []reference.Reference{
				{Index: "versatile", Name: "Versatile"},
			},
		},
		WeaponCategory: "Martial",
		WeaponRange:    "Melee",
		CategoryRange:  "Martial Melee",
		Damage: inventory.Damage{
			DamageDice: "1d8",
			DamageType: reference.Reference{
				Index: "slashing",
				Name:  "Slashing",
			},
		},
		Range: inventory.WeaponRange{
			Normal: 5,
			Long:   5,
		},
		TwoHandedDamage: &inventory.Damage{
			DamageDice: "1d10",
			DamageType: reference.Reference{
				Index: "slashing",
				Name:  "Slashing",
			},
		},
	}
}

// TestInventoryTestSuite runs the test suite
func TestInventoryTestSuite(t *testing.T) {
	suite.Run(t, new(InventoryTestSuite))
}

// Test Item methods
func (suite *InventoryTestSuite) TestItemGetEndpoint() {
	endpoint := suite.sampleItem.GetEndpoint()
	assert.Equal(suite.T(), "equipment/abacus", endpoint)
}

func (suite *InventoryTestSuite) TestItemGetEndpointWithDifferentIndex() {
	item := inventory.Item{
		BaseEquipment: inventory.BaseEquipment{
			Index: "rope-hempen-50-feet",
		},
	}
	endpoint := item.GetEndpoint()
	assert.Equal(suite.T(), "equipment/rope-hempen-50-feet", endpoint)
}

func (suite *InventoryTestSuite) TestItemPrint() {
	// This test verifies Print doesn't panic
	assert.NotPanics(suite.T(), func() {
		suite.sampleItem.Print()
	})
}

// Test Armor methods
func (suite *InventoryTestSuite) TestArmorGetEndpoint() {
	endpoint := suite.sampleArmor.GetEndpoint()
	assert.Equal(suite.T(), "equipment/padded", endpoint)
}

func (suite *InventoryTestSuite) TestArmorGetEndpointWithDifferentIndex() {
	armor := inventory.Armor{
		BaseEquipment: inventory.BaseEquipment{
			Index: "plate",
		},
	}
	endpoint := armor.GetEndpoint()
	assert.Equal(suite.T(), "equipment/plate", endpoint)
}

func (suite *InventoryTestSuite) TestArmorPrint() {
	assert.NotPanics(suite.T(), func() {
		suite.sampleArmor.Print()
	})
}

func (suite *InventoryTestSuite) TestArmorClassStructure() {
	ac := suite.sampleArmor.ArmorClass
	assert.Equal(suite.T(), 11, ac.Base)
	assert.True(suite.T(), ac.DexBonus)
}

func (suite *InventoryTestSuite) TestArmorStealthDisadvantage() {
	assert.True(suite.T(), suite.sampleArmor.StealthDisadvantage)

	// Test armor without stealth disadvantage
	leatherArmor := inventory.Armor{
		ArmorCategory:       "Light",
		StealthDisadvantage: false,
	}
	assert.False(suite.T(), leatherArmor.StealthDisadvantage)
}

// Test Weapon methods
func (suite *InventoryTestSuite) TestWeaponGetEndpoint() {
	endpoint := suite.sampleWeapon.GetEndpoint()
	assert.Equal(suite.T(), "equipment/longsword", endpoint)
}

func (suite *InventoryTestSuite) TestWeaponGetEndpointWithDifferentIndex() {
	weapon := inventory.Weapon{
		BaseEquipment: inventory.BaseEquipment{
			Index: "longbow",
		},
	}
	endpoint := weapon.GetEndpoint()
	assert.Equal(suite.T(), "equipment/longbow", endpoint)
}

func (suite *InventoryTestSuite) TestWeaponPrint() {
	assert.NotPanics(suite.T(), func() {
		suite.sampleWeapon.Print()
	})
}

func (suite *InventoryTestSuite) TestWeaponDamage() {
	damage := suite.sampleWeapon.Damage
	assert.Equal(suite.T(), "1d8", damage.DamageDice)
	assert.Equal(suite.T(), "slashing", damage.DamageType.Index)
}

func (suite *InventoryTestSuite) TestWeaponTwoHandedDamage() {
	assert.NotNil(suite.T(), suite.sampleWeapon.TwoHandedDamage)
	assert.Equal(suite.T(), "1d10", suite.sampleWeapon.TwoHandedDamage.DamageDice)
}

func (suite *InventoryTestSuite) TestWeaponWithoutTwoHandedDamage() {
	dagger := inventory.Weapon{
		BaseEquipment: inventory.BaseEquipment{
			Index: "dagger",
			Name:  "Dagger",
		},
		Damage: inventory.Damage{
			DamageDice: "1d4",
		},
		TwoHandedDamage: nil,
	}
	assert.Nil(suite.T(), dagger.TwoHandedDamage)
}

func (suite *InventoryTestSuite) TestWeaponRange() {
	weaponRange := suite.sampleWeapon.Range
	assert.Equal(suite.T(), 5, weaponRange.Normal)
	assert.Equal(suite.T(), 5, weaponRange.Long)
}

// Test Inventory methods
func (suite *InventoryTestSuite) TestInventoryGetEndpoint() {
	inv := &inventory.Inventory{}
	endpoint := inv.GetEndpoint()
	assert.Equal(suite.T(), "equipment/", endpoint)
}

func (suite *InventoryTestSuite) TestInventoryPrint() {
	inv := &inventory.Inventory{
		Items:   []inventory.Item{suite.sampleItem},
		Armor:   []inventory.Armor{suite.sampleArmor},
		Weapons: []inventory.Weapon{suite.sampleWeapon},
	}

	assert.NotPanics(suite.T(), func() {
		inv.Print()
	})
}

func (suite *InventoryTestSuite) TestInventoryPrintEmpty() {
	inv := &inventory.Inventory{}

	assert.NotPanics(suite.T(), func() {
		inv.Print()
	})
}

func (suite *InventoryTestSuite) TestInventoryWithMultipleItems() {
	secondItem := inventory.Item{
		BaseEquipment: inventory.BaseEquipment{
			Index: "torch",
			Name:  "Torch",
		},
	}

	inv := &inventory.Inventory{
		Items:   []inventory.Item{suite.sampleItem, secondItem},
		Armor:   []inventory.Armor{suite.sampleArmor},
		Weapons: []inventory.Weapon{suite.sampleWeapon},
	}

	assert.Len(suite.T(), inv.Items, 2)
	assert.Len(suite.T(), inv.Armor, 1)
	assert.Len(suite.T(), inv.Weapons, 1)
}

// Test Cost structure
func (suite *InventoryTestSuite) TestCostStructure() {
	cost := suite.sampleItem.Cost
	assert.Equal(suite.T(), 2, cost.Quantity)
	assert.Equal(suite.T(), "gp", cost.Unit)
}

func (suite *InventoryTestSuite) TestCostWithDifferentUnits() {
	silverCost := inventory.Cost{
		Quantity: 50,
		Unit:     "sp",
	}
	assert.Equal(suite.T(), 50, silverCost.Quantity)
	assert.Equal(suite.T(), "sp", silverCost.Unit)

	copperCost := inventory.Cost{
		Quantity: 100,
		Unit:     "cp",
	}
	assert.Equal(suite.T(), 100, copperCost.Quantity)
	assert.Equal(suite.T(), "cp", copperCost.Unit)
}

// Test BaseEquipment fields
func (suite *InventoryTestSuite) TestBaseEquipmentDesc() {
	assert.NotEmpty(suite.T(), suite.sampleItem.Desc)
	assert.Contains(suite.T(), suite.sampleItem.Desc[0], "calculating")
}

func (suite *InventoryTestSuite) TestBaseEquipmentSpecial() {
	itemWithSpecial := inventory.Item{
		BaseEquipment: inventory.BaseEquipment{
			Special: []string{"Grants advantage on Intelligence checks"},
		},
	}
	assert.NotEmpty(suite.T(), itemWithSpecial.Special)
	assert.Len(suite.T(), itemWithSpecial.Special, 1)
}

func (suite *InventoryTestSuite) TestBaseEquipmentProperties() {
	assert.NotEmpty(suite.T(), suite.sampleWeapon.Properties)
	assert.Equal(suite.T(), "versatile", suite.sampleWeapon.Properties[0].Index)
}

func (suite *InventoryTestSuite) TestBaseEquipmentContents() {
	pack := inventory.Item{
		BaseEquipment: inventory.BaseEquipment{
			Contents: []reference.Reference{
				{Index: "rope", Name: "Rope"},
				{Index: "torch", Name: "Torch"},
			},
		},
	}
	assert.Len(suite.T(), pack.Contents, 2)
}

// Test edge cases
func (suite *InventoryTestSuite) TestEmptyIndexGetEndpoint() {
	item := inventory.Item{}
	endpoint := item.GetEndpoint()
	assert.Equal(suite.T(), "equipment/", endpoint)
}

func (suite *InventoryTestSuite) TestItemWithNilGearCategory() {
	item := inventory.Item{
		BaseEquipment: inventory.BaseEquipment{
			Index: "test",
		},
		GearCategory: nil,
	}
	assert.Nil(suite.T(), item.GearCategory)
}

func (suite *InventoryTestSuite) TestWeaponCategories() {
	assert.Equal(suite.T(), "Martial", suite.sampleWeapon.WeaponCategory)
	assert.Equal(suite.T(), "Melee", suite.sampleWeapon.WeaponRange)
	assert.Equal(suite.T(), "Martial Melee", suite.sampleWeapon.CategoryRange)
}

func (suite *InventoryTestSuite) TestArmorCategories() {
	assert.Equal(suite.T(), "Light", suite.sampleArmor.ArmorCategory)

	// Test other armor categories
	mediumArmor := inventory.Armor{ArmorCategory: "Medium"}
	heavyArmor := inventory.Armor{ArmorCategory: "Heavy"}
	shield := inventory.Armor{ArmorCategory: "Shield"}

	assert.Equal(suite.T(), "Medium", mediumArmor.ArmorCategory)
	assert.Equal(suite.T(), "Heavy", heavyArmor.ArmorCategory)
	assert.Equal(suite.T(), "Shield", shield.ArmorCategory)
}

func (suite *InventoryTestSuite) TestArmorStrengthMinimum() {
	assert.Equal(suite.T(), 0, suite.sampleArmor.StrMinimum)

	plateArmor := inventory.Armor{
		StrMinimum: 15,
	}
	assert.Equal(suite.T(), 15, plateArmor.StrMinimum)
}
