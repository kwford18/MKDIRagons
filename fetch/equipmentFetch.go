package fetch

import "fmt"

// =============================================================
// EQUIPMENT
type Equipment struct {
	Desc              []string    `json:"desc"`
	Special           []string    `json:"special"`
	Index             string      `json:"index"`
	Name              string      `json:"name"`
	EquipmentCategory Reference   `json:"equipment_category"`
	WeaponCategory    string      `json:"weapon_category"`
	WeaponRange       string      `json:"weapon_range"`
	CategoryRange     string      `json:"category_range"`
	Cost              Cost        `json:"cost"`
	Damage            Damage      `json:"damage"`
	Range             WeaponRange `json:"range"`
	Weight            int         `json:"weight"`
	Properties        []Reference `json:"properties"`
	URL               string      `json:"url"`
}

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

func (e *Equipment) GetEndpoint() string {
	return "equipment/"
}

func (e *Equipment) Print() {
	fmt.Printf("Equipment: %v", e.Name)
}

// =============================================================
// INVENTORY
type Inventory struct {
	Weapons, Armor, Items []Equipment
}

func (i *Inventory) GetEndpoint() string {
	return "equipment/"
}

func (i *Inventory) Print() {
	fmt.Printf("Inventory contains: %d armor, %d weapons, %d items\n",
		len(i.Armor), len(i.Weapons), len(i.Items))
}

func (i *Inventory) PrintAll() {
	// Equipment
	fmt.Println("Equipment:")

	fmt.Printf("    - Armor: \n")
	for _, armor := range i.Armor {
		fmt.Printf("	- %s", armor.Name)
	}

	fmt.Printf("    - Weapons: \n")
	for _, weapons := range i.Weapons {
		fmt.Printf("	- %s", weapons.Name)
	}

	fmt.Printf("    - Items: \n")
	for _, items := range i.Items {
		fmt.Printf("	- %s", items.Name)
	}
}
