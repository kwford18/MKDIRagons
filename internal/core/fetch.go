package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kwford18/MKDIRagons/internal/reference"
)

const DefaultBaseURL = "https://www.dnd5eapi.co/api/2014/"

// FetchJSON fetches and unmarshals JSON from the D&D API using default settings.
// This is the main function that production code should use.
func FetchJSON(property reference.Fetchable, input string) error {
	return FetchJSONWithClient(http.DefaultClient, DefaultBaseURL, property, input)
}

// fetchJSONWithClient is an internal function that allows dependency injection
// for testing purposes. It accepts a custom HTTP client and base URL.
func FetchJSONWithClient(client *http.Client, baseURL string, property reference.Fetchable, input string) error {
	// Format the input
	endpoint := baseURL + property.GetEndpoint()
	noSpaces := strings.ReplaceAll(input, " ", "-")
	lowercase := strings.ToLower(noSpaces)
	formattedURL := endpoint + strings.ReplaceAll(lowercase, "'", "")

	// Make the HTTP request
	resp, err := client.Get(formattedURL)
	if err != nil {
		return fmt.Errorf("FetchJSON: failed to make request to %s: %w", formattedURL, err)
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("FetchJSON: request to %s returned status %d (%s): %s",
			formattedURL, resp.StatusCode, resp.Status, string(body))
	}

	// Decode JSON response
	if err := json.NewDecoder(resp.Body).Decode(property); err != nil {
		return fmt.Errorf("FetchJSON: failed to decode JSON from %s: %w", formattedURL, err)
	}

	return nil
}
