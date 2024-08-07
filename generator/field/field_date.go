package field

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func DateFieldTemplate(f entity.Field, e entity.Entity) Template {

	template := BaseFieldTemplate(f, e)

	//base
	template.Type = "time.Time"
	template.InternalType = entity.DateFieldType
	template.GenFieldType = "TimestampFieldType"
	template.GenRandomValue = "randomvalues.GetRandomTimeValue()"
	imp := "time"
	template.Import = &imp

	//graph
	template.GraphInType = fmt.Sprintf("String%s", template.GraphRequired)
	template.GraphInTypeOptional = "String"
	template.GraphOutType = fmt.Sprintf("String%s", template.GraphRequired)
	template.GraphGenToMapper = fmt.Sprintf("i.%s.Format(\"2006-01-02\")", helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapper = fmt.Sprintf("ParseDate(i.%s)", helpers.ToCamelCase(f.Identifier))
	template.GraphGenFromMapperOptional = fmt.Sprintf("ParseDateFromPointer(i.%s)", helpers.ToCamelCase(f.Identifier))
	if !f.Required {
		template.GraphGenToMapper = fmt.Sprintf("FormatDateToPointer(i.%s)", helpers.ToCamelCase(f.Identifier))
		template.GraphGenFromMapper = template.GraphGenFromMapperOptional
	}

	//proto
	template.ProtoType = "google.protobuf.Timestamp"
	template.ProtoToMapper = fmt.Sprintf("timestamppb.New(e.%s)", helpers.ToCamelCase(f.Identifier))
	template.ProtoFromMapper = fmt.Sprintf("m.Get%s().AsTime()", strcase.ToCamel(f.Identifier))

	return template

}
