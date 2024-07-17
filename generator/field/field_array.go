package field

import (
	"fmt"

	"github.com/maykel/gpg/entity"
)

func ArrayFieldTemplate(f entity.Field, e entity.Entity) Template {
	template := BaseFieldTemplate(f, e)

	arrayType := f.ArrayConfig.Type
	arrayTypeTemplate := ResolveFieldType(entity.Field{Type: arrayType, Required: true}, e, nil)

	//base
	template.Type = fmt.Sprintf("[]%s", arrayTypeTemplate.Type)
	template.InternalType = entity.ArrayFieldType
	template.GenFieldType = "ArrayFieldType"
	template.GenRandomValue = fmt.Sprintf("randomvalues.GetArrayValue(%s)", arrayTypeTemplate.GenFieldType)

	//graph
	template.GraphInType = fmt.Sprintf("[%s!]%s", arrayTypeTemplate.GraphInType, template.GraphRequired)
	template.GraphInTypeOptional = fmt.Sprintf("[%s!]%s", arrayTypeTemplate.GraphInType, template.GraphRequired)
	template.GraphOutType = fmt.Sprintf("[%s!]%s", arrayTypeTemplate.GraphOutType, template.GraphRequired)
	template.GraphGenType = fmt.Sprintf("[]%s", arrayTypeTemplate.GraphGenType)

	//proto
	template.ProtoType = fmt.Sprintf("repeated %s", arrayTypeTemplate.ProtoType)

	return template
}
