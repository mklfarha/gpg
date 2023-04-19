package entity

type Validation struct {
	Operation Operation      `json:"operation"`
	Type      ValidationType `json:"type"`
	Rule      ValidationRule `json:"rule"`
}
