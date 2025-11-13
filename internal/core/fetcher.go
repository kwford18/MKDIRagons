package core

import (
	"github.com/kwford18/MKDIRagons/internal/reference"
	"net/http"
)

// Fetcher is an interface for fetching JSON data
type Fetcher interface {
	FetchJSON(property reference.Fetchable, input string) error
}

// HTTPFetcher implements Fetcher using HTTP
type HTTPFetcher struct {
	Client  *http.Client
	BaseURL string
}

// NewFetcher creates a new HTTPFetcher with default settings
func NewFetcher() *HTTPFetcher {
	return &HTTPFetcher{
		Client:  http.DefaultClient,
		BaseURL: DefaultBaseURL,
	}
}

// FetchJSON implements the Fetcher interface
func (f *HTTPFetcher) FetchJSON(property reference.Fetchable, input string) error {
	return FetchJSONWithClient(f.Client, f.BaseURL, property, input)
}

// FOR INTEGRATION WITH EXISTING CODE

// DefaultFetcher is the package-level fetcher for production use
var DefaultFetcher Fetcher = NewFetcher()

// FetchJSON is the convenience function that uses DefaultFetcher
func FetchJSON(property reference.Fetchable, input string) error {
	return DefaultFetcher.FetchJSON(property, input)
}
