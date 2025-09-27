package models

import "fmt"

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

// AbilityScore Struct type to represent the player's ability scores
type AbilityScore struct {
	Strength     int
	Dexterity    int
	Constitution int
	Intelligence int
	Wisdom       int
	Charisma     int
}

// Skill single skill
type Skill struct {
	Name       string
	Bonus      int
	Ability    Ability
	Proficient bool
	Expertise  bool
}

// SkillList for each of type Skill
type SkillList struct {
	Athletics      Skill
	Acrobatics     Skill
	SleightOfHand  Skill
	Stealth        Skill
	Arcana         Skill
	History        Skill
	Investigation  Skill
	Nature         Skill
	Religion       Skill
	AnimalHandling Skill
	Insight        Skill
	Medicine       Skill
	Perception     Skill
	Survival       Skill
	Deception      Skill
	Intimidation   Skill
	Performance    Skill
	Persuasion     Skill
}

// Convert Ability enum value to string representation
func (a Ability) String() string {
	return [...]string{"Strength", "Dexterity", "Constitution", "Intelligence", "Wisdom", "Charisma"}[a]
}

// Modifier takes an ability and returns the modifier
func (ab *AbilityScore) Modifier(a Ability) int {
	var score int
	switch a {
	case Strength:
		score = ab.Strength
	case Dexterity:
		score = ab.Dexterity
	case Constitution:
		score = ab.Constitution
	case Intelligence:
		score = ab.Intelligence
	case Wisdom:
		score = ab.Wisdom
	case Charisma:
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

func (s *Skill) GetEndpoint() string {
	return "skills/" + s.Name
}

func (s *Skill) Print() {
	fmt.Printf("	- Name: %s\n", s.Name)
	fmt.Printf("	- Value: %d\n", s.Bonus)
	fmt.Printf("	- Ability: %s\n", s.Ability)
	fmt.Printf("	- Proficient: %t\n", s.Proficient)
	fmt.Printf("	- Expertise: %t\n", s.Expertise)
}

func (sl *SkillList) GetEndpoint() string {
	return "skills/"
}

func (sl *SkillList) Print() {
	skills := []Skill{
		sl.Athletics,
		sl.Acrobatics,
		sl.SleightOfHand,
		sl.Stealth,
		sl.Arcana,
		sl.History,
		sl.Investigation,
		sl.Nature,
		sl.Religion,
		sl.AnimalHandling,
		sl.Insight,
		sl.Medicine,
		sl.Perception,
		sl.Survival,
		sl.Deception,
		sl.Intimidation,
		sl.Performance,
		sl.Persuasion,
	}

	for _, skill := range skills {
		fmt.Printf("Skill: %s, Bonus: %d\n", skill.Name, skill.Bonus)
	}
}
