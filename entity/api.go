package entity

type API struct {
	Domain   string      `json:"domain"`
	Protocol APIProtocol `json:"protocol"`
}

type APIProtocol string

const (
	APIProtocolInvalid  = "invalid"
	APIProtocolAll      = "all"
	APIProtocolGraphQL  = "graphql"
	APIProtocolProtobuf = "protobuf"
)
