package field

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func UUIDFieldTemplate(f entity.Field, e entity.Entity) Template {

	template := BaseFieldTemplate(f, e)

	//base
	template.Type = "uuid.UUID"
	template.InternalType = entity.UUIDFieldType
	template.GenFieldType = "UUIDFieldType"
	template.GenRandomValue = "randomvalues.GetRandomUUIDValue()"
	imp := "github.com/gofrs/uuid"
	template.Import = &imp

	if f.Required {
		template.RepoToMapper = ".String()"
		template.RepoFromMapper = fmt.Sprintf("uuid.FromStringOrNil(model.%s)", helpers.ToCamelCase(f.Identifier))
	} else {
		template.RepoToMapper = fmt.Sprintf("mapper.StringToSqlNullString(req.%s.%s.String())", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(f.Identifier))
		template.RepoToMapperFetch = fmt.Sprintf("mapper.StringToSqlNullString(req.%s.String())", helpers.ToCamelCase(f.Identifier))
		template.RepoFromMapper = fmt.Sprintf("uuid.FromStringOrNil(mapper.SqlNullStringToString(model.%s))", helpers.ToCamelCase(f.Identifier))
	}

	//graph
	template.GraphName = strings.ReplaceAll(f.Identifier, "Uuid", "UUID")
	graphModelName := strings.ReplaceAll(helpers.ToCamelCase(f.Identifier), "Id", "ID")
	graphModelName = strings.ReplaceAll(graphModelName, "Uuid", "UUID")
	template.GraphModelName = graphModelName
	template.GraphInType = fmt.Sprintf("ID%s", template.GraphRequired)
	template.GraphInTypeOptional = "ID"
	template.GraphOutType = fmt.Sprintf("ID%s", template.GraphRequired)
	template.GraphGenType = "string"
	template.GraphGenToMapper = fmt.Sprintf("i.%s.String()", helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapper = fmt.Sprintf("UuidFromStringOrNil(i.%s)", graphModelName)
	template.GraphGenFromMapperOptional = fmt.Sprintf("UuidFromPointerString(i.%s)", graphModelName)
	template.GraphGenFromMapperParam = fmt.Sprintf("mapper.UuidFromStringOrNil(%s)", f.Identifier)

	if !f.Required {
		template.GraphGenToMapper = fmt.Sprintf("UuidToPointerString(i.%s)", helpers.ToCamelCase(f.Identifier))
		template.GraphGenFromMapper = template.GraphGenFromMapperOptional
		template.GraphGenFromMapperParam = fmt.Sprintf("mapper.UuidFromPointerString(%s)", f.Identifier)
		template.GraphGenType = "*string"
	}

	if f.StorageConfig.PrimaryKey {
		template.GraphInType = "ID"
		template.GraphGenFromMapper = fmt.Sprintf("UuidFromPointerString(i.%s)", graphModelName)
		template.GraphGenFromMapperParam = fmt.Sprintf("mapper.UuidFromPointerString(%s)", f.Identifier)
		template.GraphGenType = "*string"
	}

	//proto
	template.ProtoType = "string"
	template.ProtoToMapper = fmt.Sprintf("e.%s.String()", helpers.ToCamelCase(f.Identifier))
	template.ProtoFromMapper = fmt.Sprintf("uuid.FromStringOrNil(m.Get%s())", strcase.ToCamel(f.Identifier))

	return template

}
