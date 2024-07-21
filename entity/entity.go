package entity

import (
	"slices"
)

type Entity struct {
	Identifier    string        `json:"identifier"`
	Render        Render        `json:"render"`
	Operation     []Operation   `json:"operations"`
	Validations   []Validation  `json:"validations"`
	Fields        []Field       `json:"fields"`
	CustomQueries []CustomQuery `json:"custom_queries"`
}

func (e Entity) UsesRandomValues() bool {
	for _, f := range e.Fields {
		if slices.Contains(UsesRandomValues, f.Type) {
			return true
		}
	}
	return false
}
