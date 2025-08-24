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
	path := "templates/characters/"

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

func fetchJSON[T any](url string, target *T, wg *sync.WaitGroup, errCh chan<- error) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		errCh <- fmt.Errorf("failed request to %s: %w", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errCh <- fmt.Errorf("non-200 status code from %s: %d", url, resp.StatusCode)
		return
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		errCh <- fmt.Errorf("failed to decode %s: %w", url, err)
	}
}

func buildCharacter(base templates.TemplateCharacter, emptyRace *fetch.Race, emptyClass *fetch.Class, emptyEquipment *fetch.Equipment, emptySpells *fetch.Spell) (*templates.Character, error) {
	var wg sync.WaitGroup
	errs := make(chan error, 4)

	baseURL := "https://www.dnd5eapi.co/api/2014/"

	// Closure to build the endpoint for a fetchable property
	endpoints := func(property fetch.Fetchable) string {
		return baseURL + property.GetEndpoint()
	}

	lower := func(base string) string {
		return strings.ToLower(base)
	}

	// Goroutine to fetch the data for a Character property
	wg.Add(4)
	go fetchJSON(endpoints(emptyRace)+lower(base.Race), emptyRace, &wg, errs)
	go fetchJSON(endpoints(emptyClass)+lower(base.Class), emptyClass, &wg, errs)
	go fetchJSON(endpoints(emptyEquipment)+lower(base.Equipment), emptyEquipment, &wg, errs)
	go fetchJSON(endpoints(emptySpells)+lower(base.Spells), emptySpells, &wg, errs)
	wg.Wait()

	// Handle errs if any exist
	close(errs)
	var collected []error
	for err := range errs {
		collected = append(collected, err)
	}
	if len(collected) > 0 {
		return nil, fmt.Errorf("failed to build character: %v", collected)
	}

	// Return built Character struct reference
	return &templates.Character{
		Name:      base.Name,
		Race:      *emptyRace,
		Class:     *emptyClass,
		Equipment: *emptyEquipment,
		Spells:    *emptySpells,
	}, nil
}

func main() {
	fileName, err := fileArgs()
	if err != nil {
		panic(err)
	}

	t := tomlParse(fileName)

	var race fetch.Race
	var class fetch.Class
	var equipment fetch.Equipment
	var spells fetch.Spell

	character, err := buildCharacter(t, &race, &class, &equipment, &spells)
	if err != nil {
		log.Fatalf("Error building character: %v\n", err)
	}

	character.Print()
}
