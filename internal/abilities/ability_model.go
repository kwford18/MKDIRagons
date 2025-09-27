package abilities

import (
	"fmt"
	"github.com/kwford18/MKDIRagons/internal/core"
)

// AbilityScore Struct type to represent the player's ability scores
type AbilityScore struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

// Modifier takes an ability and returns the modifier
func (ab *AbilityScore) Modifier(a core.Ability) int {
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
	return (score - 10) / 2
}

// Print for Fetchable interface methods
func (ab *AbilityScore) Print() {
	fmt.Println("Ability Scores:")
	fmt.Printf("    - Strength:     %d\n", ab.Strength)
	fmt.Printf("    - Dexterity:    %d\n", ab.Dexterity)
	fmt.Printf("    - Constitution: %d\n", ab.Constitution)
	fmt.Printf("    - Wisdom:       %d\n", ab.Wisdom)
	fmt.Printf("    - Intelligence: %d\n", ab.Intelligence)
	fmt.Printf("    - Charisma:     %d\n", ab.Charisma)
	fmt.Println()
}

func (ab *AbilityScore) GetEndpoint() string {
	return "ability-scores/"
}
