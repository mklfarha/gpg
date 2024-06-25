package core

import (
	"context"
	"fmt"
	"os/exec"
	"path"

	"github.com/gertd/go-pluralize"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

type EntityTemplate struct {
	ProjectName string
	Package     string
	Imports     []string
	EntityName  string
	Identifier  string
	Fields      []field.Template
	JSON        bool
	JSONField   field.Template
}

type EnumTemplate struct {
	Package       string
	EnumName      string
	EnumNameUpper string
	Values        []string
	Options       []entity.OptionValue
}

func GenerateCoreEntities(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG] Generating core entities\n")
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	entitiesDir := path.Join(projectDir, generator.CORE_ENTITY_DIR)
	allimports := map[string]any{}
	for _, e := range project.Entities {
		fmt.Printf("----[GPG] Generating entity: %s\n", e.Identifier)
		entityDir := path.Join(entitiesDir, e.Identifier)
		entityTemplate, entityImports := resolveEntityTemplate(e, project)
		for imp, _ := range entityImports {
			allimports[imp] = struct{}{}
		}
		generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(entityDir, fmt.Sprintf("%s.go", e.Identifier)),
			TemplateName: path.Join("core", "entity"),
			Data:         entityTemplate,
		})
		generateEnums(ctx, entityDir, e)
		generateJSONEntities(ctx, entityDir, e, project)
	}

	for imp, _ := range allimports {
		cmd := exec.Command("go", "get", imp)
		cmd.Dir = projectDir
		err := cmd.Run()
		if err != nil {
			fmt.Printf("error running go get %s\n", imp)
		}
	}

	entityTypesDir := path.Join(entitiesDir, "types")
	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(entityTypesDir, "field_types.go"),
		TemplateName: path.Join("core", "field_types"),
	})

	return nil
}

func resolveEntityTemplate(e entity.Entity, project entity.Project) (EntityTemplate, map[string]any) {
	fields, imports := field.ResolveFieldsAndImports(e.Fields, e, nil)
	return EntityTemplate{
		ProjectName: project.Identifier,
		Package:     e.Identifier,
		EntityName:  helpers.ToCamelCase(e.Identifier),
		Identifier:  e.Identifier,
		Fields:      fields,
		Imports:     helpers.MapKeys(imports),
	}, imports
}

func generateJSONEntities(ctx context.Context, entityDir string, e entity.Entity, project entity.Project) {
	pl := pluralize.NewClient()
	for _, f := range e.Fields {
		if f.Type == entity.JSONFieldType {
			singularIdentifier := pl.Singular(f.Identifier)
			fmt.Printf("----[GPG] Generating json entity: %s\n", singularIdentifier)
			fields, imports := field.ResolveFieldsAndImports(f.JSONConfig.Fields, e, &singularIdentifier)
			field := field.ResolveFieldType(f, e, nil)
			entityTemplate := EntityTemplate{
				ProjectName: project.Identifier,
				Package:     e.Identifier,
				EntityName:  helpers.ToCamelCase(singularIdentifier),
				Fields:      fields,
				Imports:     helpers.MapKeys(imports),
				JSON:        true,
				JSONField:   field,
			}
			generator.GenerateFile(ctx, generator.FileRequest{
				OutputFile:   path.Join(entityDir, fmt.Sprintf("%s.go", singularIdentifier)),
				TemplateName: path.Join("core", "entity"),
				Data:         entityTemplate,
			})
		}
	}
}
