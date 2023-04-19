package entity

type Auth struct {
	BasicUsername string `json:"basic_username"`
	BasicPassword string `json:"basic_password"`
	JWTKey        string `json:"jwt_key"`
}
