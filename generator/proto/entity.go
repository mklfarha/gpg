package proto

import (
	"context"
	"fmt"
	"path"
	"text/template"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

func generateEntityProtoFile(
	ctx context.Context,
	protoDir string,
	project entity.Project,
	e entity.Entity,
	parentIdentifier string,
	dependantEntity *field.Template) (ProtoEntityTemplate, []entity.Field, error) {

	fields := []field.Template{}
	searchable := false
	enums := map[string]ProtoEnumTemplate{}
	nested := []entity.Field{}
	entityTemplate := ProtoEntityTemplate{}
	var err error
	pl := pluralize.NewClient()
	imports := map[string]interface{}{}
	if len(e.Fields) > 0 {
		for _, f := range e.Fields {
			fieldTemplate := field.ResolveFieldType(f, e, dependantEntity)
			fields = append(fields, fieldTemplate)
			for _, field := range fields {
				if field.Enum {
					if _, found := enums[field.ProtoType]; !found {
						enums[field.ProtoType] = ProtoEnumTemplate{
							Field:   field,
							Many:    field.EnumMany,
							Options: field.ProtoEnumOptions,
						}
					}
				}
				if f.Type == entity.JSONFieldType {
					if len(f.JSONConfig.Fields) > 0 {
						nested = append(nested, f)
						imports[fmt.Sprintf("%s.proto", strcase.ToSnake(fieldTemplate.ProtoType))] = nil
					}
				}
				if f.Type == entity.DateFieldType || f.Type == entity.DateTimeFieldType {
					imports["google/protobuf/timestamp.proto"] = nil
				}
			}

			if f.StorageConfig.Search {
				searchable = true
			}
		}
		primaryKey := field.ResolveFieldType(helpers.EntityPrimaryKey(e), e, nil)

		finalIdentifier := strcase.ToSnake(e.Identifier)
		if dependantEntity != nil && dependantEntity.EntityIdentifier != "" {
			finalIdentifier = fmt.Sprintf("%s_%s", dependantEntity.EntityIdentifier, dependantEntity.SingularIdentifier)
		}
		entityTemplate = ProtoEntityTemplate{
			ProjectIdentifier:     project.Identifier,
			ParentIdentifier:      parentIdentifier,
			OrignalIdentifier:     e.Identifier,
			FinalIdentifier:       finalIdentifier,
			FinalIdentifierPlural: pl.Plural(finalIdentifier),
			Name:                  helpers.ToCamelCase(finalIdentifier),
			NamePlural:            pl.Plural(helpers.ToCamelCase(finalIdentifier)),
			Type:                  dependantEntity.Type,
			Fields:                fields,
			PrimaryKey:            primaryKey,
			Search:                searchable,
			Enums:                 enums,
			Imports:               imports,
		}

		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:      path.Join(protoDir, "proto", fmt.Sprintf("%s.proto", finalIdentifier)),
			TemplateName:    path.Join("proto", "proto_model"),
			Data:            entityTemplate,
			DisableGoFormat: true,
			Funcs: template.FuncMap{
				"Inc": helpers.Inc,
			},
		})
	}
	return entityTemplate, nested, err
}

func generateProtoFiles(ctx context.Context, protoDir string, project entity.Project) (standaloneEntityTemplates []ProtoEntityTemplate, dependantEntityTemplates map[string][]ProtoEntityTemplate, returnErr error) {

	standaloneEntityTemplates = []ProtoEntityTemplate{}
	dependantEntityTemplates = map[string][]ProtoEntityTemplate{}
	//generate entities/models
	fmt.Printf("--[GPG][Proto] Generating Entities\n")
	for _, e := range project.Entities {
		template, nested, err := generateEntityProtoFile(ctx, protoDir, project, e, e.Identifier, &field.Template{
			Identifier: e.Identifier,
			Type:       helpers.ToCamelCase(e.Identifier),
		})
		if err != nil {
			returnErr = err
			return
		}
		standaloneEntityTemplates = append(standaloneEntityTemplates, template)
		nestedTemplates := []ProtoEntityTemplate{}
		for _, n := range nested {
			ft := field.ResolveFieldType(n, e, &field.Template{
				Identifier: n.Identifier,
			})
			nestedTemplate, _, err := generateEntityProtoFile(ctx, protoDir, project, entity.Entity{
				Identifier: ft.Identifier,
				Fields:     n.JSONConfig.Fields,
			}, e.Identifier, &ft)
			if err != nil {
				returnErr = err
				return
			}
			nestedTemplates = append(nestedTemplates, nestedTemplate)
		}
		dependantEntityTemplates[template.OrignalIdentifier] = nestedTemplates
	}

	//generate project service definition
	fmt.Printf("--[GPG][Proto] Generating Service Definition\n")
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(protoDir, "proto", fmt.Sprintf("service_%s.proto", project.Identifier)),
		TemplateName: path.Join("proto", "proto_service"),
		Data: ProtoServiceTemplate{
			Identifier: project.Identifier,
			Name:       helpers.ToCamelCase(project.Identifier),
			Entities:   standaloneEntityTemplates,
		},
		DisableGoFormat: true,
		Funcs: template.FuncMap{
			"Inc": helpers.Inc,
		},
	})

	if err != nil {
		returnErr = err
		return
	}

	return
}
