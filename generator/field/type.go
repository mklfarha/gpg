package field

import "github.com/maykel/gpg/entity"

type Template struct {
	Identifier                 string
	Name                       string
	Type                       string
	EntityIdentifier           string
	InternalType               entity.FieldType
	IsPrimary                  bool
	Required                   bool
	Tags                       string
	Import                     *string
	JSON                       bool
	JSONMany                   bool
	JSONRaw                    bool
	Custom                     bool
	Generated                  bool
	GeneratedFuncInsert        string
	GeneratedFuncUpdate        string
	Enum                       bool
	EnumMany                   bool
	RepoToMapper               string
	RepoFromMapper             string
	GraphName                  string
	GraphModelName             string
	GraphOutType               string
	GraphInType                string
	GraphInTypeOptional        string
	GraphGenType               string
	GraphGenToMapper           string
	GraphGenFromMapperParam    string
	GraphGenFromMapper         string
	GraphGenFromMapperOptional string
	ProtoType                  string   // the type in the proto file
	ProtoName                  string   // the field name in the proto file
	ProtoEnumOptions           []string // enum options
	ProtoToMapper              string   // used in mapper to map from entity to proto
	ProtoFromMapper            string   // user in mapper tp map from proto to entity
	ProtoGenName               string   // field name in generated code by protoc
}
