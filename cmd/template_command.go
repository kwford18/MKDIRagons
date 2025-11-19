package cmd

import (
	"fmt"

	"github.com/kwford18/MKDIRagons/template"
	"github.com/spf13/cobra"
)

var emptyCmd = &cobra.Command{
	Use:   "empty",
	Short: "Generate an empty TOML template file",
	Long:  "Generates an empty TOML file in the template directory. This is a template to scaffold D&D characters with.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Generate empty TOML if needed
		if err := template.GenerateEmptyTOML(); err != nil {
			return fmt.Errorf("error generating empty TOML: %w", err)
		}

		fmt.Println("Generated empty TOML file in toml-characters/")
		return nil
	},
}

func init() {
	// Add the build command to the root
	rootCmd.AddCommand(emptyCmd)
}
