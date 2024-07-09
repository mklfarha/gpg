package entity

import "fmt"

type Project struct {
	Identifier                string          `json:"identifier"`
	Render                    Render          `json:"render"`
	Entities                  []Entity        `json:"entities"`
	Database                  DB              `json:"database"`
	Auth                      []Auth          `json:"auth"`
	API                       API             `json:"api"`
	AWS                       AWS             `json:"aws"`
	Protocol                  ProjectProtocol `json:"protocol"`
	DisableSelectCombinations bool            `json:"select_combinations"`
}

type API struct {
	URL string `json:"url"`
}

type ProjectProtocol string

const (
	ProjectProtocolInvalid  = "invalid"
	ProjectProtocolAll      = "all"
	ProjectProtocolGraphQL  = "graphql"
	ProjectProtocolProtobuf = "protobuf"
)

func (p Project) HasBasicAuth() bool {
	found, config := p.AuthByType(BASIC_AUTH_TYPE)
	if found && config.Config.Basic != nil {
		return true
	}
	return false
}

func (p Project) BasicAuth() Auth {
	_, config := p.AuthByType(BASIC_AUTH_TYPE)
	return config
}

func (p Project) HasJWTAuth() bool {
	found, config := p.AuthByType(JWT_SERVER_AUTH_TYPE)
	if found && config.Config.JWT != nil {
		return true
	}
	return false
}

func (p Project) JWTAuth() Auth {
	_, config := p.AuthByType(JWT_SERVER_AUTH_TYPE)
	return config
}

func (p Project) HasKeycloakAuth() bool {
	found, config := p.AuthByType(KEYCLOAK_AUTH_TYPE)
	if found && config.Config.JWT != nil {
		return true
	}
	return false
}

func (p Project) KeycloakAuth() Auth {
	_, config := p.AuthByType(KEYCLOAK_AUTH_TYPE)
	return config
}

func (p Project) AuthByType(t AuthType) (bool, Auth) {
	for _, a := range p.Auth {
		if a.Type == t {
			return true, a
		}
	}
	return false, Auth{}
}

func (p Project) AuthImport() string {
	if p.HasJWTAuth() && p.JWTAuth().Config.JWT != nil {
		return fmt.Sprintf("auth \"%s/auth/jwtserver\"", p.Identifier)
	}

	if p.HasKeycloakAuth() && p.KeycloakAuth().Config.Keycloak != nil {
		return fmt.Sprintf("auth \"%s/auth/keycloak\"", p.Identifier)
	}
	return ""
}
