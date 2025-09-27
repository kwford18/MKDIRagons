package core

// Ability Enum for easier & more consistent lookup and assignment
type Ability int

const (
	Strength Ability = iota
	Dexterity
	Constitution
	Intelligence
	Wisdom
	Charisma
)

// Convert Ability enum value to string representation
func (a Ability) String() string {
	return [...]string{"Strength", "Dexterity", "Constitution", "Intelligence", "Wisdom", "Charisma"}[a]
}
