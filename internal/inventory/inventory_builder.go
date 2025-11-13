package inventory

import (
	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/templates"
	"sync"
)

// FetchInventoryWithFetcher allows using a custom fetcher for testing
func FetchInventoryWithFetcher(fetcher core.Fetcher, base *templates.TemplateCharacter, inv *Inventory) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	errs := make(chan error, len(base.Inventory.Armor)+len(base.Inventory.Weapons)+len(base.Inventory.Items))

	// Fetch all armor in parallel
	for _, armorName := range base.Inventory.Armor {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			var armor Armor
			if err := fetcher.FetchJSON(&armor, name); err != nil {
				errs <- err
				return
			}
			mu.Lock()
			inv.Armor = append(inv.Armor, armor)
			mu.Unlock()
		}(armorName)
	}

	// Fetch all weapons in parallel
	for _, weaponName := range base.Inventory.Weapons {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			var weapon Weapon
			if err := fetcher.FetchJSON(&weapon, name); err != nil {
				errs <- err
				return
			}
			mu.Lock()
			inv.Weapons = append(inv.Weapons, weapon)
			mu.Unlock()
		}(weaponName)
	}

	// Fetch all items in parallel
	for _, itemName := range base.Inventory.Items {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			var item Item
			if err := fetcher.FetchJSON(&item, name); err != nil {
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

// FetchInventory uses the default fetcher for production
func FetchInventory(base *templates.TemplateCharacter, inv *Inventory) error {
	return FetchInventoryWithFetcher(core.DefaultFetcher, base, inv)
}
