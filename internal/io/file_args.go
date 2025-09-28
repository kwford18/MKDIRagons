package io

import (
	"fmt"
	"os"
	"path/filepath"
)

func FileArgs() (string, error) {
	if len(os.Args) > 2 {
		return "Too many args", fmt.Errorf("too many args: %s", os.Args)
	} else if len(os.Args) < 2 {
		return "Too few args", fmt.Errorf("too few args: %s", os.Args)
	}

	// Load file
	path := "templates/toml-characters/"

	// Check if file extension was provided
	if filepath.Ext(os.Args[1]) != ".toml" && filepath.Ext(os.Args[1]) != "" {
		return "Incorrect file type", fmt.Errorf("incorrect file type: %s", filepath.Ext(os.Args[1]))
	}

	file := path + os.Args[1]
	return file, nil
}
