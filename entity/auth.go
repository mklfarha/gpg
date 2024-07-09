package entity

type Auth struct {
	Type   AuthType   `json:"type"`
	Config AuthConfig `json:"config"`
}

type AuthType string

const (
	BASIC_AUTH_TYPE      = "basic"
	JWT_SERVER_AUTH_TYPE = "jwt"
	KEYCLOAK_AUTH_TYPE   = "keycloak"
)

type AuthConfig struct {
	Basic    *BasicAuthConfig
	JWT      *JWTConfig
	Keycloak *KeycloakConfig
}

type BasicAuthConfig struct {
	BasicUsername string `json:"basic_username"`
	BasicPassword string `json:"basic_password"`
}

type JWTConfig struct {
	JWTKey string `json:"jwt_key"`
}

type KeycloakConfig struct {
	Realm        string `json:"realm"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
