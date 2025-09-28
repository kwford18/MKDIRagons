package stats

import (
	"fmt"
	"math/rand/v2"

	"github.com/kwford18/MKDIRagons/internal/abilities"
	"github.com/kwford18/MKDIRagons/internal/class"
	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/internal/inventory"
)

func averageHP(hitDie int) (int, error) {
	switch hitDie {
	case 6:
		return 4, nil
	case 8:
		return 5, nil
	case 10:
		return 6, nil
	case 12:
		return 7, nil
	default:
		return 0, fmt.Errorf("invalid hit die provided")
	}
}

// BuildStats generates the combat stats of a character and populates + returns the Stats struct with values
func BuildStats(level int, abilityScores abilities.AbilityScore, class class.Class, rollHP bool, armor *inventory.Armor) (Stats, error) {
	var HP int

	// Ability score bonuses to avoid repeat calls
	con_bonus := abilityScores.Modifier(core.Constitution)
	dex_bonus := abilityScores.Modifier(core.Dexterity)

	if rollHP {
		// Build HP based on level
		for i := 1; i <= level; i++ {
			HP += rand.IntN(class.HitDie) + 1 + con_bonus
			// fmt.Printf("HP at %d Level: %d\n", i, HP)
		}
	} else {
		avgHP, err := averageHP(class.HitDie) // Average HP if not rolling
		if err != nil {
			return Stats{}, err
		}
		for i := 1; i <= level; i++ {
			HP += avgHP + con_bonus
		}
	}

	// If character has armor, use that for AC
	// If they are a Barbarian or Monk, use Unarmored Defense
	// Otherwise it is 10 + dex bonus
	// Calculate AC
	var AC int
	if armor != nil {
		// Armor equipped
		AC = armor.ArmorClass.Base
		if armor.ArmorClass.DexBonus {
			AC += dex_bonus
		}
	} else if class.Name == "Barbarian" {
		// Unarmored Defense for Barbarian
		AC = 10 + dex_bonus + con_bonus
	} else if class.Name == "Monk" {
		// Unarmored Defense for Monk
		AC = 10 + dex_bonus + abilityScores.Modifier(core.Wisdom) // Only time calculating wisdom mod, no need to pre-compute
	} else {
		// Default AC
		AC = 10 + dex_bonus
	}
	return Stats{
		HP:     HP,
		TempHP: 0,
		AC:     AC,
		Speed:  30,
	}, nil
}
