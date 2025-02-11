package core

import (
	"context"
	"fmt"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/core/events"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

type upsertModuleTemplate struct {
	Package             string
	EntityName          string
	EntityIdentifier    string
	ProjectIdentifier   string
	ProjectModule       string
	PrimaryKey          field.Template
	Fields              []field.Template
	Imports             []string
	HasVersionField     bool
	VersionField        field.Template
	ShouldPublishEvents bool
	HasArrayField       bool
}

func generateUpsert(ctx context.Context, req coreSubModuleRequest) error {
	fmt.Printf("--[GPG] Generating core module upsert: %s\n", req.Entity.Identifier)
	primaryKey := field.ResolveFieldType(helpers.EntityPrimaryKey(req.Entity), req.Entity, nil)
	hasArrayField := false
	for _, f := range req.Fields {
		if f.InternalType == entity.ArrayFieldType {
			hasArrayField = true
		}
	}
	upsertTemplate := upsertModuleTemplate{
		Package:             req.Entity.Identifier,
		ProjectIdentifier:   req.Project.Identifier,
		ProjectModule:       req.Project.Module,
		EntityIdentifier:    req.Entity.Identifier,
		EntityName:          helpers.ToCamelCase(req.Entity.Identifier),
		PrimaryKey:          primaryKey,
		Fields:              req.Fields,
		Imports:             helpers.MapKeys(req.Imports),
		ShouldPublishEvents: events.ShouldPublishEvents(req.Project, req.Entity.Identifier),
		HasArrayField:       hasArrayField,
	}

	versionField := VersionField(req.Fields)
	if versionField != nil {
		upsertTemplate.HasVersionField = true
		upsertTemplate.VersionField = *versionField
	}

	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, "types", "upsert.go"),
		TemplateName: path.Join("core", "core_module_upsert_types"),
		Data:         upsertTemplate,
	})
	if err != nil {
		return err
	}
	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, "upsert_insert.go"),
		TemplateName: path.Join("core", "core_module_upsert_insert"),
		Data:         upsertTemplate,
	})
	if err != nil {
		return err
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, "upsert_update.go"),
		TemplateName: path.Join("core", "core_module_upsert_update"),
		Data:         upsertTemplate,
	})
	if err != nil {
		return err
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, "upsert.go"),
		TemplateName: path.Join("core", "core_module_upsert"),
		Data:         upsertTemplate,
	})
	if err != nil {
		return err
	}
	return nil
}

func VersionField(fields []field.Template) *field.Template {
	for _, f := range fields {
		if f.Identifier == "version" && f.InternalType == entity.IntFieldType {
			return &f
		}
	}
	return nil
}
