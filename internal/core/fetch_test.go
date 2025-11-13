package core_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/class"
	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/internal/inventory"
	"github.com/kwford18/MKDIRagons/internal/race"
	"github.com/kwford18/MKDIRagons/internal/reference"
	"github.com/kwford18/MKDIRagons/internal/spells"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test helper to create a mock server with custom response
func CreateMockServer(t *testing.T, statusCode int, response interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if response != nil {
			err := json.NewEncoder(w).Encode(response)
			require.NoError(t, err)
		}
	}))
}

func TestFetchJSON_Class_Success(t *testing.T) {
	// Create realistic mock response based on actual API
	mockResponse := map[string]interface{}{
		"index":   "wizard",
		"name":    "Wizard",
		"hit_die": 6,
		"url":     "/api/2014/classes/wizard",
		"proficiencies": []map[string]string{
			{"index": "daggers", "name": "Daggers", "url": "/api/2014/proficiencies/daggers"},
			{"index": "darts", "name": "Darts", "url": "/api/2014/proficiencies/darts"},
		},
		"saving_throws": []map[string]string{
			{"index": "int", "name": "INT", "url": "/api/2014/ability-scores/int"},
			{"index": "wis", "name": "WIS", "url": "/api/2014/ability-scores/wis"},
		},
		"starting_equipment": []map[string]interface{}{
			{
				"equipment": map[string]string{
					"index": "spellbook",
					"name":  "Spellbook",
					"url":   "/api/2014/equipment/spellbook",
				},
				"quantity": 1,
			},
		},
		"spellcasting": map[string]interface{}{
			"level": 1,
			"spellcasting_ability": map[string]string{
				"index": "int",
				"name":  "INT",
				"url":   "/api/2014/ability-scores/int",
			},
			"info": []map[string]interface{}{
				{
					"name": "Cantrips",
					"desc": []string{"At 1st level, you know three cantrips..."},
				},
			},
		},
		"class_levels": "/api/2014/classes/wizard/levels",
		"spells":       "/api/2014/classes/wizard/spells",
	}

	server := CreateMockServer(t, http.StatusOK, mockResponse)
	defer server.Close()

	testClass := &class.Class{}
	err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", testClass, "wizard")

	require.NoError(t, err)
	assert.Equal(t, "wizard", testClass.Index)
	assert.Equal(t, "Wizard", testClass.Name)
	assert.Equal(t, 6, testClass.HitDie)
	assert.Len(t, testClass.SavingThrows, 2)
	assert.Equal(t, 1, testClass.Spellcasting.Level)
	assert.Equal(t, "int", testClass.Spellcasting.SpellcastingAbility.Index)
}

func TestFetchJSON_Race_Success(t *testing.T) {
	mockResponse := map[string]interface{}{
		"index": "elf",
		"name":  "Elf",
		"speed": 30,
		"size":  "Medium",
		"url":   "/api/races/elf",
		"ability_bonuses": []map[string]interface{}{
			{
				"ability_score": map[string]string{
					"index": "dex",
					"name":  "DEX",
					"url":   "/api/ability-scores/dex",
				},
				"bonus": 2,
			},
		},
		"languages": []map[string]string{
			{"index": "common", "name": "Common", "url": "/api/languages/common"},
			{"index": "elvish", "name": "Elvish", "url": "/api/languages/elvish"},
		},
	}

	server := CreateMockServer(t, http.StatusOK, mockResponse)
	defer server.Close()

	testRace := &race.Race{}
	err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", testRace, "elf")

	require.NoError(t, err)
	assert.Equal(t, "elf", testRace.Index)
	assert.Equal(t, "Elf", testRace.Name)
	assert.Equal(t, 30, testRace.Speed)
	assert.Equal(t, "Medium", testRace.Size)
	assert.Len(t, testRace.AbilityBonuses, 1)
	assert.Equal(t, 2, testRace.AbilityBonuses[0].Bonus)
}

func TestFetchJSON_Spell_Success(t *testing.T) {
	mockResponse := map[string]interface{}{
		"index":        "fireball",
		"name":         "Fireball",
		"level":        3,
		"range":        "150 feet",
		"components":   []string{"V", "S", "M"},
		"duration":     "Instantaneous",
		"casting_time": "1 action",
		"url":          "/api/spells/fireball",
		"desc": []string{
			"A bright streak flashes from your pointing finger...",
		},
		"school": map[string]string{
			"index": "evocation",
			"name":  "Evocation",
			"url":   "/api/magic-schools/evocation",
		},
		"classes": []map[string]string{
			{"index": "sorcerer", "name": "Sorcerer", "url": "/api/classes/sorcerer"},
			{"index": "wizard", "name": "Wizard", "url": "/api/classes/wizard"},
		},
	}

	server := CreateMockServer(t, http.StatusOK, mockResponse)
	defer server.Close()

	spell := &spells.Spell{}
	err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", spell, "fireball")

	require.NoError(t, err)
	assert.Equal(t, "fireball", spell.Index)
	assert.Equal(t, "Fireball", spell.Name)
	assert.Equal(t, 3, spell.Level)
	assert.NotEmpty(t, spell.Desc)
	assert.Equal(t, "evocation", spell.School.Index)
	assert.Len(t, spell.Classes, 2)
}

func TestFetchJSON_Item_Success(t *testing.T) {
	mockResponse := map[string]interface{}{
		"index":  "abacus",
		"name":   "Abacus",
		"weight": 2,
		"url":    "/api/equipment/abacus",
		"equipment_category": map[string]string{
			"index": "adventuring-gear",
			"name":  "Adventuring Gear",
			"url":   "/api/equipment-categories/adventuring-gear",
		},
		"cost": map[string]interface{}{
			"quantity": 2,
			"unit":     "gp",
		},
	}

	server := CreateMockServer(t, http.StatusOK, mockResponse)
	defer server.Close()

	item := &inventory.Item{}
	err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", item, "abacus")

	require.NoError(t, err)
	assert.Equal(t, "abacus", item.Index)
	assert.Equal(t, "Abacus", item.Name)
	assert.Equal(t, 2, item.Cost.Quantity)
	assert.Equal(t, "gp", item.Cost.Unit)
}

func TestFetchJSON_Armor_Success(t *testing.T) {
	mockResponse := map[string]interface{}{
		"index":                "chain-mail",
		"name":                 "Chain Mail",
		"armor_category":       "Heavy",
		"str_minimum":          13,
		"stealth_disadvantage": true,
		"weight":               55,
		"url":                  "/api/equipment/chain-mail",
		"equipment_category": map[string]string{
			"index": "armor",
			"name":  "Armor",
			"url":   "/api/equipment-categories/armor",
		},
		"armor_class": map[string]interface{}{
			"base":      16,
			"dex_bonus": false,
		},
		"cost": map[string]interface{}{
			"quantity": 75,
			"unit":     "gp",
		},
	}

	server := CreateMockServer(t, http.StatusOK, mockResponse)
	defer server.Close()

	armor := &inventory.Armor{}
	err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", armor, "chain-mail")

	require.NoError(t, err)
	assert.Equal(t, "chain-mail", armor.Index)
	assert.Equal(t, "Chain Mail", armor.Name)
	assert.Equal(t, "Heavy", armor.ArmorCategory)
	assert.Equal(t, 16, armor.ArmorClass.Base)
	assert.False(t, armor.ArmorClass.DexBonus)
	assert.True(t, armor.StealthDisadvantage)
}

func TestFetchJSON_Weapon_Success(t *testing.T) {
	mockResponse := map[string]interface{}{
		"index":           "longsword",
		"name":            "Longsword",
		"weapon_category": "Martial",
		"weapon_range":    "Melee",
		"category_range":  "Martial Melee",
		"weight":          3,
		"url":             "/api/equipment/longsword",
		"equipment_category": map[string]string{
			"index": "weapon",
			"name":  "Weapon",
			"url":   "/api/equipment-categories/weapon",
		},
		"damage": map[string]interface{}{
			"damage_dice": "1d8",
			"damage_type": map[string]string{
				"index": "slashing",
				"name":  "Slashing",
				"url":   "/api/damage-types/slashing",
			},
		},
		"range": map[string]interface{}{
			"normal": 5,
			"long":   5,
		},
		"cost": map[string]interface{}{
			"quantity": 15,
			"unit":     "gp",
		},
	}

	server := CreateMockServer(t, http.StatusOK, mockResponse)
	defer server.Close()

	weapon := &inventory.Weapon{}
	err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", weapon, "longsword")

	require.NoError(t, err)
	assert.Equal(t, "longsword", weapon.Index)
	assert.Equal(t, "Longsword", weapon.Name)
	assert.Equal(t, "Martial", weapon.WeaponCategory)
	assert.Equal(t, "1d8", weapon.Damage.DamageDice)
	assert.Equal(t, "slashing", weapon.Damage.DamageType.Index)
}

func TestFetchJSON_InputFormatting(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedURL string
	}{
		{
			name:        "Spaces to dashes",
			input:       "magic missile",
			expectedURL: "/spells/magic-missile",
		},
		{
			name:        "Uppercase to lowercase",
			input:       "FIREBALL",
			expectedURL: "/spells/fireball",
		},
		{
			name:        "Apostrophes removed",
			input:       "tasha's hideous laughter",
			expectedURL: "/spells/tashas-hideous-laughter",
		},
		{
			name:        "Mixed formatting",
			input:       "Hunter's Mark",
			expectedURL: "/spells/hunters-mark",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestedPath string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestedPath = r.URL.Path
				w.WriteHeader(http.StatusOK)
				err := json.NewEncoder(w).Encode(map[string]string{
					"index": "test",
					"name":  "Test",
				})
				if err != nil {
					return
				}
			}))
			defer server.Close()

			spell := &spells.Spell{}
			_ = core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", spell, tt.input)

			assert.Equal(t, tt.expectedURL, requestedPath, "URL should be properly formatted")
		})
	}
}

func TestFetchJSON_NotFound(t *testing.T) {
	server := CreateMockServer(t, http.StatusNotFound, map[string]string{
		"error": "Not found",
	})
	defer server.Close()

	testClass := &class.Class{}
	err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", testClass, "nonexistent")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "status 404")
}

func TestFetchJSON_ServerError(t *testing.T) {
	server := CreateMockServer(t, http.StatusInternalServerError, map[string]string{
		"error": "Internal server error",
	})
	defer server.Close()

	spell := &spells.Spell{}
	err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", spell, "fireball")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "status 500")
}

func TestFetchJSON_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("invalid json{{{"))
		if err != nil {
			return
		}
	}))
	defer server.Close()

	testRace := &race.Race{}
	err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", testRace, "elf")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode JSON")
}

func TestFetchJSON_NetworkError(t *testing.T) {
	// Use an invalid URL to simulate network error
	testClass := &class.Class{}
	err := core.FetchJSONWithClient(http.DefaultClient, "http://invalid-url-that-does-not-exist.local/", testClass, "wizard")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to make request")
}

func TestFetchJSON_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("{}"))
		if err != nil {
			return
		}
	}))
	defer server.Close()

	item := &inventory.Item{}
	err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", item, "test")

	// Should not error, but fields should be empty/zero values
	require.NoError(t, err)
	assert.Empty(t, item.Index)
	assert.Empty(t, item.Name)
}

func TestFetchJSON_AllTypes_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		fetchable    reference.Fetchable
		input        string
		mockResponse map[string]interface{}
		validate     func(t *testing.T, f reference.Fetchable)
	}{
		{
			name:      "Class",
			fetchable: &class.Class{},
			input:     "wizard",
			mockResponse: map[string]interface{}{
				"index":   "wizard",
				"name":    "Wizard",
				"hit_die": 6,
			},
			validate: func(t *testing.T, f reference.Fetchable) {
				testClass := f.(*class.Class)
				assert.Equal(t, "wizard", testClass.Index)
				assert.Equal(t, 6, testClass.HitDie)
			},
		},
		{
			name:      "Race",
			fetchable: &race.Race{},
			input:     "elf",
			mockResponse: map[string]interface{}{
				"index": "elf",
				"name":  "Elf",
				"speed": 30,
				"size":  "Medium",
			},
			validate: func(t *testing.T, f reference.Fetchable) {
				testRace := f.(*race.Race)
				assert.Equal(t, "elf", testRace.Index)
				assert.Equal(t, 30, testRace.Speed)
			},
		},
		{
			name:      "Spell",
			fetchable: &spells.Spell{},
			input:     "fireball",
			mockResponse: map[string]interface{}{
				"index": "fireball",
				"name":  "Fireball",
				"level": 3,
			},
			validate: func(t *testing.T, f reference.Fetchable) {
				spell := f.(*spells.Spell)
				assert.Equal(t, "fireball", spell.Index)
				assert.Equal(t, 3, spell.Level)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := CreateMockServer(t, http.StatusOK, tt.mockResponse)
			defer server.Close()

			err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", tt.fetchable, tt.input)
			require.NoError(t, err)

			tt.validate(t, tt.fetchable)
		})
	}
}
