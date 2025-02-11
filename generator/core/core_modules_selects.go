package core

import (
	"context"
	"fmt"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/core/repo"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

type fetchModuleTemplate struct {
	Package           string
	EntityName        string
	EntityIdentifier  string
	ProjectIdentifier string
	ProjectModule     string
	Select            repo.SchemaSelectStatement
	CustomQuery       entity.CustomQuery
	Fields            []field.Template
	Imports           []string
	SearchFields      []field.Template
}

func generateSelects(ctx context.Context, req coreSubModuleRequest) error {
	fmt.Printf("--[GPG] Generating core module selects: %s\n", req.Entity.Identifier)
	for _, sel := range req.Selects {
		imports := map[string]any{}
		for _, f := range sel.Fields {
			if f.Field.Import != nil {
				imports[*f.Field.Import] = struct{}{}
			}
		}
		fetchTemplate := fetchModuleTemplate{
			Package:           req.Entity.Identifier,
			ProjectIdentifier: req.Project.Identifier,
			ProjectModule:     req.Project.Module,
			EntityIdentifier:  req.Entity.Identifier,
			EntityName:        helpers.ToCamelCase(req.Entity.Identifier),
			Select:            sel,
			Imports:           helpers.MapKeys(imports),
		}
		err := generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, "types", fmt.Sprintf("fetch_%s.go", helpers.ToSnakeCase(sel.Name))),
			TemplateName: path.Join("core", "core_module_fetch_types"),
			Data:         fetchTemplate,
		})
		if err != nil {
			return err
		}
		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, fmt.Sprintf("fetch_%s.go", helpers.ToSnakeCase(sel.Name))),
			TemplateName: path.Join("core", "core_module_fetch"),
			Data:         fetchTemplate,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
