package graph

import (
	"context"
	"fmt"
	"path"

	"github.com/gertd/go-pluralize"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/core/repo"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

type graphGenEntitiesResponse struct {
	EntityTemplates     []GraphEntityTemplate
	JsonEntityTemplates []GraphEntityTemplate
	EnumTemplates       []field.Template
}

func generateEntities(ctx context.Context, graphDir string, project entity.Project) (graphGenEntitiesResponse, error) {
	pl := pluralize.NewClient()
	entityTemplates := []GraphEntityTemplate{}
	jsonEntityTemplates := []GraphEntityTemplate{}
	enumTemplates := []field.Template{}
	fmt.Printf("----[GPG][GraphQL] Generating entities gqls\n")
	for _, e := range project.Entities {
		inFields := []field.Template{}
		outFields := []field.Template{}
		searchable := false

		for _, f := range e.Fields {
			fieldTemplate := field.ResolveFieldType(f, e, nil)
			if !helpers.EntityContainsOperation(f.Hidden.API, entity.SelectOperation) {
				outFields = append(outFields, fieldTemplate)
			}
			if !helpers.EntityContainsOperation(f.Hidden.API, entity.UpsertOperation) {
				inFields = append(inFields, fieldTemplate)
			}

			if f.Type == entity.OptionsSingleFieldType || f.Type == entity.OptionsManyFieldType {
				enumTemplates = append(enumTemplates, fieldTemplate)
			}

			if f.StorageConfig.Search {
				searchable = true
			}
		}

		entityTemplate := GraphEntityTemplate{
			ProjectIdentifier: project.Identifier,
			ProjectModule:     project.Module,
			Identifier:        e.Identifier,
			EntityType:        helpers.ToCamelCase(e.Identifier),
			EntityTypePlural:  pl.Plural(helpers.ToCamelCase(e.Identifier)),
			PrimaryKey:        field.ResolveFieldType(helpers.EntityPrimaryKey(e), e, nil),
			InFields:          inFields,
			OutFields:         outFields,
			Search:            searchable,
			CustomQueries:     e.CustomQueries,
		}
		err := generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:      path.Join(graphDir, "gqls", fmt.Sprintf("model_%s.graphqls", e.Identifier)),
			TemplateName:    path.Join("graph", "graph_entity"),
			Data:            entityTemplate,
			DisableGoFormat: true,
		})
		if err != nil {
			return graphGenEntitiesResponse{}, err
		}

		selects := repo.ResolveSelectStatements(project, e)
		entityTemplate.Selects = selects
		entityTemplates = append(entityTemplates, entityTemplate)

		res, err := generateJSONEntities(ctx, graphDir, project, e)
		if err != nil {
			return graphGenEntitiesResponse{}, err
		}
		jsonEntityTemplates = append(jsonEntityTemplates, res.JsonEntityTemplates...)
		enumTemplates = append(enumTemplates, res.EnumTemplates...)
	}

	return graphGenEntitiesResponse{
		EntityTemplates:     entityTemplates,
		JsonEntityTemplates: jsonEntityTemplates,
		EnumTemplates:       enumTemplates,
	}, nil
}

func generateJSONEntities(ctx context.Context, graphDir string, project entity.Project, e entity.Entity) (graphGenEntitiesResponse, error) {
	pl := pluralize.NewClient()
	jsonEntityTemplates := []GraphEntityTemplate{}
	enumTemplates := []field.Template{}
	for _, f := range e.Fields {
		if f.Type == entity.JSONFieldType {
			ft := field.ResolveFieldType(f, e, &field.Template{
				Identifier: f.Identifier,
			})
			if !f.JSONConfig.Reuse && len(f.JSONConfig.Fields) > 0 {
				fields, _ := field.ResolveFieldsAndImports(project, f.JSONConfig.Fields, e, &ft)
				jsonEntityTemplate := GraphEntityTemplate{
					ProjectIdentifier: project.Identifier,
					ProjectModule:     project.Module,
					Identifier:        f.JSONConfig.Identifier,
					EntityType:        ft.Type,
					EntityTypePlural:  pl.Plural(ft.Type),
					JSON:              true,
					JSONMany:          f.JSONConfig.Type == entity.ManyJSONConfigType,
					Required:          f.Required,
					GraphGenType:      ft.GraphGenType,
					ParentIdentifier:  e.Identifier,
					ParentEntityName:  helpers.ToCamelCase(e.Identifier),
					InFields:          fields,
					OutFields:         fields,
				}
				generator.GenerateFile(ctx, generator.FileRequest{
					OutputFile:      path.Join(graphDir, "gqls", fmt.Sprintf("model_%s.graphqls", f.JSONConfig.Identifier)),
					TemplateName:    path.Join("graph", "graph_entity"),
					Data:            jsonEntityTemplate,
					DisableGoFormat: true,
				})
				jsonEntityTemplates = append(jsonEntityTemplates, jsonEntityTemplate)
				for _, fn := range fields {
					if fn.InternalType == entity.OptionsSingleFieldType || fn.InternalType == entity.OptionsManyFieldType {
						enumTemplates = append(enumTemplates, fn)
					}

				}
			}

			if f.HasNestedJsonFields() {
				res, err := generateJSONEntities(ctx, graphDir, project, entity.Entity{
					Identifier: f.JSONConfig.Identifier,
					Fields:     f.JSONConfig.Fields,
				})
				if err != nil {
					return graphGenEntitiesResponse{}, err
				}
				jsonEntityTemplates = append(jsonEntityTemplates, res.JsonEntityTemplates...)
				enumTemplates = append(enumTemplates, res.EnumTemplates...)
			}
		}
	}

	return graphGenEntitiesResponse{
		EntityTemplates:     []GraphEntityTemplate{},
		JsonEntityTemplates: jsonEntityTemplates,
		EnumTemplates:       enumTemplates,
	}, nil
}
