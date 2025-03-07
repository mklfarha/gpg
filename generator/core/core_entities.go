package core

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

type EntityTemplate struct {
	ProjectIdentifier    string
	ProjectModule        string
	Package              string
	Imports              []string
	EntityName           string
	Identifier           string
	Fields               []field.Template
	JSON                 bool
	JSONField            field.Template
	PrimaryKeyIdentifier string
	PrimaryKeyName       string
	UsesRandomValues     bool
}

type EnumTemplate struct {
	ProjectIdentifier string
	ProjectModule     string
	Package           string
	EnumName          string
	EnumNameUpper     string
	Values            []string
	Options           []entity.OptionValue
}

func GenerateCoreEntities(ctx context.Context, rootPath string, project entity.Project) error {

	fmt.Printf("--[GPG] Generating core entities\n")
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	entitiesDir := path.Join(projectDir, generator.CORE_ENTITY_DIR)

	err := os.RemoveAll(entitiesDir)
	if err != nil {
		fmt.Printf("ERROR: Deleting entity directory\n")
	}

	allimports := map[string]any{}
	for _, e := range project.Entities {
		generateEnums(ctx, project, entitiesDir, e)
		generateJSONEntities(ctx, entitiesDir, e, project)

		fmt.Printf("----[GPG] Generating entity: %s\n", e.Identifier)
		entityDir := path.Join(entitiesDir, e.Identifier)
		entityTemplate, entityImports := resolveEntityTemplate(e, project)
		for imp := range entityImports {
			allimports[imp] = struct{}{}
		}
		generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(entityDir, fmt.Sprintf("%s.go", e.Identifier)),
			TemplateName: path.Join("core", "entity"),
			Data:         entityTemplate,
		})

	}

	for imp := range allimports {
		if !strings.Contains(imp, fmt.Sprintf("%s/", project.Identifier)) {
			cmd := exec.Command("go", "get", imp)
			cmd.Dir = projectDir
			err := cmd.Run()
			if err != nil {
				fmt.Printf("error running go get %s\n", imp)
			}
		}
	}

	entityTypesDir := path.Join(entitiesDir, "types")
	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(entityTypesDir, "types.go"),
		TemplateName: path.Join("core", "entity_types"),
	})

	entityMapperDir := path.Join(projectDir, generator.CORE_ENTITY_DIR, "mapper")
	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(entityMapperDir, "mapper.go"),
		TemplateName: path.Join("core", "entity_mapper"),
		Data:         project,
	})

	return nil
}

func resolveEntityTemplate(e entity.Entity, project entity.Project) (EntityTemplate, map[string]any) {
	fields, imports := field.ResolveFieldsAndImports(project, e.Fields, e, nil)
	primaryKey := helpers.EntityPrimaryKey(e)

	tpl := EntityTemplate{
		ProjectIdentifier:    project.Identifier,
		ProjectModule:        project.Module,
		Package:              e.Identifier,
		EntityName:           helpers.ToCamelCase(e.Identifier),
		Identifier:           e.Identifier,
		Fields:               fields,
		Imports:              helpers.MapKeys(imports),
		PrimaryKeyIdentifier: primaryKey.Identifier,
		PrimaryKeyName:       helpers.ToCamelCase(primaryKey.Identifier),
		UsesRandomValues:     e.UsesRandomValues(),
	}

	return tpl, imports
}

func generateJSONEntities(ctx context.Context, entitiesDir string, e entity.Entity, project entity.Project) {
	for _, f := range e.Fields {
		if f.Type == entity.JSONFieldType && !f.JSONConfig.Reuse && len(f.JSONConfig.Fields) > 0 {

			ft := field.ResolveFieldType(f, e, &field.Template{
				Identifier: f.Identifier,
			})

			fmt.Printf("----[GPG] Generating json entity: %s\n", f.JSONConfig.Identifier)
			entityDir := path.Join(entitiesDir, f.JSONConfig.Identifier)

			jsonEntity := entity.Entity{
				Identifier: f.JSONConfig.Identifier,
				Fields:     f.JSONConfig.Fields,
			}

			fields, imports := field.ResolveFieldsAndImports(project, f.JSONConfig.Fields, e, &ft)
			entityTemplate := EntityTemplate{
				ProjectIdentifier: project.Identifier,
				ProjectModule:     project.Module,
				Package:           f.JSONConfig.Identifier,
				Identifier:        f.JSONConfig.Identifier,
				EntityName:        ft.Type,
				Fields:            fields,
				Imports:           helpers.MapKeys(imports),
				JSON:              true,
				JSONField:         ft,
				UsesRandomValues:  jsonEntity.UsesRandomValues(),
			}

			generator.GenerateFile(ctx, generator.FileRequest{
				OutputFile:   path.Join(entityDir, fmt.Sprintf("%s.go", f.JSONConfig.Identifier)),
				TemplateName: path.Join("core", "entity"),
				Data:         entityTemplate,
			})

			if f.HasNestedJsonFields() {
				generateJSONEntities(ctx, entitiesDir, entity.Entity{
					Identifier: f.JSONConfig.Identifier,
					Fields:     f.JSONConfig.Fields,
				}, project)
			}
		}
	}
}
