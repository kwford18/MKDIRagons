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
	"sync"

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

func buildCharacter(base *templates.TemplateCharacter) (*templates.Character, error) {
	var race fetch.Race
	var class fetch.Class
	var inventory fetch.Inventory

	// Spellbook 2D slice to hold character's spells
	spellbook := make([][]fetch.Spell, len(base.Spells.Level))
	for level := range base.Spells.Level {
		spellbook[level] = make([]fetch.Spell, len(base.Spells.Level[level]))
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	errs := make(chan error, 1)

	// Helper function to wrap fetchJSON in a goroutine
	run := func(target fetch.Fetchable, input string) {
		defer wg.Done()
		if err := fetchJSON(target, input); err != nil {
			// Only first error matters here
			select {
			case errs <- err:
			default:
			}
		}
	}

	// Race, Class
	wg.Add(1)
	go run(&race, base.Race)
	wg.Add(1)
	go run(&class, base.Class)

	// Inventory
	// Loop over inventory components and spin up goroutine for each
	for _, armorName := range base.Inventory.Armor {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			var armor fetch.Equipment
			if err := fetchJSON(&armor, name); err != nil {
				select {
				case errs <- err:
				default:
				}
				return
			}

			// Append with lock to prevent data race
			mu.Lock()
			inventory.Armor = append(inventory.Armor, armor)
			mu.Unlock()
		}(armorName)
	}
	for _, weaponName := range base.Inventory.Weapons {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			var weapon fetch.Equipment
			if err := fetchJSON(&weapon, name); err != nil {
				select {
				case errs <- err:
				default:
				}
				return
			}

			// Append with lock to prevent data race
			mu.Lock()
			inventory.Weapons = append(inventory.Weapons, weapon)
			mu.Unlock()
		}(weaponName)
	}
	for _, itemName := range base.Inventory.Items {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			var item fetch.Equipment
			if err := fetchJSON(&item, name); err != nil {
				select {
				case errs <- err:
				default:
				}
				return
			}

			// Append with lock to prevent data race
			mu.Lock()
			inventory.Items = append(inventory.Armor, item)
			mu.Unlock()
		}(itemName)
	}

	// Spells
	// Loop over 2D array and spin up goroutine for each spell at index i, j
	for i := range spellbook {
		for j := range spellbook[i] {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				if err := fetchJSON(&spellbook[i][j], base.Spells.Level[i][j]); err != nil {
					select {
					case errs <- err:
					default:
					}
				}
			}(i, j)
		}
	}

	// Wait for goroutines to finish
	wg.Wait()
	close(errs)

	// Return first err if any
	if err, ok := <-errs; ok {
		return nil, err
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
		log.Fatalf("Incorrect file arguments: %v\n", err)
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
