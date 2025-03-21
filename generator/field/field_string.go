package field

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func StringFieldTemplate(f entity.Field, e entity.Entity) Template {

	template := BaseFieldTemplate(f, e)

	//base
	if f.Required {
		template.Type = "string"
		template.GenRandomValue = "randomvalues.GetRandomStringValue()"
	} else {
		template.Type = "*string"
		template.GenRandomValue = "randomvalues.GetRandomStringValuePtr()"
		template.RepoToMapper = fmt.Sprintf("mapper.StringPtrToSqlNullString(req.%s.%s)", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(f.Identifier))
		template.RepoToMapperFetch = fmt.Sprintf("mapper.StringPtrToSqlNullString(req.%s)", helpers.ToCamelCase(f.Identifier))
		template.RepoFromMapper = fmt.Sprintf("mapper.SqlNullStringToStringPtr(model.%s)", template.Name)
	}

	template.InternalType = entity.StringFieldType
	template.GenFieldType = "StringFieldType"

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
	if !f.Required {
		template.ProtoToMapper = fmt.Sprintf("StringPtrToString(e.%s)", helpers.ToCamelCase(f.Identifier))
		template.ProtoFromMapper = fmt.Sprintf("&m.%s", strcase.ToCamel(f.Identifier))
	} else {
		template.ProtoToMapper = fmt.Sprintf("e.%s", helpers.ToCamelCase(f.Identifier))
		template.ProtoFromMapper = fmt.Sprintf("m.Get%s()", strcase.ToCamel(f.Identifier))
	}

	return template
}
