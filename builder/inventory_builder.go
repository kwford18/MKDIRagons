package builder

import (
	"sync"

	"github.com/kwford18/MKDIRagons/fetch"
	"github.com/kwford18/MKDIRagons/templates"
)

// Handle concurrent fetching for Inventory components (Armor, Weapons, Items)
func fetchInventory(base *templates.TemplateCharacter, inv *fetch.Inventory) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	errs := make(chan error, len(base.Inventory.Armor)+len(base.Inventory.Weapons)+len(base.Inventory.Items))

	for _, armorName := range base.Inventory.Armor {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			var armor fetch.Equipment
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
			var weapon fetch.Equipment
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
			var item fetch.Equipment
			if err := fetchJSON(&item, name); err != nil {
				errs <- err
				return
			}
			mu.Lock()
			inv.Items = append(inv.Items, item)
			mu.Unlock()
		}(itemName)
	}

	wg.Wait()
	close(errs)

	if err, ok := <-errs; ok {
		return err
	}
	return nil
}
