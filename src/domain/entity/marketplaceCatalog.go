package entity

type MarketplaceCatalog struct {
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Description string        `json:"description"`
	ImageUrl    string        `json:"imageUrl"`
	Services    []string      `json:"services"`
	Mappings    []interface{} `json:"mapppings"`
	Steps       []string      `json:"steps"`
}
