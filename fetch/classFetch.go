package fetch

import "fmt"

type Class struct {
	Index                    string                 `json:"index"`
	Name                     string                 `json:"name"`
	HitDie                   int                    `json:"hit_die"`
	ProficiencyChoices       []ProficiencyChoice    `json:"proficiency_choices"`
	Proficiencies            []Reference            `json:"proficiencies"`
	SavingThrows             []Reference            `json:"saving_throws"`
	StartingEquipment        []StartingEquipment    `json:"starting_equipment"`
	StartingEquipmentOptions []EquipmentOptionGroup `json:"starting_equipment_options"`
	ClassLevels              string                 `json:"class_levels"`
	MultiClassing            MultiClassing          `json:"multi_classing"`
	Subclasses               []Reference            `json:"subclasses"`
	Spellcasting             Spellcasting           `json:"spellcasting"`
	Spells                   string                 `json:"spells"`
	URL                      string                 `json:"url"`
}

// --- Proficiencies and Options ---

type ProficiencyChoice struct {
	Desc   string      `json:"desc"`
	Choose int         `json:"choose"`
	Type   string      `json:"type"`
	From   OptionGroup `json:"from"`
}

type OptionGroup struct {
	OptionSetType string   `json:"option_set_type"`
	Options       []Option `json:"options"`
}

type Option struct {
	OptionType string       `json:"option_type"`
	Item       *Reference   `json:"item,omitempty"`
	Count      int          `json:"count,omitempty"`
	Of         *Reference   `json:"of,omitempty"`
	Choice     *ChoiceGroup `json:"choice,omitempty"`
}

type ChoiceGroup struct {
	Desc   string            `json:"desc"`
	Choose int               `json:"choose"`
	Type   string            `json:"type"`
	From   EquipmentCategory `json:"from"`
}

type EquipmentCategory struct {
	OptionSetType     string    `json:"option_set_type"`
	EquipmentCategory Reference `json:"equipment_category"`
}

// --- Starting Equipment ---

type StartingEquipment struct {
	Equipment Reference `json:"equipment"`
	Quantity  int       `json:"quantity"`
}

type EquipmentOptionGroup struct {
	Desc   string      `json:"desc"`
	Choose int         `json:"choose"`
	Type   string      `json:"type"`
	From   OptionGroup `json:"from"`
}

// --- Multiclassing ---

type MultiClassing struct {
	Prerequisites []Prerequisite `json:"prerequisites"`
	Proficiencies []interface{}  `json:"proficiencies"`
}

type Prerequisite struct {
	AbilityScore Reference `json:"ability_score"`
	MinimumScore int       `json:"minimum_score"`
}

// --- Spellcasting ---

type Spellcasting struct {
	Level               int                `json:"level"`
	SpellcastingAbility Reference          `json:"spellcasting_ability"`
	Info                []SpellcastingInfo `json:"info"`
}

type SpellcastingInfo struct {
	Name string   `json:"name"`
	Desc []string `json:"desc"`
}

func (c *Class) GetEndpoint() string {
	return "classes/"
}

func (c *Class) Print() {
	fmt.Printf("Class: %s", c.Name)
}
