package field

import (
	"fmt"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func DateFieldTemplate(f entity.Field, e entity.Entity) Template {

	generatedDate := f.Autogenerated.Type == entity.InsertCurrentTimestampAutogeneratedType ||
		f.Autogenerated.Type == entity.UpdateCurrentTimestampAutogeneratedType

	graphRequired := ""
	if f.Required && !generatedDate {
		graphRequired = "!"
	}

	graphGenToMapper := fmt.Sprintf("i.%s.Format(\"2006-01-02\")", helpers.ToCamelCase(f.Identifier))
	graphGenFromMapper := fmt.Sprintf("ParseDate(i.%s)", helpers.ToCamelCase(f.Identifier))
	graphGenFromMapperOptional := fmt.Sprintf("ParseDateFromPointer(i.%s)", helpers.ToCamelCase(f.Identifier))
	if !f.Required || generatedDate {
		graphGenToMapper = fmt.Sprintf("FormatDateToPointer(i.%s)", helpers.ToCamelCase(f.Identifier))
		graphGenFromMapper = graphGenFromMapperOptional
	}

	imp := "time"
	return Template{
		Identifier:                 f.Identifier,
		Name:                       helpers.ToCamelCase(f.Identifier),
		Type:                       "time.Time",
		IsPrimary:                  f.StorageConfig.PrimaryKey,
		Required:                   f.Required,
		Tags:                       helpers.ResolveTags(f),
		Import:                     &imp,
		Custom:                     false,
		Generated:                  f.Autogenerated.Type != entity.InvalidAutogeneratedType,
		GeneratedFuncInsert:        resolveGeneratedFuncInsert(e, f),
		GeneratedFuncUpdate:        resolveGeneratedFuncUpdate(e, f),
		Enum:                       false,
		RepoFromMapper:             fmt.Sprintf("model.%s", helpers.ToCamelCase(f.Identifier)),
		GraphName:                  f.Identifier,
		GraphModelName:             helpers.ToCamelCase(f.Identifier),
		GraphInType:                fmt.Sprintf("String%s", graphRequired),
		GraphInTypeOptional:        "String",
		GraphOutType:               fmt.Sprintf("String%s", graphRequired),
		GraphGenToMapper:           graphGenToMapper,
		GraphGenFromMapper:         graphGenFromMapper,
		GraphGenFromMapperOptional: graphGenFromMapperOptional,
	}
}
