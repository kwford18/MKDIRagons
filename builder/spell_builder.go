package builder

import (
	"sync"

	"github.com/kwford18/MKDIRagons/models"
	"github.com/kwford18/MKDIRagons/templates"
)

// Initialize the spellbook 2D array to hold character spells
func initSpellbook(base *templates.TemplateCharacter) [][]models.Spell {
	spellbook := make([][]models.Spell, len(base.Spells.Level))
	for level := range base.Spells.Level {
		spellbook[level] = make([]models.Spell, len(base.Spells.Level[level]))
	}
	return spellbook
}

// Handle concurrent fetching for spells into the 2D spellbook
func fetchSpells(base *templates.TemplateCharacter, spellbook [][]models.Spell) error {
	var wg sync.WaitGroup
	errs := make(chan error, len(spellbook))

	for i := range spellbook {
		for j := range spellbook[i] {
			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				if err := fetchJSON(&spellbook[i][j], base.Spells.Level[i][j]); err != nil {
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
