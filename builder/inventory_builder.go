package builder

import (
	"sync"

	"github.com/kwford18/MKDIRagons/models"
	"github.com/kwford18/MKDIRagons/templates"
)

// Handle concurrent fetching for Inventory components (Armor, Weapons, Items)
func fetchInventory(base *templates.TemplateCharacter, inv *models.Inventory) error {
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Buffered channel to collect goroutine errors
	errs := make(chan error, len(base.Inventory.Armor)+len(base.Inventory.Weapons)+len(base.Inventory.Items))

	// Loop through inventory and fetch respective JSON
	for _, armorName := range base.Inventory.Armor {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			var armor models.Armor
			if err := fetchJSON(&armor, name); err != nil {
				errs <- err
				return
			}
			mu.Lock()
			inv.Armor = append(inv.Armor, armor)
			mu.Unlock()
		}(armorName)
	}

	for _, weaponName := range base.Inventory.Weapons {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			var weapon models.Weapon
			if err := fetchJSON(&weapon, name); err != nil {
				errs <- err
				return
			}
			mu.Lock()
			inv.Weapons = append(inv.Weapons, weapon)
			mu.Unlock()
		}(weaponName)
	}

	for _, itemName := range base.Inventory.Items {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			var item models.Item
			if err := fetchJSON(&item, name); err != nil {
				errs <- err
				return
			}
			mu.Lock()
			inv.Items = append(inv.Items, item)
			mu.Unlock()
		}(itemName)
	}

	// Wait for all goroutines to finish and close channel
	wg.Wait()
	close(errs)

	// Return first error if any
	if err, ok := <-errs; ok {
		return err
	}

	return nil
}
