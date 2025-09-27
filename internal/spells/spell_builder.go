package spells

import (
	"github.com/kwford18/MKDIRagons/internal/core"
	"sync"

	"github.com/kwford18/MKDIRagons/templates"
)

// InitSpellbook initializes 2D spellbook array to hold character spells
func InitSpellbook(base *templates.TemplateCharacter) [][]Spell {
	spellbook := make([][]Spell, len(base.Spells.Level))
	for level := range base.Spells.Level {
		spellbook[level] = make([]Spell, len(base.Spells.Level[level]))
	}
	return spellbook
}

// FetchSpells concurrently fetches spells into the 2D spellbook
func FetchSpells(base *templates.TemplateCharacter, spellbook [][]Spell) error {
	var wg sync.WaitGroup
	errs := make(chan error, len(spellbook))

	for i := range spellbook {
		for j := range spellbook[i] {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				if err := core.FetchJSON(&spellbook[i][j], base.Spells.Level[i][j]); err != nil {
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
