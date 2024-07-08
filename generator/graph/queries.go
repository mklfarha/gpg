package graph

import (
	"context"
	"fmt"
	"path"
	"text/template"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/core"
	"github.com/maykel/gpg/generator/field"
)

func generateQueries(ctx context.Context,
	graphDir string,
	entityTemplates []GraphEntityTemplate,
	project entity.Project) error {

	fmt.Printf("----[GPG][GraphQL] Generating queries\n")
	err := generator.GenerateFile(ctx, generator.FileRequest{
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

	if err != nil {
		return err
	}

	fmt.Printf("----[GPG][GraphQL] Generating mutations\n")
	return generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(graphDir, "gqls", "mutations.graphqls"),
		TemplateName: path.Join("graph", "graph_mutations"),
		Data: GraphQueriesTemplate{
			Entities: entityTemplates,
		},
		DisableGoFormat: true,
	})
}
