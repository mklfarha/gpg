package entity

type JSONConfig struct {
	Type       JSONConfigType `json:"type,omitempty"`
	Identifier string         `json:"identifier"`
	Reuse      bool           `json:"reuse"`
	Fields     []Field        `json:"fields,omitempty"`
}
