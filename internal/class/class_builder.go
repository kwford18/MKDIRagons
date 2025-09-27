package class

import (
	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/templates"
	"sync"
)

func FetchClass(base *templates.TemplateCharacter, class *Class) error {
	var wg sync.WaitGroup
	errs := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := core.FetchJSON(class, base.Class); err != nil {
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
