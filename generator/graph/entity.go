package graph

import (
	"context"
	"fmt"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/core"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

type graphGenEntitiesResponse struct {
	EntityTemplates     []GraphEntityTemplate
	JsonEntityTemplates []GraphEntityTemplate
	EnumTemplates       []field.Template
}

func generateEntities(ctx context.Context, graphDir string, project entity.Project) (graphGenEntitiesResponse, error) {
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

			if f.Type == entity.JSONFieldType {
				ft := field.ResolveFieldType(f, e, &field.Template{
					Identifier: f.Identifier,
				})
				if len(f.JSONConfig.Fields) > 0 {
					fields, _ := field.ResolveFieldsAndImports(f.JSONConfig.Fields, e, &ft)
					jsonEntityTemplate := GraphEntityTemplate{
						Identifier:       f.Identifier,
						EntityType:       ft.Type,
						JSON:             true,
						JSONMany:         f.JSONConfig.Type == entity.ManyJSONConfigType,
						Required:         f.Required,
						GraphGenType:     ft.GraphGenType,
						ParentIdentifier: e.Identifier,
						ParentEntityName: helpers.ToCamelCase(e.Identifier),
						InFields:         fields,
						OutFields:        fields,
					}
					generator.GenerateFile(ctx, generator.FileRequest{
						OutputFile:      path.Join(graphDir, "gqls", fmt.Sprintf("model_%s_%s.graphqls", e.Identifier, f.Identifier)),
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
			}

			if f.Type == entity.OptionsSingleFieldType || f.Type == entity.OptionsManyFieldType {
				enumTemplates = append(enumTemplates, fieldTemplate)
			}

			if f.StorageConfig.Search {
				searchable = true
			}
		}
		entityTemplate := GraphEntityTemplate{
			Identifier:    e.Identifier,
			EntityType:    helpers.ToCamelCase(e.Identifier),
			PrimaryKey:    field.ResolveFieldType(helpers.EntityPrimaryKey(e), e, nil),
			InFields:      inFields,
			OutFields:     outFields,
			Search:        searchable,
			CustomQueries: e.CustomQueries,
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

		selects := core.ResolveSelectStatements(project, e)
		entityTemplate.Selects = selects
		entityTemplates = append(entityTemplates, entityTemplate)
	}

	return graphGenEntitiesResponse{
		EntityTemplates:     entityTemplates,
		JsonEntityTemplates: jsonEntityTemplates,
		EnumTemplates:       enumTemplates,
	}, nil
}