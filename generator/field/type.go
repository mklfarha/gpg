package field

type Template struct {
	Identifier                 string
	Name                       string
	Type                       string
	IsPrimary                  bool
	Required                   bool
	Tags                       string
	Import                     *string
	JSON                       bool
	JSONMany                   bool
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
	ProtoType                  string
	ProtoName                  string
	ProtoEnumOptions           []string
}
