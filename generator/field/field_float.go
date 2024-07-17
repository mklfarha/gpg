package field

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func FloatFieldTemplate(f entity.Field, e entity.Entity) Template {

	template := BaseFieldTemplate(f, e)

	//base
	template.Type = "float64"
	template.InternalType = entity.FloatFieldType
	template.GenFieldType = "FloatFieldType"
	template.GenRandomValue = "randomvalues.GetRandomFloatValue()"

	//graph
	template.GraphInType = fmt.Sprintf("Float%s", template.GraphRequired)
	template.GraphInTypeOptional = "Float"
	template.GraphOutType = fmt.Sprintf("Float%s", template.GraphRequired)
	template.GraphGenType = "float64"
	template.GraphGenToMapper = fmt.Sprintf("i.%s", helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapper = fmt.Sprintf("i.%s", helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapperOptional = fmt.Sprintf("FloatFromPointer(i.%s)", helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapperParam = f.Identifier
	if !f.Required {
		template.GraphGenType = "*float64"
		template.GraphGenToMapper = fmt.Sprintf("&i.%s", helpers.ToCamelCase(f.Identifier))
		template.GraphGenFromMapper = template.GraphGenFromMapperOptional
		template.GraphGenFromMapperParam = fmt.Sprintf("mapper.FloatFromPointer(%s)", f.Identifier)

	}

	//proto
	template.ProtoType = "double"
	template.ProtoToMapper = fmt.Sprintf("e.%s", strcase.ToCamel(f.Identifier))
	template.ProtoFromMapper = fmt.Sprintf("m.%s", strcase.ToCamel(f.Identifier))

	return template
}
