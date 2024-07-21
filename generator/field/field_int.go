package field

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func IntFieldTemplate(f entity.Field, e entity.Entity) Template {

	template := BaseFieldTemplate(f, e)

	//base
	template.Type = "int64"
	template.InternalType = entity.IntFieldType
	template.GenFieldType = "IntFieldType"
	template.GenRandomValue = "randomvalues.GetRandomIntValue()"

	//graph
	template.GraphInType = fmt.Sprintf("Int%s", template.GraphRequired)
	template.GraphInTypeOptional = "Int"
	template.GraphOutType = fmt.Sprintf("Int%s", template.GraphRequired)
	template.GraphGenType = "int64"
	template.GraphGenToMapper = fmt.Sprintf("i.%s", helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapper = fmt.Sprintf("int64(i.%s)", helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapperOptional = fmt.Sprintf("IntFromPointer(i.%s)", helpers.ToCamelCase(f.Identifier))
	if !f.Required {
		template.GraphGenToMapper = fmt.Sprintf("IntPointer(i.%s)", helpers.ToCamelCase(f.Identifier))
		template.GraphGenFromMapper = template.GraphGenFromMapperOptional
	}

	//proto
	template.ProtoType = "int64"
	template.ProtoToMapper = fmt.Sprintf("int64(e.%s)", helpers.ToCamelCase(f.Identifier))
	template.ProtoFromMapper = fmt.Sprintf("int64(m.Get%s())", strcase.ToCamel(f.Identifier))

	return template
}
