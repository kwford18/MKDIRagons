package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/kwford18/MKDIRagons/fetch"
	"github.com/kwford18/MKDIRagons/templates"
)

// Builds the file path
func fileArgs() (string, error) {
	if len(os.Args) > 2 {
		return "Too many args", errors.New("too many args")
	} else if len(os.Args) < 2 {
		return "Too few args", errors.New("too few args")
	}

	// Load file
	path := "templates/"

	// Check if file extension was provided
	if filepath.Ext(os.Args[1]) != ".toml" && filepath.Ext(os.Args[1]) != "" {
		return "Incorrect file type", errors.New("incorrect file type")
	}

	file := path + os.Args[1]
	return file, nil
}

// Parses the provided TOML file into the Template struct
func tomlParse(fileName string) templates.TemplateCharacter {
	var t templates.TemplateCharacter
	_, err := toml.DecodeFile(fileName, &t)
	if err != nil {
		panic(err)
	}

	return t
}

func saveJSON(character *templates.Character) {
	// Ensure ./characters exists
	dir := "characters"
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	fileName := character.Name + ".json"

	// Build file path inside ./characters
	filePath := filepath.Join(dir, fileName)

	// Convert to pretty JSON
	data, err := json.MarshalIndent(character, "", "  ")
	if err != nil {
		panic(err)
	}

	// Open file with truncate + write-only + create
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write JSON
	if _, err := file.Write(data); err != nil {
		panic(err)
	}
}

func fetchJSON(property fetch.Fetchable, input string) error {
	baseURL := "https://www.dnd5eapi.co/api/2014/"

	// Format
	endpoint := baseURL + property.GetEndpoint()
	no_spaces := strings.ReplaceAll(input, " ", "-")
	lowercase := strings.ToLower(no_spaces)
	formatted_url := endpoint + strings.ReplaceAll(lowercase, "'", "")

	fmt.Printf("Formatted URL: %s\n", formatted_url)

	resp, err := http.Get(formatted_url)
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

func buildCharacter(base *templates.TemplateCharacter) (*templates.Character, error) {
	var race fetch.Race
	var class fetch.Class
	var inventory fetch.Inventory

	// Spellbook 2D slice to hold character's spells
	spellbook := make([][]fetch.Spell, len(base.Spells.Level))
	for level := range base.Spells.Level {
		spellbook[level] = make([]fetch.Spell, len(base.Spells.Level[level]))
	}

	// TODO: Make concurrent

	// Name, Class
	if err := fetchJSON(&race, base.Race); err != nil {
		return nil, err
	}
	if err := fetchJSON(&class, base.Class); err != nil {
		return nil, err
	}

	// Inventory
	for _, armorName := range base.Inventory.Armor {
		var armor fetch.Equipment
		if err := fetchJSON(&armor, armorName); err != nil {
			return nil, err
		}
		inventory.Armor = append(inventory.Armor, armor)
	}
	for _, weaponName := range base.Inventory.Weapons {
		var weapon fetch.Equipment
		if err := fetchJSON(&weapon, weaponName); err != nil {
			return nil, err
		}
		inventory.Weapons = append(inventory.Weapons, weapon)
	}
	for _, itemName := range base.Inventory.Items {
		var item fetch.Equipment
		if err := fetchJSON(&item, itemName); err != nil {
			return nil, err
		}
		inventory.Items = append(inventory.Items, item)
	}

	// Spells
	for i := range spellbook {
		for j := range spellbook[i] {
			if err := fetchJSON(&spellbook[i][j], base.Spells.Level[i][j]); err != nil {
				return nil, err
			}
		}
	}

	return &templates.Character{
		Name:      base.Name,
		Race:      race,
		Class:     class,
		Inventory: inventory,
		Spells:    spellbook,
	}, nil
}

func main() {
	fileName, err := fileArgs()
	if err != nil {
		panic(err)
	}

	base := tomlParse(fileName)

	fmt.Println()
	base.Print()
	fmt.Println()

	character, err := buildCharacter(&base)
	if err != nil {
		log.Fatalf("Error building character: %v\n", err)
	}

	fmt.Println()
	character.Print()

	saveJSON(character)
}
