package builder

import (
	"math/rand/v2"

	"github.com/kwford18/MKDIRagons/models"
)

func buildStats(level int, ability_scores models.AbilityScore, class models.Class, armor *models.Armor) models.Stats {
	var HP int
	// Build HP based on level
	for i := 1; i <= level; i++ {
		HP += rand.IntN(class.HitDie) + 1 + ability_scores.Modifier(models.Constitution)
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
			AC += ability_scores.Modifier(models.Dexterity)
		}
	} else if class.Name == "Barbarian" {
		// Unarmored Defense for Barbarian
		AC = 10 + ability_scores.Modifier(models.Dexterity) + ability_scores.Modifier(models.Constitution)
	} else if class.Name == "Monk" {
		// Unarmored Defense for Monk
		AC = 10 + ability_scores.Modifier(models.Dexterity) + ability_scores.Modifier(models.Wisdom)
	} else {
		// Default AC
		AC = 10 + ability_scores.Modifier(models.Dexterity)
	}
	return models.Stats{
		HP:     HP,
		TempHP: 0,
		AC:     AC,
		Speed:  30,
	}
}
