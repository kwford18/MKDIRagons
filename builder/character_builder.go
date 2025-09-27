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
	if err := fetchRace(base, &race); err != nil {
		return nil, err
	}
	if err := fetchClass(base, &class); err != nil {
		return nil, err
	}
	if err := fetchInventory(base, &inventory); err != nil {
		return nil, err
	}
	if err := fetchSpells(base, spellbook); err != nil {
		return nil, err
	}

	// Build ability scores & skills
	abilityScores := buildAbilityScores(base, race)
	skillList := buildSkillList(base)

	// Build Combat Stats
	var firstArmor *models.Armor
	if len(inventory.Armor) > 0 {
		firstArmor = &inventory.Armor[0]
	}
	combatStats := buildStats(base.Level, abilityScores, class, firstArmor)

	return &models.Character{
		Name:          base.Name,
		Level:         base.Level,
		Race:          race,
		Class:         class,
		Stats:         combatStats,
		AbilityScores: abilityScores,
		Skills:        skillList,
		Proficiencies: base.Proficiencies,
		Inventory:     inventory,
		Spells:        spellbook,
	}, nil
}
