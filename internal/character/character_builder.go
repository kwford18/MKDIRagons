package character

import (
	"sync"

	"github.com/kwford18/MKDIRagons/internal/abilities"
	"github.com/kwford18/MKDIRagons/internal/class"
	"github.com/kwford18/MKDIRagons/internal/core"
	"github.com/kwford18/MKDIRagons/internal/inventory"
	"github.com/kwford18/MKDIRagons/internal/race"
	"github.com/kwford18/MKDIRagons/internal/skills"
	"github.com/kwford18/MKDIRagons/internal/spells"
	"github.com/kwford18/MKDIRagons/internal/stats"
	"github.com/kwford18/MKDIRagons/template"
)

// BuildCharacterWithFetcher builds a character using a custom fetcher (for testing)
func BuildCharacterWithFetcher(fetcher core.Fetcher, base *template.Character, rollHP bool) (*Character, error) {
	var playerRace race.Race
	var playerClass class.Class
	var playerInventory inventory.Inventory
	spellbook := spells.InitSpellbook(base)

	// Concurrent fetch of all character components
	var wg sync.WaitGroup
	errs := make(chan error, 4) // Buffer size = number of parallel fetch operations

	// Fetch race
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := race.FetchRaceWithFetcher(fetcher, base, &playerRace); err != nil {
			errs <- err
		}
	}()

	// Fetch class
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := class.FetchClassWithFetcher(fetcher, base, &playerClass); err != nil {
			errs <- err
		}
	}()

	// Fetch inventory (internally fetches multiple items in parallel)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := inventory.FetchInventoryWithFetcher(fetcher, base, &playerInventory); err != nil {
			errs <- err
		}
	}()

	// Fetch spells (internally fetches multiple spells in parallel)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := spells.FetchSpellsWithFetcher(fetcher, base, spellbook); err != nil {
			errs <- err
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()
	close(errs)

	// Check if any fetch operation failed
	if err, ok := <-errs; ok {
		return nil, err
	}

	// Build ability scores, saves, & skills
	abilityScores := abilities.BuildAbilityScores(base, playerRace)
	savingThrows := abilities.BuildSavingThrows(base, abilityScores, &playerClass)
	skillList := skills.BuildSkillList(base)

	// Build Combat Stats
	var firstArmor *inventory.Armor
	if len(playerInventory.Armor) > 0 {
		firstArmor = &playerInventory.Armor[0]
	}
	combatStats, err := stats.BuildStats(base.Level, abilityScores, playerClass, rollHP, firstArmor)
	if err != nil {
		return nil, err
	}

	return &Character{
		Name:          base.Name,
		Level:         base.Level,
		Race:          playerRace,
		Class:         playerClass,
		Stats:         combatStats,
		AbilityScores: abilityScores,
		SavingThrows:  savingThrows,
		Skills:        skillList,
		Proficiencies: base.Proficiencies,
		Inventory:     playerInventory,
		Spells:        spellbook,
	}, nil
}

// BuildCharacter builds a character using the default fetcher (for production)
func BuildCharacter(base *template.Character, rollHP bool) (*Character, error) {
	return BuildCharacterWithFetcher(core.DefaultFetcher, base, rollHP)
}
