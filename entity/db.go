package entity

type DB struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Driver   string `json:"driver"`
	Version  string `json:"version"`
}
