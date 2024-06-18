package core

import (
	"context"
	"fmt"
	"path"

	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

type upsertModuleTemplate struct {
	Package          string
	EntityName       string
	EntityIdentifier string
	ProjectName      string
	PrimaryKey       field.Template
	Fields           []field.Template
	Imports          []string
}

func generateUpsert(ctx context.Context, req coreSubModuleRequest) error {
	fmt.Printf("--[GPG] Generating core module upsert: %s\n", req.Entity.Identifier)
	primaryKey := field.ResolveFieldType(helpers.EntityPrimaryKey(req.Entity), req.Entity, nil)
	upsertTemplate := upsertModuleTemplate{
		Package:          req.Entity.Identifier,
		ProjectName:      req.Project.Identifier,
		EntityIdentifier: req.Entity.Identifier,
		EntityName:       helpers.ToCamelCase(req.Entity.Identifier),
		PrimaryKey:       primaryKey,
		Fields:           req.Fields,
		Imports:          helpers.MapKeys(req.Imports),
	}

	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, "types", "upsert.go"),
		TemplateName: path.Join("core", "core_module_upsert_types"),
		Data:         upsertTemplate,
	})
	if err != nil {
		return err
	}
	return generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, "upsert.go"),
		TemplateName: path.Join("core", "core_module_upsert"),
		Data:         upsertTemplate,
	})
}
