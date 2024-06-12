package entity

type Auth struct {
	Enabled       bool   `json:"enabled"`
	BasicUsername string `json:"basic_username"`
	BasicPassword string `json:"basic_password"`
	JWTKey        string `json:"jwt_key"`
}
