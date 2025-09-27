package core

import (
	"encoding/json"
	"fmt"
	"github.com/kwford18/MKDIRagons/internal/reference"
	"net/http"
	"strings"
)

func FetchJSON(property reference.Fetchable, input string) error {
	baseURL := "https://www.dnd5eapi.co/api/2014/"

	// Format
	endpoint := baseURL + property.GetEndpoint()
	noSpaces := strings.ReplaceAll(input, " ", "-")
	lowercase := strings.ToLower(noSpaces)
	formattedURL := endpoint + strings.ReplaceAll(lowercase, "'", "")

	// fmt.Printf("Formatted URL: %s\n", formatted_url)

	resp, err := http.Get(formattedURL)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(property); err != nil {
		fmt.Printf("Error: %+v\n", err)
		return err
	}

	// fmt.Printf("Fetched: %v\n", property)

	return nil
}
