package abilities

import (
	"fmt"
	"github.com/kwford18/MKDIRagons/internal/core"
	"math"
)

// AbilityScores Struct type to represent the player's ability scores
type AbilityScores struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

// Modifier takes an ability and returns the modifier
func (ab *AbilityScores) Modifier(a core.Ability) int {
	// Get the value for a given ability score
	var score int
	switch a {
	case core.Strength:
		score = ab.Strength
	case core.Dexterity:
		score = ab.Dexterity
	case core.Constitution:
		score = ab.Constitution
	case core.Intelligence:
		score = ab.Intelligence
	case core.Wisdom:
		score = ab.Wisdom
	case core.Charisma:
		score = ab.Charisma
	}

	// Formula to calculate the bonus: (ability score value - 10)/2, rounded down
	return int(math.Floor(float64(score-10) / 2.0))
}

// Print for Fetchable interface methods
func (ab *AbilityScores) Print() {
	fmt.Printf("    - Strength:     %d\n", ab.Strength)
	fmt.Printf("    - Dexterity:    %d\n", ab.Dexterity)
	fmt.Printf("    - Constitution: %d\n", ab.Constitution)
	fmt.Printf("    - Wisdom:       %d\n", ab.Wisdom)
	fmt.Printf("    - Intelligence: %d\n", ab.Intelligence)
	fmt.Printf("    - Charisma:     %d\n", ab.Charisma)
	fmt.Println()
}

func (ab *AbilityScores) GetEndpoint() string {
	return "ability-scores/"
}
