package inventory

import (
	"github.com/kwford18/MKDIRagons/internal/core"
	"sync"

	"github.com/kwford18/MKDIRagons/templates"
)

// FetchInventory concurrently fetches for Inventory components (Armor, Weapons, Items)
func FetchInventory(base *templates.TemplateCharacter, inv *Inventory) error {
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Buffered channel to collect goroutine errors
	errs := make(chan error, len(base.Inventory.Armor)+len(base.Inventory.Weapons)+len(base.Inventory.Items))

	// Loop through inventory and fetch respective JSON
	for _, armorName := range base.Inventory.Armor {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			var armor Armor
			if err := core.FetchJSON(&armor, name); err != nil {
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
			var weapon Weapon
			if err := core.FetchJSON(&weapon, name); err != nil {
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
			var item Item
			if err := core.FetchJSON(&item, name); err != nil {
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
