package class

import (
	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/template"
)

// FetchClassWithFetcher fetches this race's data from the API
func FetchClassWithFetcher(fetcher core.Fetcher, base *template.Character, class *Class) error {
	return fetcher.FetchJSON(class, base.Class)
}

// FetchClass allows using a custom fetcher for testing
func FetchClass(base *template.Character, class *Class) error {
	return core.FetchJSON(class, base.Class)
}
