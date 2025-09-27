package reference

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
