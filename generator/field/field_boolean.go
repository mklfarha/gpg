package field

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func BooleanFieldTemplate(f entity.Field, e entity.Entity) Template {

	template := BaseFieldTemplate(f, e)

	//base
	template.Type = "bool"
	template.InternalType = entity.BooleanFieldType
	template.GenFieldType = "BooleanFieldType"
	template.GenRandomValue = "randomvalues.GetRandomBoolValue()"

	//graph
	template.GraphInType = fmt.Sprintf("Boolean%s", template.GraphRequired)
	template.GraphInTypeOptional = "Boolean"
	template.GraphOutType = fmt.Sprintf("Boolean%s", template.GraphRequired)
	template.GraphGenType = "bool"
	template.GraphGenToMapper = fmt.Sprintf("&i.%s", helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapperParam = f.Identifier
	template.GraphGenFromMapper = fmt.Sprintf("i.%s", helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapperOptional = fmt.Sprintf("BoolFromPointer(i.%s)", helpers.ToCamelCase(f.Identifier))
	if f.Required {
		template.GraphGenToMapper = fmt.Sprintf("i.%s", helpers.ToCamelCase(f.Identifier))
	}

	//proto
	template.ProtoType = "bool"
	template.ProtoToMapper = fmt.Sprintf("e.%s", helpers.ToCamelCase(f.Identifier))
	template.ProtoFromMapper = fmt.Sprintf("m.Get%s()", strcase.ToCamel(f.Identifier))

	return template
}
