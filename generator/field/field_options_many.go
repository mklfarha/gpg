package field

import (
	"fmt"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func OptionsManyFieldTemplate(f entity.Field, e entity.Entity, dependantEntity *Template) Template {

	template := BaseFieldTemplate(f, e)

	//base
	pl := pluralize.NewClient()
	name := f.Identifier
	if dependantEntity != nil && dependantEntity.EntityIdentifier != "" {
		name = fmt.Sprintf("%s_%s", pl.Singular(dependantEntity.Identifier), f.Identifier)
	}
	template.Type = pl.Singular(helpers.ToCamelCase(name))
	template.InternalType = entity.OptionsManyFieldType
	template.GenFieldType = "MultiEnumFieldType"
	template.GenRandomValue = fmt.Sprintf("randomvalues.GetRandomOptionsValues[%s](%d)", template.Type, len(f.OptionValues))
	template.Custom = true
	template.Enum = true
	template.EnumMany = true

	// graph
	template.GraphInType = fmt.Sprintf("[String%s]%s", template.GraphRequired, template.GraphRequired)
	template.GraphInTypeOptional = "[String]"
	template.GraphOutType = fmt.Sprintf("[String%s]%s", template.GraphRequired, template.GraphRequired)
	template.GraphGenType = "[]string"
	template.GraphGenToMapper = fmt.Sprintf("Map%sToModel(i.%s)", pl.Singular(helpers.ToCamelCase(name)), helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapper = fmt.Sprintf("Map%sFromModel(i.%s)", pl.Singular(helpers.ToCamelCase(name)), helpers.ToCamelCase(f.Identifier))

	// proto
	protoType := helpers.ToCamelCase(fmt.Sprintf("%s_%s", pl.Singular(e.Identifier), pl.Singular(f.Identifier)))
	if dependantEntity != nil && dependantEntity.EntityIdentifier != "" {
		protoType = helpers.ToCamelCase(fmt.Sprintf("%s_%s", dependantEntity.EntityIdentifier, protoType))
	}
	template.ProtoType = protoType
	template.ProtoEnumOptions = helpers.ProtoEnumOptions(protoType, f.OptionValues)
	template.ProtoToMapper = fmt.Sprintf("%sSliceToProto(e.%s)", protoType, helpers.ToCamelCase(f.Identifier))
	template.ProtoFromMapper = fmt.Sprintf("%sSliceFromProto(m.Get%s())", protoType, strcase.ToCamel(f.Identifier))

	return template
}
