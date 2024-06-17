package proto

import (
	"context"
	"fmt"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/helpers"
)

func generateServer(ctx context.Context, protoDir string, project entity.Project, standaloneEntities []ProtoEntityTemplate) error {
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(protoDir, "server", "server.go"),
		TemplateName: path.Join("proto", "server"),
		Data: ProtoServiceTemplate{
			Identifier: project.Identifier,
			Name:       helpers.ToCamelCase(project.Identifier),
		},
		DisableGoFormat: false,
	})
	if err != nil {
		return err
	}

	for _, se := range standaloneEntities {
		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:      path.Join(protoDir, "server", fmt.Sprintf("create_%s.go", se.Identifier)),
			TemplateName:    path.Join("proto", "server_create_entity"),
			Data:            se,
			DisableGoFormat: false,
		})
		if err != nil {
			return err
		}
		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:      path.Join(protoDir, "server", fmt.Sprintf("update_%s.go", se.Identifier)),
			TemplateName:    path.Join("proto", "server_update_entity"),
			Data:            se,
			DisableGoFormat: false,
		})
		if err != nil {
			return err
		}
	}

	return err
}
