# MKDIRagons
*A fast, flexible CLI tool for building D&D 5e characters from TOML files.*  
  
Build a basic character using the blank TOML file, which is then parsed into a feature-complete 5e character!  
Data is fetched from the [Unofficial 5e API](https://www.dnd5eapi.co/).  
Future updates will support scraping missing data from [5e Wikidot](https://dnd5e.wikidot.com/)
  
### Built Using Cobra
MKDIRagons is built on top of the [Cobra](https://github.com/spf13/cobra) library, which offers a suite of tools for building CLIs.
This gives it the power and extendability needed to support the complexity of a D&D character, without compromising on speed or readability

### Unit Testing using Testify
MKDIRagons has wide (and still increasing) test coverage using the [Testify](https://github.com/stretchr/testify) library.
Testify offers tools that are especially helpful with a project like this. Namely:
- Easy assertions
- Mock structures
- Testing suite interfaces and functions
  
MKDIRagons also utilizes the httptest core library for mock servers.
## Currently supported commands:  
- `build`: Build a JSON character file using the provided TOML file  
- `load`:  Load the provided JSON character and print their information  
- `empty`: Generate an empty TOML template file for making characters in toml-empty/  
  
The build & load commands both support the `--file/-f` flags, which are  
for providing a path to the file to build from/load respectively. 
Build looks for the **toml-characters** directory by default, but this can be overridden with the -f flag.  
Run MKDIRagons empty to generate a template file in this directory.  

The build command also supports the following two switches:  
- `--print/-p`:  Print the content of the built file (equivalent to loading
  file)  
- `--rollHP/-r`: Use random number generation for calculating health
  instead of the hit die average  

# Example Usage

### Generating an Empty Character
``` bash
MKDIRagons empty
```

### Building a Character
Building with the toml-characters/ directory and the long or short flag:
``` bash
MKDIRagons build --file example_character.toml

MKDIRagons build -f example_character
```

The toml-characters/ directory can be overridden as such
``` bash
MKDIRagons build -f path/to/character.toml
```

Printing the character stats during build using -p switch
``` bash
MKDIRagons build -f example_character.toml -p
```

Rolling character HP during build using -r switch
``` bash
MKDIRagons build -f example_character.toml -r
```

### Loading a Character
Load from the generated characters/ directory
``` bash
MKDIRagons load -f example_character.json
```

Load from the provided directory (works identical to building)
``` bash
MKDIRagons load -f path/to/example_character.json
```