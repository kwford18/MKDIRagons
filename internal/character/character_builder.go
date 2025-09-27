package character

import (
	"github.com/kwford18/MKDIRagons/internal/abilities"
	"github.com/kwford18/MKDIRagons/internal/class"
	"github.com/kwford18/MKDIRagons/internal/inventory"
	"github.com/kwford18/MKDIRagons/internal/race"
	"github.com/kwford18/MKDIRagons/internal/skills"
	"github.com/kwford18/MKDIRagons/internal/spells"
	"github.com/kwford18/MKDIRagons/internal/stats"

	"github.com/kwford18/MKDIRagons/templates"
)

func BuildCharacter(base *templates.TemplateCharacter) (*Character, error) {
	var playerRace race.Race
	var playerClass class.Class
	var playerInventory inventory.Inventory

	spellbook := spells.InitSpellbook(base)

	// Concurrent fetch
	if err := race.FetchRace(base, &playerRace); err != nil {
		return nil, err
	}
	if err := class.FetchClass(base, &playerClass); err != nil {
		return nil, err
	}
	if err := inventory.FetchInventory(base, &playerInventory); err != nil {
		return nil, err
	}
	if err := spells.FetchSpells(base, spellbook); err != nil {
		return nil, err
	}

	// Build ability scores & skills
	abilityScores := abilities.BuildAbilityScores(base, playerRace)
	skillList := skills.BuildSkillList(base)

	// Build Combat Stats
	var firstArmor *inventory.Armor
	if len(playerInventory.Armor) > 0 {
		firstArmor = &playerInventory.Armor[0]
	}
	combatStats := stats.BuildStats(base.Level, abilityScores, playerClass, firstArmor)

	return &Character{
		Name:          base.Name,
		Level:         base.Level,
		Race:          playerRace,
		Class:         playerClass,
		Stats:         combatStats,
		AbilityScores: abilityScores,
		Skills:        skillList,
		Proficiencies: base.Proficiencies,
		Inventory:     playerInventory,
		Spells:        spellbook,
	}, nil
}
