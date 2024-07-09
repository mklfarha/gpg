package entity

type Auth struct {
	Type          AuthType `json:"type"`
	Enabled       bool     `json:"enabled"`
	BasicUsername string   `json:"basic_username"`
	BasicPassword string   `json:"basic_password"`
	JWTKey        string   `json:"jwt_key"`
}

type AuthType string

const (
	BASIC_IN_MEMORT_JWT_SERVER = "basic-jwt"
)
