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

Command   Description
  --------- -----------------------------------------
`build`   Build a JSON character from a TOML file
`load`    Load & display a JSON character file
`empty`   Generate an empty TOML template

### Global Flags

-   `--file, -f` --- Provide a file path instead of using the default
    directory

### Build Flags

-   `--print, -p` --- Print the generated character to the console\
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

## Roadmap

-   Expanded test coverage
-   Wikidot scraping
-   TUI support for interactive building
-   Character sheet export
-   Multi-classing support
-   Homebrew plugin system

## License

MIT License
