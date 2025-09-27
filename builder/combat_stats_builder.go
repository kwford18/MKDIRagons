package builder

import (
	"math/rand/v2"

	"github.com/kwford18/MKDIRagons/models"
)

func buildHP(level, constitution int, class models.Class) int {
	var HP int
	for i := 1; i <= level; i++ {
		HP += rand.IntN(class.HitDie) + 1 + constitution
	}

	return HP
}

// func buildAC(dexterity int, armor *models.Equipment) int {
// 	var AC int

// 	if armor != nil {
// 		armor.
// 	}

// 	return AC
// }
