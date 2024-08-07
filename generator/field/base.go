package field

import (
	"fmt"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/helpers"
)

func BaseFieldTemplate(f entity.Field, e entity.Entity) Template {
	pl := pluralize.NewClient()
	generatedFuncInsert, generatedFuncInsertCustom := resolveGeneratedFuncInsert(e, f)
	generatedFuncUpdate, generatedFuncUpdateCustom := resolveGeneratedFuncUpdate(e, f)

	graphModelName := strings.ReplaceAll(strcase.ToCamel(f.Identifier), "Url", "URL")
	if graphModelName == "Uuid" {
		graphModelName = "UUID"
	}
	graphModelName = strings.ReplaceAll(graphModelName, "Json", "JSON")
	graphModelName = strings.ReplaceAll(graphModelName, "Https", "HTTPS")
	graphModelName = strings.ReplaceAll(graphModelName, "Http", "HTTP")
	graphGenFromMapper := fmt.Sprintf("i.%s", graphModelName)
	graphGenToMapper := fmt.Sprintf("i.%s", helpers.ToCamelCase(f.Identifier))
	if !f.Required && f.Type != entity.ArrayFieldType {
		graphGenToMapper = fmt.Sprintf("&i.%s", helpers.ToCamelCase(f.Identifier))
	}

	graphRequired := ""
	if f.Required {
		graphRequired = "!"
	}

	finalName := strings.ReplaceAll(helpers.ToCamelCase(f.Identifier), "Json", "JSON")
	return Template{
		Identifier:         f.Identifier,
		SingularIdentifier: pl.Singular(f.Identifier),
		Name:               finalName,
		EntityIdentifier:   e.Identifier,
		IsPrimary:          f.StorageConfig.PrimaryKey,
		Required:           f.Required,
		Tags:               helpers.ResolveTags(f),
		Custom:             false,

		Generated:             f.Autogenerated.Type != entity.InvalidAutogeneratedType,
		GeneratedFuncInsert:   generatedFuncInsert,
		GeneratedInsertCustom: generatedFuncInsertCustom,
		GeneratedFuncUpdate:   generatedFuncUpdate,
		GeneratedUpdateCustom: generatedFuncUpdateCustom,

		RepoToMapper:   "",
		RepoFromMapper: fmt.Sprintf("model.%s", helpers.ToCamelCase(f.Identifier)),

		GraphRequired:              graphRequired,
		GraphName:                  f.Identifier,
		GraphModelName:             graphModelName,
		GraphGenToMapper:           graphGenToMapper,
		GraphGenFromMapper:         graphGenFromMapper,
		GraphGenFromMapperOptional: graphGenFromMapper,
		GraphGenFromMapperParam:    f.Identifier,

		ProtoName:       helpers.ToSnakeCase(f.Identifier),
		ProtoToMapper:   fmt.Sprintf("e.%s", helpers.ToCamelCase(f.Identifier)),
		ProtoFromMapper: fmt.Sprintf("m.Get%s()", strcase.ToCamel(f.Identifier)),
		ProtoGenName:    strcase.ToCamel(f.Identifier),
	}
}
