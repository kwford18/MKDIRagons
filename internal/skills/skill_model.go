package skills

import (
	"fmt"
	"github.com/kwford18/MKDIRagons/internal/core"
)

// Skill single skill
type Skill struct {
	Name       string
	Bonus      int
	Ability    core.Ability
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
