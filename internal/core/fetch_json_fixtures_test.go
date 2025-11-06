package core_test

import (
	"net/http"
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

// TestFetchJSON_RealJSON_Wizard tests with actual API response
func TestFetchJSON_RealJSON_Wizard(t *testing.T) {
	mockResponse := core.LoadFixture(t, "wizard.json")

	server := CreateMockServer(t, http.StatusOK, mockResponse)
	defer server.Close()

	class := &class.Class{}
	err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", class, "wizard")
	require.NoError(t, err)

	// Basic fields
	assert.Equal(t, "wizard", class.Index)
	assert.Equal(t, "Wizard", class.Name)
	assert.Equal(t, 6, class.HitDie)
	assert.Equal(t, "/api/2014/classes/wizard", class.URL)

	// Proficiency choices (complex nested structure)
	require.Len(t, class.ProficiencyChoices, 1)
	assert.Equal(t, "Choose two from Arcana, History, Insight, Investigation, Medicine, and Religion",
		class.ProficiencyChoices[0].Desc)
	assert.Equal(t, 2, class.ProficiencyChoices[0].Choose)
	assert.Equal(t, "proficiencies", class.ProficiencyChoices[0].Type)
	assert.Equal(t, "options_array", class.ProficiencyChoices[0].From.OptionSetType)
	assert.Len(t, class.ProficiencyChoices[0].From.Options, 6)

	// Proficiencies
	assert.Len(t, class.Proficiencies, 7)
	assert.Equal(t, "daggers", class.Proficiencies[0].Index)
	assert.Equal(t, "saving-throw-int", class.Proficiencies[5].Index)

	// Saving throws
	require.Len(t, class.SavingThrows, 2)
	assert.Equal(t, "int", class.SavingThrows[0].Index)
	assert.Equal(t, "wis", class.SavingThrows[1].Index)

	// Starting equipment
	require.Len(t, class.StartingEquipment, 1)
	assert.Equal(t, "spellbook", class.StartingEquipment[0].Equipment.Index)
	assert.Equal(t, 1, class.StartingEquipment[0].Quantity)

	// Starting equipment options (very complex nested structure)
	require.Len(t, class.StartingEquipmentOptions, 3)
	firstOption := class.StartingEquipmentOptions[0]
	assert.Equal(t, "(a) a quarterstaff or (b) a dagger", firstOption.Desc)
	assert.Equal(t, 1, firstOption.Choose)
	assert.Len(t, firstOption.From.Options, 2)

	// Spellcasting
	assert.Equal(t, 1, class.Spellcasting.Level)
	assert.Equal(t, "int", class.Spellcasting.SpellcastingAbility.Index)
	require.Len(t, class.Spellcasting.Info, 6)
	assert.Equal(t, "Cantrips", class.Spellcasting.Info[0].Name)
	assert.Contains(t, class.Spellcasting.Info[0].Desc[0], "three cantrips")

	// Multiclassing
	require.Len(t, class.MultiClassing.Prerequisites, 1)
	assert.Equal(t, "int", class.MultiClassing.Prerequisites[0].AbilityScore.Index)
	assert.Equal(t, 13, class.MultiClassing.Prerequisites[0].MinimumScore)

	// Subclasses
	require.Len(t, class.Subclasses, 1)
	assert.Equal(t, "evocation", class.Subclasses[0].Index)
}

func TestFetchJSON_RealJSON_Dwarf(t *testing.T) {
	mockResponse := core.LoadFixture(t, "dwarf.json")

	server := CreateMockServer(t, http.StatusOK, mockResponse)
	defer server.Close()

	race := &race.Race{}
	err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", race, "dwarf")
	require.NoError(t, err)

	// Basic fields
	assert.Equal(t, "dwarf", race.Index)
	assert.Equal(t, "Dwarf", race.Name)
	assert.Equal(t, 25, race.Speed)
	assert.Equal(t, "Medium", race.Size)
	assert.Equal(t, "/api/2014/races/dwarf", race.URL)

	// Ability bonuses
	require.Len(t, race.AbilityBonuses, 1)
	assert.Equal(t, "con", race.AbilityBonuses[0].AbilityScore.Index)
	assert.Equal(t, "CON", race.AbilityBonuses[0].AbilityScore.Name)
	assert.Equal(t, 2, race.AbilityBonuses[0].Bonus)

	// Text descriptions
	assert.Contains(t, race.Alignment, "Most dwarves are lawful")
	assert.Contains(t, race.Age, "350 years")
	assert.Contains(t, race.SizeDescription, "4 and 5 feet tall")
	assert.Contains(t, race.LanguageDesc, "Common and Dwarvish")

	// Languages
	require.Len(t, race.Languages, 2)
	assert.Equal(t, "common", race.Languages[0].Index)
	assert.Equal(t, "dwarvish", race.Languages[1].Index)

	// Traits
	require.Len(t, race.Traits, 5)
	assert.Equal(t, "darkvision", race.Traits[0].Index)
	assert.Equal(t, "dwarven-resilience", race.Traits[1].Index)
	assert.Equal(t, "stonecunning", race.Traits[2].Index)

	// Subraces
	require.Len(t, race.Subraces, 1)
	assert.Equal(t, "hill-dwarf", race.Subraces[0].Index)
}

func TestFetchJSON_RealJSON_Fireball(t *testing.T) {
	mockResponse := core.LoadFixture(t, "fireball.json")

	server := CreateMockServer(t, http.StatusOK, mockResponse)
	defer server.Close()

	spell := &spells.Spell{}
	err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", spell, "fireball")
	require.NoError(t, err)

	// Basic fields
	assert.Equal(t, "fireball", spell.Index)
	assert.Equal(t, "Fireball", spell.Name)
	assert.Equal(t, 3, spell.Level)
	assert.Equal(t, "/api/2014/spells/fireball", spell.URL)

	// Spell details
	assert.Equal(t, "150 feet", spell.Range)
	assert.Equal(t, "Instantaneous", spell.Duration)
	assert.Equal(t, "1 action", spell.CastingTime)
	assert.False(t, spell.Ritual)
	assert.False(t, spell.Concentration)

	// Components
	require.Len(t, spell.Components, 3)
	assert.Contains(t, spell.Components, "V")
	assert.Contains(t, spell.Components, "S")
	assert.Contains(t, spell.Components, "M")

	// Description
	require.Len(t, spell.Desc, 2)
	assert.Contains(t, spell.Desc[0], "bright streak flashes")
	assert.Contains(t, spell.Desc[1], "fire spreads around corners")

	// Higher level casting
	require.Len(t, spell.HigherLevel, 1)
	assert.Contains(t, spell.HigherLevel[0], "4th level or higher")

	// Damage (Note: Real API uses damage_at_slot_level, not damage_at_character_level)
	require.NotNil(t, spell.Damage)
	assert.Equal(t, "fire", spell.Damage.DamageType.Index)
	// Your struct might need updating if it expects DamageAtCharacterLevel
	// The real API has DamageAtSlotLevel for spells

	// School
	assert.Equal(t, "evocation", spell.School.Index)
	assert.Equal(t, "Evocation", spell.School.Name)

	// Classes
	require.Len(t, spell.Classes, 2)
	assert.Equal(t, "sorcerer", spell.Classes[0].Index)
	assert.Equal(t, "wizard", spell.Classes[1].Index)

	// Subclasses
	require.Len(t, spell.Subclasses, 2)
	assert.Equal(t, "lore", spell.Subclasses[0].Index)
	assert.Equal(t, "fiend", spell.Subclasses[1].Index)
}

func TestFetchJSON_RealJSON_PaddedArmor(t *testing.T) {
	mockResponse := core.LoadFixture(t, "padded-armor.json")

	server := CreateMockServer(t, http.StatusOK, mockResponse)
	defer server.Close()

	armor := &inventory.Armor{}
	err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", armor, "padded-armor")
	require.NoError(t, err)

	// Basic fields
	assert.Equal(t, "padded-armor", armor.Index)
	assert.Equal(t, "Padded Armor", armor.Name)
	assert.Equal(t, 8, armor.Weight)
	assert.Equal(t, "/api/2014/equipment/padded-armor", armor.URL)

	// Armor specific
	assert.Equal(t, "Light", armor.ArmorCategory)
	assert.Equal(t, 11, armor.ArmorClass.Base)
	assert.True(t, armor.ArmorClass.DexBonus)
	assert.Equal(t, 0, armor.StrMinimum)
	assert.True(t, armor.StealthDisadvantage)

	// Cost
	assert.Equal(t, 5, armor.Cost.Quantity)
	assert.Equal(t, "gp", armor.Cost.Unit)

	// Equipment category
	assert.Equal(t, "armor", armor.EquipmentCategory.Index)
	assert.Equal(t, "Armor", armor.EquipmentCategory.Name)

	// Empty arrays (verify they don't cause issues)
	assert.Empty(t, armor.Desc)
	assert.Empty(t, armor.Special)
	assert.Empty(t, armor.Contents)
	assert.Empty(t, armor.Properties)
}

// Table-driven test using all fixtures
func TestFetchJSON_AllRealJSON(t *testing.T) {
	tests := []struct {
		name         string
		fixture      string
		fetchable    reference.Fetchable
		input        string
		validateFunc func(t *testing.T, f reference.Fetchable)
	}{
		{
			name:      "Wizard Class",
			fixture:   "wizard.json",
			fetchable: &class.Class{},
			input:     "wizard",
			validateFunc: func(t *testing.T, f reference.Fetchable) {
				class := f.(*class.Class)
				assert.Equal(t, "wizard", class.Index)
				assert.Equal(t, 6, class.HitDie)
				assert.Len(t, class.ProficiencyChoices, 1)
				assert.Len(t, class.Spellcasting.Info, 6)
			},
		},
		{
			name:      "Dwarf Race",
			fixture:   "dwarf.json",
			fetchable: &race.Race{},
			input:     "dwarf",
			validateFunc: func(t *testing.T, f reference.Fetchable) {
				race := f.(*race.Race)
				assert.Equal(t, "dwarf", race.Index)
				assert.Equal(t, 25, race.Speed)
				assert.Len(t, race.Traits, 5)
				assert.Len(t, race.Languages, 2)
			},
		},
		{
			name:      "Fireball Spell",
			fixture:   "fireball.json",
			fetchable: &spells.Spell{},
			input:     "fireball",
			validateFunc: func(t *testing.T, f reference.Fetchable) {
				spell := f.(*spells.Spell)
				assert.Equal(t, "fireball", spell.Index)
				assert.Equal(t, 3, spell.Level)
				assert.Len(t, spell.Classes, 2)
				assert.NotNil(t, spell.Damage)
			},
		},
		{
			name:      "Padded Armor",
			fixture:   "padded-armor.json",
			fetchable: &inventory.Armor{},
			input:     "padded-armor",
			validateFunc: func(t *testing.T, f reference.Fetchable) {
				armor := f.(*inventory.Armor)
				assert.Equal(t, "padded-armor", armor.Index)
				assert.Equal(t, "Light", armor.ArmorCategory)
				assert.True(t, armor.StealthDisadvantage)
				assert.True(t, armor.ArmorClass.DexBonus)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockResponse := core.LoadFixture(t, tt.fixture)
			server := CreateMockServer(t, http.StatusOK, mockResponse)
			defer server.Close()

			err := core.FetchJSONWithClient(http.DefaultClient, server.URL+"/", tt.fetchable, tt.input)
			require.NoError(t, err)

			tt.validateFunc(t, tt.fetchable)
		})
	}
}

// Test direct unmarshaling from fixture (no HTTP involved)
// This tests that your structs can properly unmarshal the real API JSON
func TestUnmarshal_RealJSON_Direct(t *testing.T) {
	t.Run("Wizard", func(t *testing.T) {
		class := &class.Class{}
		core.LoadFixtureInto(t, "wizard.json", class)

		assert.Equal(t, "wizard", class.Index)
		assert.Equal(t, 6, class.HitDie)
		assert.NotEmpty(t, class.ProficiencyChoices)
		assert.Len(t, class.Spellcasting.Info, 6)
	})

	t.Run("Dwarf", func(t *testing.T) {
		race := &race.Race{}
		core.LoadFixtureInto(t, "dwarf.json", race)

		assert.Equal(t, "dwarf", race.Index)
		assert.Equal(t, 25, race.Speed)
		assert.Len(t, race.Traits, 5)
		assert.Len(t, race.Languages, 2)
	})

	t.Run("Fireball", func(t *testing.T) {
		spell := &spells.Spell{}
		core.LoadFixtureInto(t, "fireball.json", spell)

		assert.Equal(t, "fireball", spell.Index)
		assert.Equal(t, 3, spell.Level)
		assert.NotNil(t, spell.Damage)
		assert.Len(t, spell.Classes, 2)
	})

	t.Run("Padded Armor", func(t *testing.T) {
		armor := &inventory.Armor{}
		core.LoadFixtureInto(t, "padded-armor.json", armor)

		assert.Equal(t, "padded-armor", armor.Index)
		assert.Equal(t, "Light", armor.ArmorCategory)
		assert.True(t, armor.StealthDisadvantage)
	})
}
