package spells

import (
	"encoding/json"
	"fmt"

	"bytes"

	"github.com/kwford18/MKDIRagons/internal/reference"
)

// UnmarshalJSON handles the 5e API's inconsistent damage field types.
func (s *Spell) UnmarshalJSON(data []byte) error {
	type Alias Spell
	aux := &struct {
		Damage json.RawMessage `json:"damage,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if len(aux.Damage) == 0 {
		s.Damage = nil
		return nil
	}

	trimmed := bytes.TrimSpace(aux.Damage)

	// Scenario 1: Empty array (e.g., Alarm) or null
	if bytes.Equal(trimmed, []byte("[]")) || bytes.Equal(trimmed, []byte("null")) {
		s.Damage = nil
		return nil
	}

	// Scenario 2: Array of objects (e.g., Acid Splash)
	if bytes.HasPrefix(trimmed, []byte("[")) {
		var dmgArray []SpellDamage
		if err := json.Unmarshal(trimmed, &dmgArray); err != nil {
			return err
		}
		if len(dmgArray) > 0 {
			s.Damage = &dmgArray[0] // Assign the first damage object
		}
		return nil
	}

	// Scenario 3: Single object (e.g., Fireball)
	if bytes.HasPrefix(trimmed, []byte("{")) {
		s.Damage = &SpellDamage{}
		if err := json.Unmarshal(trimmed, s.Damage); err != nil {
			return err
		}
		return nil
	}

	return nil
}

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
	DamageType             reference.Reference `json:"damage_type"`
	DamageAtSlotLevel      map[string]string   `json:"damage_at_slot_level,omitempty"`
	DamageAtCharacterLevel map[string]string   `json:"damage_at_character_level,omitempty"`
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
