package spells

import (
	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/template"
	"sync"
)

// InitSpellbook initializes 2D spellbook array to hold character spells
func InitSpellbook(base *template.Character) [][]Spell {
	spellbook := make([][]Spell, len(base.Spells.Level))
	for level := range base.Spells.Level {
		spellbook[level] = make([]Spell, len(base.Spells.Level[level]))
	}
	return spellbook
}

// FetchSpellsWithFetcher concurrently fetches spells using a custom fetcher
func FetchSpellsWithFetcher(fetcher core.Fetcher, base *template.Character, spellbook [][]Spell) error {
	if fetcher == nil {
		panic("FetchSpellsWithFetcher: fetcher cannot be nil")
	}
	if base == nil {
		panic("FetchSpellsWithFetcher: base TemplateCharacter cannot be nil")
	}
	if spellbook == nil {
		panic("FetchSpellsWithFetcher: spellbook cannot be nil")
	}

	var wg sync.WaitGroup

	// Calculate total number of spells for buffer size
	totalSpells := 0
	for i := range spellbook {
		totalSpells += len(spellbook[i])
	}
	errs := make(chan error, totalSpells)

	for i := range spellbook {
		for j := range spellbook[i] {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				if err := fetcher.FetchJSON(&spellbook[i][j], base.Spells.Level[i][j]); err != nil {
					errs <- err
				}
			}(i, j)
		}
	}

	wg.Wait()
	close(errs)

	if err, ok := <-errs; ok {
		return err
	}
	return nil
}

// FetchSpells concurrently fetches spells using the default fetcher
func FetchSpells(base *template.Character, spellbook [][]Spell) error {
	return FetchSpellsWithFetcher(core.DefaultFetcher, base, spellbook)
}
