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
	return generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(graphDir, "mapper", "mapper.go"),
		TemplateName: path.Join("graph", "graph_mapper"),
		Data: GraphQueriesTemplate{
			ProjectName:  project.Identifier,
			Entities:     res.EntityTemplates,
			JSONEntities: res.JsonEntityTemplates,
			Enums:        res.EnumTemplates,
		},
	})
}
