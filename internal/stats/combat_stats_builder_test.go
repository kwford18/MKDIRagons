package stats_test

import (
	"github.com/kwford18/MKDIRagons/internal/stats"
	"testing"

	"github.com/kwford18/MKDIRagons/internal/abilities"
	"github.com/kwford18/MKDIRagons/internal/class"
	"github.com/kwford18/MKDIRagons/internal/inventory"
	"github.com/stretchr/testify/assert"
)

// setupAbilities is a helper to create an AbilityScores struct.
// Adjust the field initialization below if your AbilityScores struct uses a map or different field names.
func setupAbilities(str, dex, con, intel, wis, cha int) abilities.AbilityScores {
	return abilities.AbilityScores{
		Strength:     str,
		Dexterity:    dex,
		Constitution: con,
		Intelligence: intel,
		Wisdom:       wis,
		Charisma:     cha,
	}
}

// setupArmor is a helper to create an Armor struct for testing.
func setupArmor(baseAC int, dexBonus bool) *inventory.Armor {
	return &inventory.Armor{
		ArmorClass: inventory.ArmorClass{ // Assuming inventory package has this nested struct based on usage
			Base:     baseAC,
			DexBonus: dexBonus,
		},
	}
}

func TestBuildStats_HP_Average(t *testing.T) {
	// Setup: Level 3 Fighter (d10), Con 14 (+2), Dex 10 (+0)
	scores := setupAbilities(10, 10, 14, 10, 10, 10)
	cls := class.Class{Name: "Fighter", HitDie: 10}
	level := 3

	// Expected HP Calculation:
	// HitDie 10 -> Average is 6 (from switch case)
	// Modifier is +2
	// Per Level: 6 + 2 = 8
	// Total Level 3: 8 * 3 = 24

	testStats, err := stats.BuildStats(level, scores, cls, false, nil)

	assert.NoError(t, err)
	assert.Equal(t, 24, testStats.HP)
}

func TestBuildStats_HP_InvalidHitDie(t *testing.T) {
	scores := setupAbilities(10, 10, 10, 10, 10, 10)
	cls := class.Class{Name: "Broken", HitDie: 20} // 20 is not in switch case

	_, err := stats.BuildStats(1, scores, cls, false, nil)

	assert.Error(t, err)
	assert.Equal(t, "invalid hit die provided", err.Error())
}

func TestBuildStats_HP_Rolled_Range(t *testing.T) {
	// Since rand is non-deterministic, we test that the result falls within possible bounds.
	// Level 10, HitDie 6, Con 10 (+0)
	// Min HP (all 1s): 1 * 10 = 10
	// Max HP (all 6s): 6 * 10 = 60

	scores := setupAbilities(10, 10, 10, 10, 10, 10)
	cls := class.Class{Name: "Wizard", HitDie: 6}
	level := 10

	testStats, err := stats.BuildStats(level, scores, cls, true, nil)

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, testStats.HP, 10)
	assert.LessOrEqual(t, testStats.HP, 60)
}

func TestBuildStats_AC_Default(t *testing.T) {
	// No Armor, Standard Class
	// AC = 10 + DexMod
	// Dex 14 (+2) -> AC 12

	scores := setupAbilities(10, 14, 10, 10, 10, 10)
	cls := class.Class{Name: "Fighter", HitDie: 10}

	testStats, _ := stats.BuildStats(1, scores, cls, false, nil)
	assert.Equal(t, 12, testStats.AC)
}

func TestBuildStats_AC_Barbarian_Unarmored(t *testing.T) {
	// Barbarian Unarmored Defense: 10 + Dex + Con
	// Dex 14 (+2), Con 16 (+3) -> AC 15

	scores := setupAbilities(10, 14, 16, 10, 10, 10)
	cls := class.Class{Name: "Barbarian", HitDie: 12}

	testStats, _ := stats.BuildStats(1, scores, cls, false, nil)
	assert.Equal(t, 15, testStats.AC)
}

func TestBuildStats_AC_Monk_Unarmored(t *testing.T) {
	// Monk Unarmored Defense: 10 + Dex + Wis
	// Dex 14 (+2), Wis 16 (+3) -> AC 15

	scores := setupAbilities(10, 14, 10, 10, 16, 10)
	cls := class.Class{Name: "Monk", HitDie: 8}

	testStats, _ := stats.BuildStats(1, scores, cls, false, nil)
	assert.Equal(t, 15, testStats.AC)
}

func TestBuildStats_AC_WithArmor_NoDex(t *testing.T) {
	// Heavy Armor (e.g., Plate): Base 18, No Dex Bonus
	// Dex 14 (+2) should be ignored

	scores := setupAbilities(10, 14, 10, 10, 10, 10)
	cls := class.Class{Name: "Fighter", HitDie: 10}
	armor := setupArmor(18, false)

	testStats, _ := stats.BuildStats(1, scores, cls, false, armor)
	assert.Equal(t, 18, testStats.AC)
}

func TestBuildStats_AC_WithArmor_WithDex(t *testing.T) {
	// Light Armor (e.g., Leather): Base 11, + Dex Bonus
	// Dex 14 (+2) -> AC 13

	scores := setupAbilities(10, 14, 10, 10, 10, 10)
	cls := class.Class{Name: "Rogue", HitDie: 8}
	armor := setupArmor(11, true)

	testStats, _ := stats.BuildStats(1, scores, cls, false, armor)
	assert.Equal(t, 13, testStats.AC)
}

func TestBuildStats_Speed(t *testing.T) {
	// Verify default speed is 30
	scores := setupAbilities(10, 10, 10, 10, 10, 10)
	cls := class.Class{Name: "Fighter", HitDie: 10}

	testStats, _ := stats.BuildStats(1, scores, cls, false, nil)
	assert.Equal(t, 30, testStats.Speed)
}
