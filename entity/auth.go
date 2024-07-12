package entity

type Auth struct {
	Enabled bool       `json:"enabled"`
	Type    AuthType   `json:"type"`
	Config  AuthConfig `json:"config"`
}

type AuthType string

const (
	BASIC_AUTH_TYPE      = "basic"
	JWT_SERVER_AUTH_TYPE = "jwt"
	KEYCLOAK_AUTH_TYPE   = "keycloak"
)

type AuthConfig struct {
	Basic    *BasicAuthConfig `json:"basic"`
	JWT      *JWTConfig       `json:"jwt"`
	Keycloak *KeycloakConfig  `json:"keycloak"`
}

type BasicAuthConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTConfig struct {
	Key string `json:"key"`
}

type KeycloakConfig struct {
	Hostname     string `json:"hostname" yaml:"hostname"`
	Realm        string `json:"realm" yaml:"realm"`
	ClientID     string `json:"client_id" yaml:"client_id"`
	ClientSecret string `json:"client_secret" yaml:"client_secret"`
}
