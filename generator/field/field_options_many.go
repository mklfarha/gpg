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
	template.Type = pl.Singular(helpers.ToCamelCase(name))
	template.InternalType = entity.OptionsManyFieldType
	template.GenFieldType = "MultiEnumFieldType"
	template.GenRandomValue = fmt.Sprintf("randomvalues.GetRandomOptionsValues[%s](%d)", template.Type, len(f.OptionValues))
	template.Custom = true
	template.Enum = true
	template.EnumMany = true
	template.RepoToMapper = fmt.Sprintf("entity.%sSliceToJSON(req.%s.%s)", template.Type, helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(f.Identifier))
	template.RepoFromMapper = fmt.Sprintf("entity.JSONTo%sSlice(model.%s)",
		template.Type,
		template.Type,
	)

	// graph
	template.GraphInType = fmt.Sprintf("[String%s]%s", template.GraphRequired, template.GraphRequired)
	template.GraphInTypeOptional = "[String]"
	template.GraphOutType = fmt.Sprintf("[String%s]%s", template.GraphRequired, template.GraphRequired)
	template.GraphGenType = "[]string"
	template.GraphGenToMapper = fmt.Sprintf("Map%s%sToModel(i.%s)",
		helpers.ToCamelCase(e.Identifier),
		pl.Singular(helpers.ToCamelCase(name)),
		helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapper = fmt.Sprintf("Map%s%sFromModel(i.%s)",
		helpers.ToCamelCase(e.Identifier),
		pl.Singular(helpers.ToCamelCase(name)),
		helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapperOptional = fmt.Sprintf("Map%s%sFromModelOptional(i.%s)",
		helpers.ToCamelCase(e.Identifier),
		pl.Singular(helpers.ToCamelCase(name)),
		helpers.ToCamelCase(f.Identifier))
	if dependantEntity != nil {
		template.GraphGenToMapper = fmt.Sprintf("Map%s%sToModel(i.%s)", pl.Singular(helpers.ToCamelCase(dependantEntity.Identifier)), pl.Singular(helpers.ToCamelCase(name)), helpers.ToCamelCase(f.Identifier))
		template.GraphGenFromMapper = fmt.Sprintf("Map%s%sFromModel(i.%s)", pl.Singular(helpers.ToCamelCase(dependantEntity.Identifier)), pl.Singular(helpers.ToCamelCase(name)), helpers.ToCamelCase(f.Identifier))
	}

	// proto
	protoType := helpers.ToCamelCase(fmt.Sprintf("%s_%s", pl.Singular(e.Identifier), pl.Singular(f.Identifier)))
	template.ProtoType = protoType
	template.ProtoEnumOptions = helpers.ProtoEnumOptions(protoType, f.OptionValues)
	template.ProtoToMapper = fmt.Sprintf("%sSliceToProto(e.%s)", protoType, helpers.ToCamelCase(f.Identifier))
	template.ProtoFromMapper = fmt.Sprintf("%sSliceFromProto(m.Get%s())", protoType, strcase.ToCamel(f.Identifier))

	return template
}
