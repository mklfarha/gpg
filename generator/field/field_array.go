package field

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func ArrayFieldTemplate(f entity.Field, e entity.Entity) Template {
	template := BaseFieldTemplate(f, e)

	arrayType := f.ArrayConfig.Type
	arrayTypeTemplate := ResolveFieldType(entity.Field{Type: arrayType, Required: true}, e, nil)

	//base
	template.Type = fmt.Sprintf("[]%s", arrayTypeTemplate.Type)
	template.Import = arrayTypeTemplate.Import
	template.InternalType = entity.ArrayFieldType
	template.GenFieldType = "ArrayFieldType"
	template.ArrayInternalType = arrayTypeTemplate.InternalType
	template.ArrayGenFieldType = arrayTypeTemplate.GenFieldType
	template.GenRandomValue = fmt.Sprintf("[]%s{}", arrayTypeTemplate.Type)
	template.RepoFromMapper = fmt.Sprintf("mapper.JSONTo%sSlice(%s)",
		helpers.ToCamelCase(arrayTypeTemplate.InternalType.String()),
		template.RepoFromMapper,
	)
	template.RepoToMapper = fmt.Sprintf("mapper.SliceToJSON(req.%s.%s)", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(f.Identifier))
	template.Array = true

	//graph
	template.GraphInType = fmt.Sprintf("[%s]%s", arrayTypeTemplate.GraphInType, template.GraphRequired)
	template.GraphInTypeOptional = fmt.Sprintf("[%s]%s", arrayTypeTemplate.GraphInType, template.GraphRequired)
	template.GraphOutType = fmt.Sprintf("[%s]%s", arrayTypeTemplate.GraphOutType, template.GraphRequired)
	template.GraphGenType = fmt.Sprintf("[]%s", arrayTypeTemplate.GraphGenType)
	if arrayTypeTemplate.InternalType == entity.UUIDFieldType ||
		arrayTypeTemplate.InternalType == entity.IntFieldType ||
		arrayTypeTemplate.InternalType == entity.DateFieldType ||
		arrayTypeTemplate.InternalType == entity.DateTimeFieldType {
		// entity to model
		template.GraphGenToMapper = fmt.Sprintf("Map%sSlice(i.%s)",
			helpers.ToCamelCase(arrayTypeTemplate.InternalType.String()),
			helpers.ToCamelCase(f.Identifier))
		// model to entity
		template.GraphGenFromMapperParam = ""
		template.GraphGenFromMapper = fmt.Sprintf("MapTo%sSlice(i.%s)",
			helpers.ToCamelCase(arrayTypeTemplate.InternalType.String()),
			strcase.ToCamel(f.Identifier))
		template.GraphGenFromMapperOptional = template.GraphGenFromMapper
	}

	//proto
	template.ProtoType = fmt.Sprintf("repeated %s", arrayTypeTemplate.ProtoType)
	if arrayTypeTemplate.InternalType == entity.UUIDFieldType {
		template.ProtoToMapper = fmt.Sprintf("UUIDSliceToStringSlice(e.%s)", helpers.ToCamelCase(f.Identifier))
		template.ProtoFromMapper = fmt.Sprintf("StringSliceToUUIDSlice(m.Get%s())", strcase.ToCamel(f.Identifier))
	}
	if arrayTypeTemplate.InternalType == entity.DateFieldType || arrayTypeTemplate.InternalType == entity.DateTimeFieldType {
		template.ProtoToMapper = fmt.Sprintf("TimeSliceToProtoTimeSlice(e.%s)", helpers.ToCamelCase(f.Identifier))
		template.ProtoFromMapper = fmt.Sprintf("ProtoTimeSliceToTimeSlice(m.Get%s())", strcase.ToCamel(f.Identifier))
	}

	return template
}
