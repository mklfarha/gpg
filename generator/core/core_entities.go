package core

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

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
	ProjectName   string
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

	err := os.RemoveAll(entitiesDir)
	if err != nil {
		fmt.Printf("ERROR: Deleting module directory\n")
	}

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
		generateEnums(ctx, project, entityDir, e)
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
		OutputFile:   path.Join(entityTypesDir, "types.go"),
		TemplateName: path.Join("core", "types"),
	})

	entityRandomValuesDir := path.Join(entitiesDir, "randomvalues")
	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(entityRandomValuesDir, "randomvalues.go"),
		TemplateName: path.Join("core", "random_values"),
		Data: struct {
			ProjectName string
		}{
			ProjectName: project.Identifier,
		},
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
	for _, f := range e.Fields {
		if f.Type == entity.JSONFieldType {
			ft := field.ResolveFieldType(f, e, &field.Template{
				Identifier: f.Identifier,
			})
			fmt.Printf("----[GPG] Generating json entity: %s\n", ft.SingularIdentifier)
			fields, imports := field.ResolveFieldsAndImports(f.JSONConfig.Fields, e, &ft)
			entityTemplate := EntityTemplate{
				ProjectName: project.Identifier,
				Package:     e.Identifier,
				EntityName:  helpers.ToCamelCase(ft.SingularIdentifier),
				Fields:      fields,
				Imports:     helpers.MapKeys(imports),
				JSON:        true,
				JSONField:   ft,
			}
			generator.GenerateFile(ctx, generator.FileRequest{
				OutputFile:   path.Join(entityDir, fmt.Sprintf("%s.go", ft.SingularIdentifier)),
				TemplateName: path.Join("core", "entity"),
				Data:         entityTemplate,
			})
		}
	}
}
