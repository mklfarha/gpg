package entity

type Project struct {
	Identifier string   `json:"identifier"`
	Render     Render   `json:"render"`
	Entities   []Entity `json:"entities"`
	Database   DB       `json:"database"`
	Auth       Auth     `json:"auth"`
	API        API      `json:"api"`
	AWS        AWS      `json:"aws"`
}

type API struct {
	URL string `json:"url"`
}
