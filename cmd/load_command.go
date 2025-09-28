package cmd

import (
	"fmt"
	"strings"

	"github.com/kwford18/MKDIRagons/internal/io"
	"github.com/spf13/cobra"
)

var loadFile string

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load a character from JSON",
	Long:  "Load and print to the console a character from JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the --file flag
		if loadFile == "" {
			return fmt.Errorf("please provide a JSON file with --file")
		}
		if !strings.Contains(loadFile, "/") {
			loadFile = "characters/" + loadFile
		}

		// Load the character from JSON
		char, err := io.LoadCharacter(loadFile)
		if err != nil {
			return fmt.Errorf("failed to load character: %w", err)
		}

		// Print the loaded character
		char.Print()

		return nil
	},
}

func init() {
	// Add the build command to the root
	rootCmd.AddCommand(loadCmd)

	// --load -l flag for providing filepath to
	loadCmd.Flags().StringVarP(&loadFile, "file", "f", "characters/", "Path to the JSON file")

}
