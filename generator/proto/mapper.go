package proto

import (
	"context"
	"fmt"
	"path"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/helpers"
)

func generateEntityMapper(ctx context.Context, dir string, et ProtoEntityTemplate) error {
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(dir, fmt.Sprintf("%s.go", strcase.ToSnake(et.FinalIdentifier))),
		TemplateName:    path.Join("proto", "model_mapper"),
		Data:            et,
		DisableGoFormat: false,
		Funcs: template.FuncMap{
			"Inc": helpers.Inc,
		},
	})

	if err != nil {
		return err
	}
	return nil
}

func generateMappers(ctx context.Context, protoDir string, _ entity.Project, standaloneEntities []ProtoEntityTemplate, dependantEntities map[string][]ProtoEntityTemplate) error {
	fmt.Printf("--[GPG][Proto] Generating mappers\n")
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(protoDir, "mapper", "json.go"),
		TemplateName:    path.Join("proto", "mapper_json"),
		DisableGoFormat: false,
	})
	if err != nil {
		return err
	}

	for _, et := range standaloneEntities {
		dir := path.Join(protoDir, "mapper", et.FinalIdentifier)
		err := generateEntityMapper(ctx, dir, et)
		if err != nil {
			return err
		}

		nestedTemplates := dependantEntities[et.FinalIdentifier]
		for _, net := range nestedTemplates {
			err := generateEntityMapper(ctx, dir, net)
			if err != nil {
				return err
			}
		}
	}
	return err
}
