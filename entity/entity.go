package entity

type Entity struct {
	Identifier    string        `json:"identifier"`
	Render        Render        `json:"render"`
	Operation     []Operation   `json:"operations"`
	Validations   []Validation  `json:"validations"`
	Fields        []Field       `json:"fields"`
	CustomQueries []CustomQuery `json:"custom_queries"`
}
