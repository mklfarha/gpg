package proto

import (
	"context"
	"fmt"
	"path"
	"text/template"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/core"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

type ProtoEntityTemplate struct {
	ProjectIdentifier string
	Identifier        string
	EntityName        string
	Fields            []field.Template
	Selects           []core.RepoSchemaSelectStatement
	CustomQueries     []entity.CustomQuery
	Search            bool
	Enums             []ProtoEnumTemplate
}

type ProtoEnumTemplate struct {
	Field   field.Template
	Options []string
}

func Generate(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG][Proto] Generating Directory\n")
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	protoDir := path.Join(projectDir, generator.PROTO_DIR)

	pl := pluralize.NewClient()

	//generate model
	for _, e := range project.Entities {
		nested, err := handleEntity(ctx, protoDir, project, e)
		if err != nil {
			return err
		}
		for _, n := range nested {
			_, err := handleEntity(ctx, protoDir, project, entity.Entity{
				Identifier: pl.Singular(n.Identifier),
				Fields:     n.JSONConfig.Fields,
			})
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func handleEntity(ctx context.Context, protoDir string, project entity.Project, e entity.Entity) ([]entity.Field, error) {
	fields := []field.Template{}
	searchable := false
	enums := []ProtoEnumTemplate{}
	nested := []entity.Field{}
	var err error
	if len(e.Fields) > 0 {
		for _, f := range e.Fields {
			fieldTemplate := field.ResolveFieldType(f, e, nil)
			fields = append(fields, fieldTemplate)
			for _, field := range fields {
				if field.Enum {
					enums = append(enums, ProtoEnumTemplate{
						Field:   field,
						Options: field.ProtoEnumOptions,
					})
				}
				if f.Type == entity.JSONFieldType {
					nested = append(nested, f)
				}
			}

			if f.StorageConfig.Search {
				searchable = true
			}
		}
		entityTemplate := ProtoEntityTemplate{
			ProjectIdentifier: project.Identifier,
			Identifier:        e.Identifier,
			EntityName:        helpers.ToCamelCase(e.Identifier),
			Fields:            fields,
			Search:            searchable,
			CustomQueries:     e.CustomQueries,
			Enums:             enums,
		}
		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:      path.Join(protoDir, "proto", fmt.Sprintf("%s.proto", e.Identifier)),
			TemplateName:    path.Join("proto", "model"),
			Data:            entityTemplate,
			DisableGoFormat: true,
			Funcs: template.FuncMap{
				"Inc": helpers.Inc,
			},
		})
	}
	return nested, err
}
