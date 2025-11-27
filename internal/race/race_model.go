package race

import (
	"fmt"
	"github.com/kwford18/MKDIRagons/internal/reference"
)

// AbilityBonus represents the bonus information for an ability score.
type AbilityBonus struct {
	AbilityScore reference.Reference `json:"ability_score"`
	Bonus        int                 `json:"bonus"`
}

// Race represents a generic D&D 5e race.
type Race struct {
	Index           string                `json:"index"`
	Name            string                `json:"name"`
	Speed           int                   `json:"speed"`
	AbilityBonuses  []AbilityBonus        `json:"ability_bonuses"`
	Age             string                `json:"age"`
	Alignment       string                `json:"alignment"`
	Size            string                `json:"size"`
	SizeDescription string                `json:"size_description"`
	Languages       []reference.Reference `json:"languages"`
	LanguageDesc    string                `json:"language_desc"`
	Traits          []reference.Reference `json:"traits"`
	Subraces        []reference.Reference `json:"subraces"`
	URL             string                `json:"url"`
}

func (r *Race) GetEndpoint() string {
	return "races/"
}

func (r *Race) Print() {
	fmt.Printf("Race: %s\n", r.Name)
}
