# MKDIRagons

*A CLI tool for building D&D 5e characters from TOML files.*

## Overview

MKDIRagons converts a simple character description written in TOML into
a fully populated Dungeons & Dragons 5e character. It fetches rules data
from the [Unofficial 5e API](https://www.dnd5eapi.co/) and produces clean JSON output.

Future updates will include optional scraping for missing fields using
5e Wikidot.

## Installation

``` bash
go install github.com/kwford18/MKDIRagons@latest
```

## Features

-   Build characters from TOML
-   Load and display JSON character files
-   Generate empty TOML templates
-   Supports random HP rolling
-   Built using Cobra for robust CLI structure
-   Test coverage using Testify and httptest

## Commands

| Command | Description                             |
|---------|-----------------------------------------|
| `build` | Build a JSON character from a TOML file |
| `load`  | Load & display a JSON character file    |
| `empty` | Generate an empty TOML template         |


### Global Flags

-   `--file, -f` --- Provide a file path instead of using the default
    directory

### Build Flags

-   `--print, -p` --- Print the generated character to the console
-   `--rollHP, -r` --- Roll HP instead of using averages

## Default Directories

-   TOML input: `toml-characters/`
-   Empty templates: `toml-characters/`
-   Generated JSON output: `characters/`
---
## Example Usage

### Generate an Empty Template

``` bash
MKDIRagons empty
```

### Build a Character

``` bash
MKDIRagons build -f example_character.toml
```

### Print Character During Build

``` bash
MKDIRagons build -f example_character.toml -p
```

### Roll HP When Building

``` bash
MKDIRagons build -f example_character.toml -r
```

### Load a Character

``` bash
MKDIRagons load -f example_character.json
```

## Example TOML File Structure

<details>
<summary>The empty/template TOML file</summary>
Note that all fields should be filled in by the user.

``` TOML
name = ""
level = 1
race = ""
class = ""
proficiencies = []

[ability_scores]
Strength = 10
Dexterity = 10
Constitution = 10
Wisdom = 10
Intelligence = 10
Charisma = 10

[inventory]
Weapons = []
Armor = []
Items = []

[spells]
Level = [[], [], [], [], [], [], [], [], [], []]
```
</details>

<details>
<summary>The Level 5 Cleric from the Examples directory</summary>

``` TOML
name = "Leki"
level = 5
race = "Tiefling"
subrace = ""
class = "Cleric"
subclass = "Life"
proficiencies = ["Arcana", "Insight"]

[ability_scores]
strength = 10
dexterity = 10
constitution = 14
wisdom = 15
intelligence = 12
charisma = 13

[inventory]
weapons = ["Rapier"]
armor = ["Padded Armor", "Leather Armor"]
items = ["Amulet", "Abacus", "Acid Vial", "alchemists fire flask"]

[spells]
level = [
["Acid Splash"],                 # Cantrip
["Alarm", "Identify", "Shield"], # Level 1
["Continual Flame"],             # Level 2
[],                              # Level 3
[],                              # Level 4
[],                              # Level 5
[],                              # Level 6
[],                              # Level 7
[],                              # Level 8
[],                              # Level 9
]
```
</details>

## Current Limitations
- Limited to the 5e API, which exclusively has the 2014 5e content
- No styling options for viewing a character
- Editing a character requires rebuilding character or directly modifying JSON
- Limited to the terminal

## Roadmap

-   Expanded test coverage
-   Wikidot scraping for data beyond the API
-   TUI support for interactive building
-   Character sheet export
-   Multi-classing support
-   Homebrew plugin system

## License

MIT License