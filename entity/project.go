package entity

type Project struct {
	Identifier string `json:"identifier"`
	Module     string `json:"module"`
	Render     Render `json:"render"`

	Database DB       `json:"database"`
	Entities []Entity `json:"entities"`
	Auth     []Auth   `json:"auth"`
	API      API      `json:"api"`
	Events   Events   `json:"events"`

	DisableSelectCombinations bool `json:"select_combinations"`
	AWS                       AWS  `json:"aws"`
}
