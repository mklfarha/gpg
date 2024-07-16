package entity

type API struct {
	Domain   string      `json:"domain"`
	Protocol APIProtocol `json:"protocol"`
	HTTPPort string      `json:"httpport"`
	GRPCPort string      `json:"grpcport"`
}

type APIProtocol string

const (
	APIProtocolInvalid  = "invalid"
	APIProtocolAll      = "all"
	APIProtocolGraphQL  = "graphql"
	APIProtocolProtobuf = "protobuf"
)
