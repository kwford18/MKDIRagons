package builder

import (
	"github.com/kwford18/MKDIRagons/models"
	"github.com/kwford18/MKDIRagons/templates"
	"sync"
)

func fetchClass(base *templates.TemplateCharacter, class *models.Class) error {
	var wg sync.WaitGroup
	errs := make(chan error, 2)

	wg.Add(1)
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
