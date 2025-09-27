package models

type Reference struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

// Fetchable Interface for getting endpoints & printing different 5e API resources
type Fetchable interface {
	GetEndpoint() string
	Print()
}

// CharacterModel Interface for common methods between models.Character and templates.TemplateCharacter
type CharacterModel interface {
	Print()
	ProficiencyBonus() int
	GetSkillAbility(string) Ability
	Modifier(Ability) int
}
