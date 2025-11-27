package spells

import (
	"fmt"
	"github.com/kwford18/MKDIRagons/internal/reference"
)

type Spell struct {
	Index         string   `json:"index"`
	Name          string   `json:"name"`
	Desc          []string `json:"desc"`
	HigherLevel   []string `json:"higher_level,omitempty"`
	Range         string   `json:"range"`
	Components    []string `json:"components"`
	Material      string   `json:"material,omitempty"`
	Ritual        bool     `json:"ritual"`
	Duration      string   `json:"duration"`
	Concentration bool     `json:"concentration"`
	CastingTime   string   `json:"casting_time"`
	Level         int      `json:"level"`
	AttackType    string   `json:"attack_type,omitempty"`

	// Complex nested objects (pointers handle null/missing keys)
	Damage       *SpellDamage       `json:"damage,omitempty"`
	DC           *SpellDC           `json:"dc,omitempty"`
	AreaOfEffect *SpellAreaOfEffect `json:"area_of_effect,omitempty"`

	School     reference.Reference   `json:"school"`
	Classes    []reference.Reference `json:"classes"`
	Subclasses []reference.Reference `json:"subclasses"`
	URL        string                `json:"url"`
}

type SpellDamage struct {
	DamageType        reference.Reference `json:"damage_type"`
	DamageAtSlotLevel map[string]string   `json:"damage_at_slot_level"`
}

type SpellDC struct {
	DCType    reference.Reference `json:"dc_type"`
	DCSuccess string              `json:"dc_success"`
}

type SpellAreaOfEffect struct {
	Type string `json:"type"`
	Size int    `json:"size"`
}

func (s *Spell) GetEndpoint() string {
	return "spells/"
}

func (s *Spell) Print() {
	fmt.Printf("Spell: %s", s.Name)
}
