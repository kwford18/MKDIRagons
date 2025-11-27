package spells_test

import (
	"encoding/json"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/reference"
	"github.com/kwford18/MKDIRagons/internal/spells"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpell_GetEndpoint(t *testing.T) {
	s := &spells.Spell{}
	assert.Equal(t, "spells/", s.GetEndpoint(), "Endpoint should be 'spells/'")
}

func TestSpell_JSON_Omitempty(t *testing.T) {
	// Scenario: A simple spell like "Wish" or "Shield" that doesn't have Damage, DC, or AoE.
	// We want to ensure the resulting JSON does NOT contain "null" for these fields,
	// but rather omits the keys entirely.
	spell := spells.Spell{
		Index: "wish",
		Name:  "Wish",
		Level: 9,
		// Explicitly leaving Damage, DC, AreaOfEffect, and Material as nil/empty
	}

	data, err := json.Marshal(spell)
	require.NoError(t, err)

	jsonString := string(data)

	// Assert mandatory fields exist
	assert.Contains(t, jsonString, `"index":"wish"`)
	assert.Contains(t, jsonString, `"name":"Wish"`)

	// Assert optional pointers are OMITTED (not just null)
	assert.NotContains(t, jsonString, `"damage"`, "Field 'damage' should be omitted when nil")
	assert.NotContains(t, jsonString, `"dc"`, "Field 'dc' should be omitted when nil")
	assert.NotContains(t, jsonString, `"area_of_effect"`, "Field 'area_of_effect' should be omitted when nil")
	assert.NotContains(t, jsonString, `"material"`, "Field 'material' should be omitted when empty")
}

func TestSpell_JSON_FullSerialization(t *testing.T) {
	// Scenario: A complex spell like "Fireball" with all optional fields populated.
	spell := spells.Spell{
		Index:    "fireball",
		Name:     "Fireball",
		Level:    3,
		Material: "Guano",
		Damage: &spells.SpellDamage{
			DamageType: reference.Reference{Index: "fire", Name: "Fire"},
			DamageAtSlotLevel: map[string]string{
				"3": "8d6",
			},
		},
		DC: &spells.SpellDC{
			DCType:    reference.Reference{Index: "dex", Name: "DEX"},
			DCSuccess: "half",
		},
		AreaOfEffect: &spells.SpellAreaOfEffect{
			Type: "sphere",
			Size: 20,
		},
	}

	data, err := json.Marshal(spell)
	require.NoError(t, err)

	jsonString := string(data)

	// Assert fields are present
	assert.Contains(t, jsonString, `"material":"Guano"`)

	// Check nested Damage existence
	assert.Contains(t, jsonString, `"damage":`)
	assert.Contains(t, jsonString, `"damage_type":`)
	assert.Contains(t, jsonString, `"8d6"`)

	// Check nested DC existence
	assert.Contains(t, jsonString, `"dc":`)
	assert.Contains(t, jsonString, `"dc_success":"half"`)

	// Check nested AoE existence
	assert.Contains(t, jsonString, `"area_of_effect":`)
	assert.Contains(t, jsonString, `"size":20`)
}
