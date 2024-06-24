package graph

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path"
	"text/template"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/core"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

type GraphEntityTemplate struct {
	Identifier       string
	EntityName       string
	JSON             bool
	JSONMany         bool
	Required         bool
	ParentIdentifier string
	ParentEntityName string
	GraphGenType     string
	PrimaryKey       field.Template
	InFields         []field.Template
	OutFields        []field.Template
	Selects          []core.RepoSchemaSelectStatement
	CustomQueries    []entity.CustomQuery
	Search           bool
}

type GraphQueriesTemplate struct {
	ProjectName  string
	Entities     []GraphEntityTemplate
	JSONEntities []GraphEntityTemplate
	Enums        []field.Template
}

func GenerateGraph(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG][GraphQL] Generating GraphQL\n")
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	graphDir := path.Join(projectDir, generator.GRAPH_DIR)

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(graphDir, "gqlgen.yml"),
		TemplateName:    path.Join("graph", "graph_yaml"),
		DisableGoFormat: true,
	})

	pl := pluralize.NewClient()

	entityTemplates := []GraphEntityTemplate{}
	jsonEntityTemplates := []GraphEntityTemplate{}
	enumTemplates := []field.Template{}
	fmt.Printf("----[GPG][GraphQL] Generating gqls\n")
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
				if len(f.JSONConfig.Fields) > 0 {
					fields, _ := field.ResolveFieldsAndImports(f.JSONConfig.Fields, e, &f.Identifier)
					jsonEntityTemplate := GraphEntityTemplate{
						Identifier:       f.Identifier,
						EntityName:       pl.Singular(helpers.ToCamelCase(f.Identifier)),
						JSON:             true,
						JSONMany:         f.JSONConfig.Type == entity.ManyJSONConfigType,
						Required:         f.Required,
						GraphGenType:     fmt.Sprintf("%s%s", helpers.ToCamelCase(e.Identifier), pl.Singular(helpers.ToCamelCase(f.Identifier))),
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
			EntityName:    helpers.ToCamelCase(e.Identifier),
			PrimaryKey:    field.ResolveFieldType(helpers.EntityPrimaryKey(e), e, nil),
			InFields:      inFields,
			OutFields:     outFields,
			Search:        searchable,
			CustomQueries: e.CustomQueries,
		}
		generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:      path.Join(graphDir, "gqls", fmt.Sprintf("model_%s.graphqls", e.Identifier)),
			TemplateName:    path.Join("graph", "graph_entity"),
			Data:            entityTemplate,
			DisableGoFormat: true,
		})

		selects := core.ResolveSelectStatements(project, e)
		entityTemplate.Selects = selects
		entityTemplates = append(entityTemplates, entityTemplate)
	}

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(graphDir, "gqls", "queries.graphqls"),
		TemplateName: path.Join("graph", "graph_queries"),
		Data: GraphQueriesTemplate{
			Entities: entityTemplates,
		},
		Funcs: template.FuncMap{
			"CustomQueryInputFields": func(cq entity.CustomQuery) map[string]field.Template {
				inputFields, _, _ := core.GetCustomQueryFields(cq.Condition, project)
				return inputFields
			},
		},
		DisableGoFormat: true,
	})
	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(graphDir, "gqls", "mutations.graphqls"),
		TemplateName: path.Join("graph", "graph_mutations"),
		Data: GraphQueriesTemplate{
			Entities: entityTemplates,
		},
		DisableGoFormat: true,
	})

	fmt.Printf("----[GPG][GraphQL] GQLGEN generate\n")
	cmd := exec.Command("go", "run", "github.com/99designs/gqlgen", "generate")
	cmd.Dir = graphDir
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		fmt.Println("gqlgen result: " + out.String())
	}

	fmt.Printf("----[GPG][GraphQL] Override resolver\n")
	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(graphDir, "resolver.go"),
		TemplateName: path.Join("graph", "graph_resolver"),
		Data:         project,
	})

	fmt.Printf("----[GPG][GraphQL] Mapper\n")
	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(graphDir, "mapper", "mapper.go"),
		TemplateName: path.Join("graph", "graph_mapper"),
		Data: GraphQueriesTemplate{
			ProjectName:  project.Identifier,
			Entities:     entityTemplates,
			JSONEntities: jsonEntityTemplates,
			Enums:        enumTemplates,
		},
	})

	fmt.Printf("----[GPG][GraphQL] Override queries resolvers\n")
	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(graphDir, "queries.resolvers.go"),
		TemplateName: path.Join("graph", "graph_queries_resolver"),
		Data: GraphQueriesTemplate{
			ProjectName: project.Identifier,
			Entities:    entityTemplates,
		},
		Funcs: template.FuncMap{
			"CustomQueryInputFields": func(cq entity.CustomQuery) map[string]field.Template {
				inputFields, _, _ := core.GetCustomQueryFields(cq.Condition, project)
				return inputFields
			},
		},
		DisableGoFormat: true,
	})

	fmt.Printf("----[GPG][GraphQL] Override mutations resolvers\n")
	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(graphDir, "mutations.resolvers.go"),
		TemplateName: path.Join("graph", "graph_mutations_resolver"),
		Data: GraphQueriesTemplate{
			ProjectName: project.Identifier,
			Entities:    entityTemplates,
		},
	})

	fmt.Printf("----[GPG][GraphQL] Server\n")
	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(graphDir, "server.go"),
		TemplateName: path.Join("graph", "graph_server"),
		Data:         project,
	})

	return nil
}
