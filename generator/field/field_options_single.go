package field

import (
	"fmt"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func OptionsSingleFieldTemplate(f entity.Field, e entity.Entity, dependantEntity *Template) Template {
	template := BaseFieldTemplate(f, e)

	//base
	pl := pluralize.NewClient()
	name := f.Identifier
	template.Type = pl.Singular(helpers.ToCamelCase(name))
	template.InternalType = entity.OptionsSingleFieldType
	template.GenFieldType = "SingleEnumFieldType"
	template.GenRandomValue = fmt.Sprintf("randomvalues.GetRandomOptionValue[%s](%d)", template.Type, len(f.OptionValues))
	template.Custom = true
	template.Enum = true
	template.RepoToMapper = ".ToInt64()"
	template.RepoFromMapper = fmt.Sprintf("entity.%s(model.%s)", helpers.ToCamelCase(f.Identifier), helpers.ToCamelCase(f.Identifier))

	//graph
	template.GraphInType = fmt.Sprintf("String%s", template.GraphRequired)
	template.GraphInTypeOptional = "String"
	template.GraphOutType = fmt.Sprintf("String%s", template.GraphRequired)
	template.GraphGenType = "string"
	template.GraphGenToMapper = fmt.Sprintf("i.%s.StringPtr()", helpers.ToCamelCase(f.Identifier))
	if f.Required {
		template.GraphGenToMapper = fmt.Sprintf("i.%s.String()", helpers.ToCamelCase(f.Identifier))
	}
	template.GraphGenFromMapperParam = fmt.Sprintf("%sentity.%sFromString(%s)", e.Identifier, helpers.ToCamelCase(f.Identifier), f.Identifier)
	template.GraphGenFromMapper = fmt.Sprintf("%s.%sFromString(i.%s)", e.Identifier, pl.Singular(helpers.ToCamelCase(name)), helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapperOptional = fmt.Sprintf("%s.%sFromPointerString(i.%s)", e.Identifier, pl.Singular(helpers.ToCamelCase(name)), helpers.ToCamelCase(f.Identifier))
	if dependantEntity != nil {
		template.GraphGenFromMapper = fmt.Sprintf("%s.%sFromString(i.%s)", pl.Singular(dependantEntity.JSONIdentifier), pl.Singular(helpers.ToCamelCase(name)), helpers.ToCamelCase(f.Identifier))
		template.GraphGenFromMapperOptional = fmt.Sprintf("%s.%sFromPointerString(i.%s)", pl.Singular(dependantEntity.JSONIdentifier), pl.Singular(helpers.ToCamelCase(name)), helpers.ToCamelCase(f.Identifier))
	}

	//proto
	protoType := helpers.ToCamelCase(fmt.Sprintf("%s_%s", pl.Singular(e.Identifier), pl.Singular(f.Identifier)))
	template.ProtoType = protoType
	template.ProtoEnumOptions = helpers.ProtoEnumOptions(protoType, f.OptionValues)
	template.ProtoToMapper = fmt.Sprintf("pb.%s(e.%s)", protoType, pl.Singular(helpers.ToCamelCase(f.Identifier)))
	template.ProtoFromMapper = fmt.Sprintf("entity.%s(m.Get%s())", pl.Singular(helpers.ToCamelCase(name)), strcase.ToCamel(f.Identifier))

	return template
}
