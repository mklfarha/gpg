package core

import (
	"context"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/core/repo"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

type coreSubModuleRequest struct {
	Project      entity.Project
	Entity       entity.Entity
	ModuleDir    string
	Fields       []field.Template
	Imports      map[string]any
	Selects      []repo.SchemaSelectStatement
	SearchFields []field.Template
}

func GenerateCoreModules(ctx context.Context, rootPath string, project entity.Project) error {
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	moduleDir := path.Join(projectDir, generator.CORE_MODULE_DIR)

	err := os.RemoveAll(moduleDir)
	if err != nil {
		fmt.Printf("ERROR: Deleting module directory\n")
	}

	fmt.Printf("--[GPG] Generating core modules\n")
	for _, e := range project.Entities {
		selects := repo.ResolveSelectStatements(project, e)
		searchFields := GetSearchFields(e)
		fields, imports := field.ResolveFieldsAndImports(e.Fields, e, nil)

		req := coreSubModuleRequest{
			Project:      project,
			Entity:       e,
			ModuleDir:    moduleDir,
			Fields:       fields,
			Imports:      imports,
			Selects:      selects,
			SearchFields: searchFields,
		}

		// generate base files for entities, module and options
		err = generateBaseCoreModule(ctx, req)
		if err != nil {
			return err
		}

		//generate mappers
		err = generateMapper(ctx, req)
		if err != nil {
			return err
		}

		// generate selects
		err = generateSelects(ctx, req)
		if err != nil {
			return err
		}

		// custom queries
		err = generateCustomQueries(ctx, req)
		if err != nil {
			return err
		}

		// search
		err = generateSearch(ctx, req)
		if err != nil {
			return err
		}

		// upsert
		err = generateUpsert(ctx, req)
		if err != nil {
			return err
		}

		// list
		err = generateList(ctx, req)
		if err != nil {
			return err
		}
	}

	// generate module types
	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(projectDir, generator.CORE_DIR, "types", "types.go"),
		TemplateName: path.Join("core", "core_module_types"),
		Data:         project,
		Funcs: template.FuncMap{
			"ToCamelCase": helpers.ToCamelCase,
		},
	})
	if err != nil {
		return err
	}

	// generate main module
	return generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(projectDir, generator.CORE_DIR, "core.go"),
		TemplateName: path.Join("core", "core_main"),
		Data:         project,
		Funcs: template.FuncMap{
			"ToCamelCase": helpers.ToCamelCase,
		},
	})
}
