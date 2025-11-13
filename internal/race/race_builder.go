package race

import (
	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/template"
)

// FetchRaceWithFetcher fetches this race's data from the API
func FetchRaceWithFetcher(fetcher core.Fetcher, base *template.Character, race *Race) error {
	return fetcher.FetchJSON(race, base.Race)
}

// FetchRace allows using a custom fetcher for testing
func FetchRace(base *template.Character, race *Race) error {
	return core.FetchJSON(race, base.Race)
}
