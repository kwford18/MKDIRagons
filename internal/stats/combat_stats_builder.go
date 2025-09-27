package stats

import (
	"github.com/kwford18/MKDIRagons/internal/abilities"
	"github.com/kwford18/MKDIRagons/internal/class"
	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/internal/inventory"
	"math/rand/v2"
)

func BuildStats(level int, abilityScores abilities.AbilityScore, class class.Class, armor *inventory.Armor) Stats {
	var HP int
	// Build HP based on level
	for i := 1; i <= level; i++ {
		HP += rand.IntN(class.HitDie) + 1 + abilityScores.Modifier(core.Constitution)
		// fmt.Printf("HP at %d Level: %d\n", i, HP)
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
			AC += abilityScores.Modifier(core.Dexterity)
		}
	} else if class.Name == "Barbarian" {
		// Unarmored Defense for Barbarian
		AC = 10 + abilityScores.Modifier(core.Dexterity) + abilityScores.Modifier(core.Constitution)
	} else if class.Name == "Monk" {
		// Unarmored Defense for Monk
		AC = 10 + abilityScores.Modifier(core.Dexterity) + abilityScores.Modifier(core.Wisdom)
	} else {
		// Default AC
		AC = 10 + abilityScores.Modifier(core.Dexterity)
	}
	return Stats{
		HP:     HP,
		TempHP: 0,
		AC:     AC,
		Speed:  30,
	}
}
