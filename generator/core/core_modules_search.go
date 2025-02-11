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

func generateSearch(ctx context.Context, req coreSubModuleRequest) error {
	fmt.Printf("--[GPG] Generating core module search: %s\n", req.Entity.Identifier)
	if len(req.SearchFields) > 0 {
		searchTemplate := fetchModuleTemplate{
			Package:           req.Entity.Identifier,
			ProjectIdentifier: req.Project.Identifier,
			ProjectModule:     req.Project.Module,
			EntityIdentifier:  req.Entity.Identifier,
			EntityName:        helpers.ToCamelCase(req.Entity.Identifier),
			Imports:           helpers.MapKeys(req.Imports),
			SearchFields:      req.SearchFields,
		}

		err := generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, "types", fmt.Sprintf("search_%s.go", req.Entity.Identifier)),
			TemplateName: path.Join("core", "core_module_search_types"),
			Data:         searchTemplate,
		})
		if err != nil {
			return err
		}
		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, fmt.Sprintf("search_%s.go", req.Entity.Identifier)),
			TemplateName: path.Join("core", "core_module_search"),
			Data:         searchTemplate,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func GetSearchFields(e entity.Entity) []field.Template {
	fields := []field.Template{}
	for _, f := range e.Fields {
		if f.StorageConfig.Search {
			fields = append(fields, field.ResolveFieldType(f, e, nil))
		}
	}
	return fields
}
