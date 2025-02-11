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

type coreModuleTemplate struct {
	Project           entity.Project
	Package           string
	ProjectIdentifier string
	ProjectModule     string
	EntityIdentifier  string
	EntityName        string
	SelectStatements  []repo.SchemaSelectStatement
	CustomQueries     []entity.CustomQuery
	SearchFields      []field.Template
}

func generateBaseCoreModule(ctx context.Context, req coreSubModuleRequest) error {
	fmt.Printf("--[GPG] Generating core module base: %s\n", req.Entity.Identifier)
	moduleTemplate := coreModuleTemplate{
		Project:           req.Project,
		Package:           req.Entity.Identifier,
		ProjectIdentifier: req.Project.Identifier,
		ProjectModule:     req.Project.Module,
		EntityIdentifier:  req.Entity.Identifier,
		EntityName:        helpers.ToCamelCase(req.Entity.Identifier),
		SelectStatements:  req.Selects,
		CustomQueries:     req.Entity.CustomQueries,
		SearchFields:      req.SearchFields,
	}
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, fmt.Sprintf("%s.go", req.Entity.Identifier)),
		TemplateName: path.Join("core", "core_module"),
		Data:         moduleTemplate,
	})
	if err != nil {
		return err
	}
	return generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, "option.go"),
		TemplateName: path.Join("core", "core_module_options"),
		Data:         moduleTemplate,
	})
}
