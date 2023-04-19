package entity

type JSONConfig struct {
	Type   JSONConfigType `json:"type,omitempty"`
	Fields []Field        `json:"fields,omitempty"`
}
