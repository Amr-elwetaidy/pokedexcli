package pokeapi

type Config struct {
	Next     *string
	Previous *string
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreas struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

func (l *LocationArea) GetName() string {
	return l.Name
}
