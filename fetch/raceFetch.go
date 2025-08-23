package fetch

import "fmt"

// AbilityBonus represents the bonus information for an ability score.
type AbilityBonus struct {
	AbilityScore Reference `json:"ability_score"`
	Bonus        int       `json:"bonus"`
}

// Race represents a generic D&D 5e race.
type Race struct {
	Index                 string         `json:"index"`
	Name                  string         `json:"name"`
	Speed                 int            `json:"speed"`
	AbilityBonuses        []AbilityBonus `json:"ability_bonuses"`
	Age                   string         `json:"age"`
	Alignment             string         `json:"alignment"`
	Size                  string         `json:"size"`
	SizeDescription       string         `json:"size_description"`
	StartingProficiencies []Reference    `json:"starting_proficiencies"`
	Languages             []Reference    `json:"languages"`
	LanguageDesc          string         `json:"language_desc"`
	Traits                []Reference    `json:"traits"`
	Subraces              []Reference    `json:"subraces"`
	URL                   string         `json:"url"`
}

func (e *Race) GetEndpoint() string {
	return "races/" + e.Index
}

func (r *Race) Print() {
	fmt.Printf("Race: %s", r.Name)
}
