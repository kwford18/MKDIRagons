package spells

import (
	"fmt"
	"github.com/kwford18/MKDIRagons/internal/reference"
)

type Spell struct {
	HigherLevel   []string              `json:"higher_level"`
	Index         string                `json:"index"`
	Name          string                `json:"name"`
	Desc          []string              `json:"desc"`
	Range         string                `json:"range"`
	Components    []string              `json:"components"`
	Ritual        bool                  `json:"ritual"`
	Duration      string                `json:"duration"`
	Concentration bool                  `json:"concentration"`
	CastingTime   string                `json:"casting_time"`
	Level         int                   `json:"level"`
	AttackType    string                `json:"attack_type"`
	Damage        *SpellDamage          `json:"damage"`
	School        reference.Reference   `json:"school"`
	Classes       []reference.Reference `json:"classes"`
	Subclasses    []reference.Reference `json:"subclasses"`
	URL           string                `json:"url"`
}

type SpellDamage struct {
	DamageType             reference.Reference `json:"damage_type"`
	DamageAtCharacterLevel map[string]string   `json:"damage_at_character_level"`
}

func (s *Spell) GetEndpoint() string {
	return "spells/"
}

func (s *Spell) Print() {
	fmt.Printf("Spell: %s", s.Name)
}
