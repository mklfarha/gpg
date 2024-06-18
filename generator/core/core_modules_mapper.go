package core

import (
	"context"
	"fmt"
	"path"

	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

type mapperModuleTemplate struct {
	Package     string
	EntityName  string
	ProjectName string
	Fields      []field.Template
	Imports     []string
}

func generateMapper(ctx context.Context, req coreSubModuleRequest) error {
	fmt.Printf("--[GPG] Generating core module mapper: %s\n", req.Entity.Identifier)
	mapperTemplate := mapperModuleTemplate{
		Package:     req.Entity.Identifier,
		ProjectName: req.Project.Identifier,
		EntityName:  helpers.ToCamelCase(req.Entity.Identifier),
		Fields:      req.Fields,
		Imports:     helpers.MapKeys(req.Imports),
	}
	return generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, "mapper.go"),
		TemplateName: path.Join("core", "core_module_mapper"),
		Data:         mapperTemplate,
	})
}
