package field

import (
	"fmt"

	"github.com/maykel/gpg/entity"
)

func StringFieldTemplate(f entity.Field, e entity.Entity) Template {

	template := BaseFieldTemplate(f, e)

	//base
	template.Type = "string"
	template.InternalType = entity.StringFieldType
	template.GenFieldType = "StringFieldType"
	template.GenRandomValue = "randomvalues.GetRandomStringValue()"

	// graph
	template.GraphInType = fmt.Sprintf("String%s", template.GraphRequired)
	template.GraphInTypeOptional = "String"
	template.GraphOutType = fmt.Sprintf("String%s", template.GraphRequired)
	template.GraphGenType = "string"
	template.GraphGenFromMapperOptional = fmt.Sprintf("StringFromPointer(i.%s)", template.GraphModelName)
	if !f.Required {
		template.GraphGenFromMapper = template.GraphGenFromMapperOptional
		template.GraphGenFromMapperParam = fmt.Sprintf("mapper.StringFromPointer(%s)", f.Identifier)
		template.GraphGenType = "*string"
	}

	//proto
	template.ProtoType = "string"

	return template
}
