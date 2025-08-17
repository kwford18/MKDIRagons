package api

type Reference struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

// Interface for getting endpoints with different 5e API resources
type Fetchable interface {
	GetEndpoint() string
}
