package models

type Reference struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

// Interface for getting endpoints with different 5e API resources
type Fetchable interface {
	GetEndpoint() string
	Print()
}

// Interface for common methods between models.Character and templates.TemplateCharacter
type CharacterModel interface {
	Print()
	ProficiencyBonus() int
	GetSkillAbility(string) Ability
	Modifier(Ability) int
}
