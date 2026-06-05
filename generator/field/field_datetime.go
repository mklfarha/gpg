package field

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func DatetimeFieldTemplate(f entity.Field, e entity.Entity) Template {

	template := BaseFieldTemplate(f, e)

	//base
	template.InternalType = entity.DateTimeFieldType
	template.GenFieldType = "TimestampFieldType"
	imp := "time"
	template.Import = &imp

	if f.Required {
		template.Type = "time.Time"
		template.GenRandomValue = "randomvalues.GetRandomTimeValue()"
	} else {
		template.Type = "*time.Time"
		template.GenRandomValue = "randomvalues.GetRandomTimeValuePtr()"
		template.RepoToMapper = fmt.Sprintf("mapper.TimePtrToSqlNullTime(req.%s.%s)", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(f.Identifier))
		template.RepoToMapperFetch = fmt.Sprintf("mapper.TimePtrToSqlNullTime(req.%s)", helpers.ToCamelCase(f.Identifier))
		template.RepoFromMapper = fmt.Sprintf("mapper.SqlNullTimeToTimePtr(model.%s)", template.Name)
	}

	//graph
	template.GraphInType = fmt.Sprintf("String%s", template.GraphRequired)
	template.GraphInTypeOptional = "String"
	template.GraphOutType = fmt.Sprintf("String%s", template.GraphRequired)

	template.GraphGenToMapper = fmt.Sprintf("i.%s.Format(\"2006-01-02 15:04:05\")", helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapper = fmt.Sprintf("ParseTime(i.%s)", helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapperOptional = fmt.Sprintf("ParseTimeFromPointer(i.%s)", helpers.ToCamelCase(f.Identifier))
	if !f.Required {
		template.GraphGenToMapper = fmt.Sprintf("FormatTimeToPointer(i.%s)", helpers.ToCamelCase(f.Identifier))
		template.GraphGenFromMapper = template.GraphGenFromMapperOptional
	}

	//proto
	template.ProtoType = "google.protobuf.Timestamp"
	if f.Required {
		template.ProtoToMapper = fmt.Sprintf("timestamppb.New(e.%s)", helpers.ToCamelCase(f.Identifier))
		template.ProtoFromMapper = fmt.Sprintf("m.Get%s().AsTime()", strcase.ToCamel(f.Identifier))
	} else {
		template.ProtoToMapper = fmt.Sprintf("TimePtrToTimestamppb(e.%s)", helpers.ToCamelCase(f.Identifier))
		template.ProtoFromMapper = fmt.Sprintf("TimestamppbToTimePtr(m.Get%s())", strcase.ToCamel(f.Identifier))
	}

	return template
}
