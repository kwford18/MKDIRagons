package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/kwford18/MKDIRagons/fetch"
	"github.com/kwford18/MKDIRagons/templates"
)

// Array of the different parts of a character
// var charParts = [...]string{"Race", "Class", "Equipment", "Spells"}

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
func tomlParse(fileName string) templates.Template {
	var t templates.Template
	_, err := toml.DecodeFile(fileName, &t)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Template: %v\n", t)

	return t
}

func buildCharacter(name string, race *fetch.Race, class *fetch.Class, equipment *fetch.Equipment, spells *fetch.Spell) *templates.Character {
	return &templates.Character{
		Name:      name,
		Race:      *race,
		Class:     *class,
		Equipment: *equipment,
		Spells:    *spells,
	}
}

func main() {
	fileName, err := fileArgs()
	if err != nil {
		panic(err)
	}

	t := tomlParse(fileName)

	// 5e API url
	baseURL := "https://www.dnd5eapi.co/api/2014/"
	classURL := strings.ToLower(baseURL + "classes/" + t.Class)
	fmt.Println(classURL)

	// GET from 5e API
	resp, err := http.Get(classURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Unmarshal into struct
	var c fetch.Class
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		panic(err)
	}
}
