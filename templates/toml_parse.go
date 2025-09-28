package templates

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Parses the provided TOML file into the Template struct
func TomlParse(fileName string) (TemplateCharacter, error) {
	var t TemplateCharacter
	_, err := toml.DecodeFile(fileName, &t)
	if err != nil {
		return t, fmt.Errorf("failed to parse file: %w", err)
	}

	return t, nil
}
