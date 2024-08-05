package generator

import (
	"context"
	"fmt"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

type CustomUpsert struct {
	ProjectIdentifier string
	EntityName        string
	EntityIdentifier  string
	FuncName          string
	Field             field.Template
}

func GenerateCustom(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG] Generating Custom\n")
	projectDir := ProjectDir(ctx, rootPath, project)
	customDir := path.Join(projectDir, CUSTOM_DIR)

	GenerateFile(ctx, FileRequest{
		OutputFile:   path.Join(customDir, "generate_uuid.go"),
		TemplateName: path.Join("custom", "custom_generate_uuid"),
	})

	GenerateFile(ctx, FileRequest{
		OutputFile:   path.Join(customDir, "generate_encrypt.go"),
		TemplateName: path.Join("custom", "custom_generate_encrypt"),
	})

	GenerateFile(ctx, FileRequest{
		OutputFile:   path.Join(customDir, "generate_time.go"),
		TemplateName: path.Join("custom", "custom_generate_time"),
	})

	for _, e := range project.Entities {
		for _, f := range e.Fields {
			if f.Autogenerated.Type == entity.CustomAutogeneratedType {
				fileName := path.Join(customDir, fmt.Sprintf("%s_%s_%s.go", e.Identifier, f.Identifier, helpers.ToSnakeCase(f.Autogenerated.FuncName)))
				if !FileExists(fileName) {
					field := field.ResolveFieldType(f, e, nil)
					GenerateFile(ctx, FileRequest{
						OutputFile:   fileName,
						TemplateName: path.Join("custom", "custom_generate_upsert"),
						Data: CustomUpsert{
							ProjectIdentifier: project.Identifier,
							EntityName:        helpers.ToCamelCase(e.Identifier),
							EntityIdentifier:  e.Identifier,
							FuncName:          f.Autogenerated.FuncName,
							Field:             field,
						},
					})
				}
			}
		}
	}

	return nil
}
