package graph

import (
	"context"
	"fmt"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
)

func generateMapper(ctx context.Context, graphDir string, project entity.Project, res graphGenEntitiesResponse) error {
	fmt.Printf("----[GPG][GraphQL] Mapper\n")
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(graphDir, "mapper", "mapper.go"),
		TemplateName: path.Join("graph", "graph_mapper"),
	})
	if err != nil {
		return err
	}

	for _, e := range res.EntityTemplates {
		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(graphDir, "mapper", fmt.Sprintf("mapper_%s.go", e.Identifier)),
			TemplateName: path.Join("graph", "graph_mapper_entity"),
			Data:         e,
		})
		if err != nil {
			return err
		}
	}

	for _, e := range res.JsonEntityTemplates {
		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(graphDir, "mapper", fmt.Sprintf("mapper_%s.go", e.Identifier)),
			TemplateName: path.Join("graph", "graph_mapper_entity_json"),
			Data:         e,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
