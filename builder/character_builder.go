package builder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kwford18/MKDIRagons/models"
	"github.com/kwford18/MKDIRagons/templates"
)

func fetchJSON(property models.Fetchable, input string) error {
	baseURL := "https://www.dnd5eapi.co/api/2014/"

	// Format
	endpoint := baseURL + property.GetEndpoint()
	noSpaces := strings.ReplaceAll(input, " ", "-")
	lowercase := strings.ToLower(noSpaces)
	formattedURL := endpoint + strings.ReplaceAll(lowercase, "'", "")

	// fmt.Printf("Formatted URL: %s\n", formatted_url)

	resp, err := http.Get(formattedURL)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(property); err != nil {
		fmt.Printf("Error: %+v\n", err)
		return err
	}

	// fmt.Printf("Fetched: %v\n", property)

	return nil
}

func BuildCharacter(base *templates.TemplateCharacter) (*models.Character, error) {
	var race models.Race
	var class models.Class
	var inventory models.Inventory

	spellbook := initSpellbook(base)

	// Concurrent fetch
	if err := fetchRaceAndClass(base, &race, &class); err != nil {
		return nil, err
	}
	if err := fetchInventory(base, &inventory); err != nil {
		return nil, err
	}
	if err := fetchSpells(base, spellbook); err != nil {
		return nil, err
	}

	// Build ability scores & skills
	ability_scores := buildAbilityScores(base, race)
	skill_list := buildSkillList(base)

	// Build Combat Stats
	var first_armor *models.Armor
	if len(inventory.Armor) > 0 {
		first_armor = &inventory.Armor[0]
	}
	combat_stats := buildStats(base.Level, ability_scores, class, first_armor)

	return &models.Character{
		Name:          base.Name,
		Level:         base.Level,
		Race:          race,
		Class:         class,
		Stats:         combat_stats,
		AbilityScores: ability_scores,
		Skills:        skill_list,
		Proficiencies: base.Proficiencies,
		Inventory:     inventory,
		Spells:        spellbook,
	}, nil
}
