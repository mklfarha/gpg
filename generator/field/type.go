package field

import "github.com/maykel/gpg/entity"

type Template struct {
	// original identifier the one the user inputed
	// this will be the column name
	Identifier string

	// this will be used in the code, to keep entity names consistent
	SingularIdentifier string

	// camel case version of singular identifier
	Name string

	// golang type
	Type string

	// parent entity identifier
	EntityIdentifier string

	// type of the field in this code
	InternalType entity.FieldType

	// type of the field in the generated code
	GenFieldType string

	// returns the functions that generates a random value
	GenRandomValue string

	// is primary key
	IsPrimary bool
	Required  bool

	// json tags
	Tags string

	// specific imports
	Import *string

	// json specific config
	JSON     bool
	JSONMany bool
	JSONRaw  bool

	// generated functions
	Custom                bool
	Generated             bool
	GeneratedInsertCustom bool
	GeneratedUpdateCustom bool
	GeneratedFuncInsert   string
	GeneratedFuncUpdate   string

	// enums
	Enum     bool
	EnumMany bool

	// repo mappers
	RepoToMapper   string
	RepoFromMapper string

	// graph mappers
	GraphRequired              string
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

	// proto mappers
	ProtoType        string   // the type in the proto file
	ProtoName        string   // the field name in the proto file
	ProtoEnumOptions []string // enum options
	ProtoToMapper    string   // used in mapper to map from entity to proto
	ProtoFromMapper  string   // user in mapper tp map from proto to entity
	ProtoGenName     string   // field name in generated code by protoc
}
