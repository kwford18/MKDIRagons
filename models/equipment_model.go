package models

import "fmt"

// Base for all types of equipment
type BaseEquipment struct {
	Desc              []string    `json:"desc"`
	Special           []string    `json:"special"`
	Index             string      `json:"index"`
	Name              string      `json:"name"`
	EquipmentCategory Reference   `json:"equipment_category"`
	Cost              Cost        `json:"cost"`
	Weight            int         `json:"weight"`
	URL               string      `json:"url"`
	Properties        []Reference `json:"properties"`
	Contents          []Reference `json:"contents"`
}

// Basic items like abacus, amulet, alchemist fire
type Item struct {
	BaseEquipment
	GearCategory    *Reference `json:"gear_category,omitempty"`
	VehicleCategory string     `json:"vehicle_category,omitempty"`
}

// Armor such as padded, leather, etc
type Armor struct {
	BaseEquipment
	ArmorCategory       string     `json:"armor_category"`
	ArmorClass          ArmorClass `json:"armor_class"`
	StrMinimum          int        `json:"str_minimum"`
	StealthDisadvantage bool       `json:"stealth_disadvantage"`
}

type ArmorClass struct {
	Base     int  `json:"base"`
	DexBonus bool `json:"dex_bonus"`
}

// Weapons such as longbow or rapier
type Weapon struct {
	BaseEquipment
	WeaponCategory  string      `json:"weapon_category"`
	WeaponRange     string      `json:"weapon_range"`
	CategoryRange   string      `json:"category_range"`
	Damage          Damage      `json:"damage"`
	Range           WeaponRange `json:"range"`
	TwoHandedDamage *Damage     `json:"two_handed_damage,omitempty"`
}

// Shared structs for various equipment types
type Cost struct {
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
}

type Damage struct {
	DamageDice string    `json:"damage_dice"`
	DamageType Reference `json:"damage_type"`
}

type WeaponRange struct {
	Normal int `json:"normal"`
	Long   int `json:"long"`
}

type Inventory struct {
	Items   []Item
	Armor   []Armor
	Weapons []Weapon
}

// Fetchable Methods
func (inv *Inventory) GetEndpoint() string {
	return "equipment/"
}

func (inv *Inventory) Print() {
	fmt.Printf("Inventory contains: %d armor, %d weapons, %d items\n",
		len(inv.Armor), len(inv.Weapons), len(inv.Items))

	fmt.Printf("    - Armor: \n")
	for _, armor := range inv.Armor {
		fmt.Printf("	- %s\n", armor.Name)
	}

	fmt.Printf("    - Weapons: \n")
	for _, weapons := range inv.Weapons {
		fmt.Printf("	- %s\n", weapons.Name)
	}

	fmt.Printf("    - Items: \n")
	for _, items := range inv.Items {
		fmt.Printf("	- %s\n", items.Name)
	}
}

func (i *Item) GetEndpoint() string {
	return "equipment/" + i.Index
}

func (i *Item) Print() {
	fmt.Printf("Item: %s\n", i.Name)
}

func (a *Armor) GetEndpoint() string {
	return "equipment/" + a.Index
}

func (a *Armor) Print() {
	fmt.Printf("Armor: %s\n", a.Name)
}

func (w *Weapon) GetEndpoint() string {
	return "equipment/" + w.Index
}

func (w *Weapon) Print() {
	fmt.Printf("Weapon: %s\n", w.Name)
}
