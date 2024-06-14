package entity

type Project struct {
	Identifier string          `json:"identifier"`
	Render     Render          `json:"render"`
	Entities   []Entity        `json:"entities"`
	Database   DB              `json:"database"`
	Auth       Auth            `json:"auth"`
	API        API             `json:"api"`
	AWS        AWS             `json:"aws"`
	Protocol   ProjectProtocol `json:"protocol"`
}

type API struct {
	URL string `json:"url"`
}

type ProjectProtocol string

const (
	ProjectProtocolInvalid  = "invalid"
	ProjectProtocolGraphQL  = "graphql"
	ProjectProtocolProtobuf = "protobuf"
)
