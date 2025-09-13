package builder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kwford18/MKDIRagons/fetch"
	"github.com/kwford18/MKDIRagons/templates"
)

func fetchJSON(property fetch.Fetchable, input string) error {
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

func BuildCharacter(base *templates.TemplateCharacter) (*templates.Character, error) {
	var race fetch.Race
	var class fetch.Class
	var inventory fetch.Inventory

	spellbook := initSpellbook(base)

	// concurrent fetch
	if err := fetchRaceAndClass(base, &race, &class); err != nil {
		return nil, err
	}
	if err := fetchInventory(base, &inventory); err != nil {
		return nil, err
	}
	if err := fetchSpells(base, spellbook); err != nil {
		return nil, err
	}

	return &templates.Character{
		Name:          base.Name,
		Race:          race,
		Class:         class,
		Proficiencies: base.Proficiencies,
		Inventory:     inventory,
		Spells:        spellbook,
	}, nil
}
