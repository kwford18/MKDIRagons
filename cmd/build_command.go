package cmd

import (
	"fmt"
	"strings"

	"github.com/kwford18/MKDIRagons/internal/character"
	"github.com/kwford18/MKDIRagons/internal/io"
	"github.com/kwford18/MKDIRagons/templates"
	"github.com/spf13/cobra"
)

var (
	buildFile string
	printChar bool
	rollHP    bool
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds a character from a TOML template file",
	Long:  "Builds a 5e character by parsing the provided TOML file. The character is saved as a JSON file.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Require the --file flag
		if buildFile == "" {
			return fmt.Errorf("please provide a TOML file with --file")
		}
		if !strings.Contains(buildFile, "/") {
			buildFile = "toml-characters/" + buildFile
		}

		base, err := templates.TomlParse(buildFile)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		char, err := character.BuildCharacter(&base, rollHP)
		if err != nil {
			return fmt.Errorf("error building character: %w", err)
		}

		if printChar {
			char.Print()
		}

		if err := io.SaveJSON(char); err != nil {
			return fmt.Errorf("failed to save character as JSON: %w", err)
		}

		return nil
	},
}

func init() {
	// Add the build command to the root
	rootCmd.AddCommand(buildCmd)

	// --file -f flag for providing filepath to TOML
	buildCmd.Flags().StringVarP(&buildFile, "file", "f", "toml-characters/", "Path to the TOML file")

	// --print -p flag for printing character info when parsing TOML
	buildCmd.Flags().BoolVarP(&printChar, "print", "p", false, "Print character info when building")

	// --rollHP -r flag for whether or not a character should roll for HP or use the average of hit die
	buildCmd.Flags().BoolVarP(&rollHP, "rollHP", "r", false, "Roll for character's HP instead of using hit die average")

}
