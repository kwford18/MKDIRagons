package builder

import (
	"sync"

	"github.com/kwford18/MKDIRagons/models"
	"github.com/kwford18/MKDIRagons/templates"
)

// Handle concurrent fetching JSON for Race & Class, as it only requires one request
func fetchRaceAndClass(base *templates.TemplateCharacter, race *models.Race, class *models.Class) error {
	var wg sync.WaitGroup
	errs := make(chan error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := fetchJSON(race, base.Race); err != nil {
			errs <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := fetchJSON(class, base.Class); err != nil {
			errs <- err
		}
	}()

	wg.Wait()
	close(errs)

	if err, ok := <-errs; ok {
		return err
	}
	return nil
}
