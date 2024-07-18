package entity

type JSONConfig struct {
	Type       JSONConfigType `json:"type,omitempty"`
	Identifier string         `json:"identifer,omitempty"`
	Reuse      bool           `json:"reuse,omitempty"`
	Fields     []Field        `json:"fields,omitempty"`
}
