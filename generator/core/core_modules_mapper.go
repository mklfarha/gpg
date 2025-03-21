package core

import (
	"context"
	"fmt"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

type mapperModuleTemplate struct {
	Package           string
	EntityName        string
	ProjectIdentifier string
	ProjectModule     string
	Fields            []field.Template
	Imports           []string
	HasArrayField     bool
	HasNullString     bool
	HasNullUUID       bool
}

func generateMapper(ctx context.Context, req coreSubModuleRequest) error {
	fmt.Printf("--[GPG] Generating core module mapper: %s\n", req.Entity.Identifier)
	hasArrayField := false
	hasNullString := false
	hasNullUUID := false
	for _, f := range req.Fields {
		if f.InternalType == entity.ArrayFieldType {
			hasArrayField = true
		}

		if f.Type == "*string" {
			hasNullString = true
		}

		if !f.Required && f.InternalType == entity.UUIDFieldType {
			hasNullUUID = true
		}
	}
	mapperTemplate := mapperModuleTemplate{
		Package:           req.Entity.Identifier,
		ProjectIdentifier: req.Project.Identifier,
		ProjectModule:     req.Project.Module,
		EntityName:        helpers.ToCamelCase(req.Entity.Identifier),
		Fields:            req.Fields,
		Imports:           helpers.MapKeys(req.Imports),
		HasArrayField:     hasArrayField,
		HasNullString:     hasNullString,
		HasNullUUID:       hasNullUUID,
	}

	return generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, "mapper.go"),
		TemplateName: path.Join("core", "core_module_mapper"),
		Data:         mapperTemplate,
	})
}
