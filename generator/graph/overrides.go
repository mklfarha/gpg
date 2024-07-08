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

func overrideGqlgenFiles(ctx context.Context,
	graphDir string,
	entityTemplates []GraphEntityTemplate,
	project entity.Project) error {

	fmt.Printf("----[GPG][GraphQL] Override resolver\n")
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(graphDir, "resolver.go"),
		TemplateName: path.Join("graph", "graph_resolver"),
		Data:         project,
	})
	if err != nil {
		return err
	}

	fmt.Printf("----[GPG][GraphQL] Override queries resolvers\n")
	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(graphDir, "queries.resolvers.go"),
		TemplateName: path.Join("graph", "graph_queries_resolver"),
		Data: GraphQueriesTemplate{
			ProjectName: project.Identifier,
			Project:     project,
			Entities:    entityTemplates,
		},
		Funcs: template.FuncMap{
			"CustomQueryInputFields": func(cq entity.CustomQuery) map[string]field.Template {
				inputFields, _, _ := core.GetCustomQueryFields(cq.Condition, project)
				return inputFields
			},
		},
		DisableGoFormat: false,
	})

	if err != nil {
		return err
	}

	fmt.Printf("----[GPG][GraphQL] Override mutations resolvers\n")
	return generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(graphDir, "mutations.resolvers.go"),
		TemplateName: path.Join("graph", "graph_mutations_resolver"),
		Data: GraphQueriesTemplate{
			ProjectName: project.Identifier,
			Entities:    entityTemplates,
		},
	})
}
