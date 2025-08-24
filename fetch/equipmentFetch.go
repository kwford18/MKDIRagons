package fetch

import "fmt"

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
