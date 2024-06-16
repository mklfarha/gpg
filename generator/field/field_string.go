package field

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func StringFieldTemplate(f entity.Field, e entity.Entity) Template {
	graphRequired := ""
	if f.Required {
		graphRequired = "!"
	}
	graphModelName := strings.ReplaceAll(helpers.ToCamelCase(f.Identifier), "Url", "URL")
	graphGenFromMapper := fmt.Sprintf("i.%s", graphModelName)
	graphGenFromMapperOptional := fmt.Sprintf("StringFromPointer(i.%s)", graphModelName)
	graphGenFromMapperParam := f.Identifier
	graphGenType := "string"
	graphGenToMapper := fmt.Sprintf("i.%s", helpers.ToCamelCase(f.Identifier))
	if !f.Required {
		graphGenToMapper = fmt.Sprintf("&i.%s", helpers.ToCamelCase(f.Identifier))
		graphGenFromMapper = graphGenFromMapperOptional
		graphGenFromMapperParam = fmt.Sprintf("mapper.StringFromPointer(%s)", f.Identifier)
		graphGenType = "*string"
	}
	return Template{
		Identifier:                 f.Identifier,
		Name:                       helpers.ToCamelCase(f.Identifier),
		Type:                       "string",
		IsPrimary:                  f.StorageConfig.PrimaryKey,
		Required:                   f.Required,
		Tags:                       helpers.ResolveTags(f),
		Import:                     nil,
		Custom:                     false,
		Generated:                  f.Autogenerated.Type != entity.InvalidAutogeneratedType,
		GeneratedFuncInsert:        resolveGeneratedFuncInsert(e, f),
		GeneratedFuncUpdate:        resolveGeneratedFuncUpdate(e, f),
		RepoToMapper:               "",
		RepoFromMapper:             fmt.Sprintf("model.%s", helpers.ToCamelCase(f.Identifier)),
		GraphName:                  f.Identifier,
		GraphModelName:             graphModelName,
		GraphInType:                fmt.Sprintf("String%s", graphRequired),
		GraphInTypeOptional:        "String",
		GraphOutType:               fmt.Sprintf("String%s", graphRequired),
		GraphGenType:               graphGenType,
		GraphGenToMapper:           graphGenToMapper,
		GraphGenFromMapperParam:    graphGenFromMapperParam,
		GraphGenFromMapper:         graphGenFromMapper,
		GraphGenFromMapperOptional: graphGenFromMapperOptional,
		ProtoType:                  "string",
		ProtoName:                  helpers.ToSnakeCase(f.Identifier),
		ProtoToMapper:              fmt.Sprintf("e.%s", helpers.ToCamelCase(f.Identifier)),
		ProtoFromMapper:            fmt.Sprintf("m.Get%s()", strcase.ToCamel(f.Identifier)),
		ProtoGenName:               strcase.ToCamel(f.Identifier),
	}
}
