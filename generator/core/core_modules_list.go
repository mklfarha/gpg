package core

import (
	"context"
	"fmt"
	"path"

	"github.com/gertd/go-pluralize"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/helpers"
)

type listData struct {
	ProjectIdentifier string
	EntityIdentifier  string
	EntityName        string
	EntityNamePlural  string
}

func generateList(ctx context.Context, req coreSubModuleRequest) error {
	fmt.Printf("--[GPG] Generating core module list: %s\n", req.Entity.Identifier)
	pl := pluralize.NewClient()
	listData := listData{
		ProjectIdentifier: req.Project.Identifier,
		EntityIdentifier:  req.Entity.Identifier,
		EntityName:        helpers.ToCamelCase(req.Entity.Identifier),
		EntityNamePlural:  pl.Plural(helpers.ToCamelCase(req.Entity.Identifier)),
	}
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, "types", "list.go"),
		TemplateName: path.Join("core", "core_module_list_types"),
		Data:         listData,
	})
	if err != nil {
		return err
	}

	return generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, "list.go"),
		TemplateName: path.Join("core", "core_module_list"),
		Data:         listData,
	})

}
