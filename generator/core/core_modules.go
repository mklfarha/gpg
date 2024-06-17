package core

import (
	"context"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

type CoreModuleTemplate struct {
	Package          string
	ProjectName      string
	EntityIdentifier string
	EntityName       string
	SelectStatements []RepoSchemaSelectStatement
	CustomQueries    []entity.CustomQuery
	SearchFields     []field.Template
}

type FetchModuleTemplate struct {
	Package          string
	EntityName       string
	EntityIdentifier string
	ProjectName      string
	Select           RepoSchemaSelectStatement
	CustomQuery      entity.CustomQuery
	Fields           []field.Template
	Imports          []string
	SearchFields     []field.Template
}

type MapperModuleTemplate struct {
	Package     string
	EntityName  string
	ProjectName string
	Fields      []field.Template
	Imports     []string
}

type UpsertModuleTemplate struct {
	Package          string
	EntityName       string
	EntityIdentifier string
	ProjectName      string
	PrimaryKey       field.Template
	Fields           []field.Template
	Imports          []string
}

func GenerateCoreModules(ctx context.Context, rootPath string, project entity.Project) {
	fmt.Printf("--[GPG] Generating core modules\n")
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	moduleDir := path.Join(projectDir, generator.CORE_MODULE_DIR)

	err := os.RemoveAll(moduleDir)
	if err != nil {
		fmt.Printf("ERROR: Deleting module directory\n")
	}

	for _, e := range project.Entities {
		fmt.Printf("--[GPG] Generating core module: %s\n", e.Identifier)
		selects := ResolveSelectStatements(project, e)
		searchFields := ResolveSearchFields(e)
		moduleTemplate := CoreModuleTemplate{
			Package:          e.Identifier,
			ProjectName:      project.Identifier,
			EntityIdentifier: e.Identifier,
			EntityName:       helpers.ToCamelCase(e.Identifier),
			SelectStatements: selects,
			CustomQueries:    e.CustomQueries,
			SearchFields:     searchFields,
		}
		generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(moduleDir, e.Identifier, fmt.Sprintf("%s.go", e.Identifier)),
			TemplateName: path.Join("core", "core_module"),
			Data:         moduleTemplate,
		})
		generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(moduleDir, e.Identifier, "option.go"),
			TemplateName: path.Join("core", "core_module_options"),
			Data:         moduleTemplate,
		})

		fields, imports := field.ResolveFieldsAndImports(e.Fields, e, nil)
		mapperTemplate := MapperModuleTemplate{
			Package:     e.Identifier,
			ProjectName: project.Identifier,
			EntityName:  helpers.ToCamelCase(e.Identifier),
			Fields:      fields,
			Imports:     helpers.MapKeys(imports),
		}
		generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(moduleDir, e.Identifier, "mapper.go"),
			TemplateName: path.Join("core", "core_module_mapper"),
			Data:         mapperTemplate,
		})

		// selects
		for _, sel := range selects {
			imports := map[string]any{}
			for _, f := range sel.Fields {
				if f.Field.Import != nil {
					imports[*f.Field.Import] = struct{}{}
				}
			}
			fetchTemplate := FetchModuleTemplate{
				Package:          e.Identifier,
				ProjectName:      project.Identifier,
				EntityIdentifier: e.Identifier,
				EntityName:       helpers.ToCamelCase(e.Identifier),
				Select:           sel,
				Imports:          helpers.MapKeys(imports),
			}
			generator.GenerateFile(ctx, generator.FileRequest{
				OutputFile:   path.Join(moduleDir, e.Identifier, "types", fmt.Sprintf("fetch_%s.go", helpers.ToSnakeCase(sel.Name))),
				TemplateName: path.Join("core", "core_module_fetch_types"),
				Data:         fetchTemplate,
			})
			generator.GenerateFile(ctx, generator.FileRequest{
				OutputFile:   path.Join(moduleDir, e.Identifier, fmt.Sprintf("fetch_%s.go", helpers.ToSnakeCase(sel.Name))),
				TemplateName: path.Join("core", "core_module_fetch"),
				Data:         fetchTemplate,
			})
		}

		// custom queries
		for _, cq := range e.CustomQueries {
			imports := map[string]any{}
			inputFields, allFields, mapping := GetCustomQueryFields(cq.Condition, project)
			for _, f := range allFields {
				if f.Import != nil {
					imports[*f.Import] = struct{}{}
				}
			}

			fetchTemplate := FetchModuleTemplate{
				Package:          e.Identifier,
				ProjectName:      project.Identifier,
				EntityIdentifier: e.Identifier,
				EntityName:       helpers.ToCamelCase(e.Identifier),
				CustomQuery:      cq,
				Imports:          helpers.MapKeys(imports),
			}

			generator.GenerateFile(ctx, generator.FileRequest{
				OutputFile:   path.Join(moduleDir, e.Identifier, "types", fmt.Sprintf("fetch_custom_%s.go", helpers.ToSnakeCase(cq.Name))),
				TemplateName: path.Join("core", "core_module_fetch_custom_types"),
				Data:         fetchTemplate,
				Funcs: template.FuncMap{
					"UniqueFields": func(cq entity.CustomQuery) map[string]field.Template {
						return inputFields
					},
				},
			})
			generator.GenerateFile(ctx, generator.FileRequest{
				OutputFile:   path.Join(moduleDir, e.Identifier, fmt.Sprintf("fetch_custom_%s.go", helpers.ToSnakeCase(cq.Name))),
				TemplateName: path.Join("core", "core_module_fetch_custom"),
				Data:         fetchTemplate,
				Funcs: template.FuncMap{
					"Fields": func(cq entity.CustomQuery) map[string]field.Template {
						return allFields
					},
					"MapToInput": func(in field.Template) field.Template {
						return mapping[in.Identifier]
					},
				},
			})
		}

		// search
		if len(searchFields) > 0 {
			searchTemplate := FetchModuleTemplate{
				Package:          e.Identifier,
				ProjectName:      project.Identifier,
				EntityIdentifier: e.Identifier,
				EntityName:       helpers.ToCamelCase(e.Identifier),
				Imports:          helpers.MapKeys(imports),
				SearchFields:     searchFields,
			}

			generator.GenerateFile(ctx, generator.FileRequest{
				OutputFile:   path.Join(moduleDir, e.Identifier, "types", fmt.Sprintf("search_%s.go", e.Identifier)),
				TemplateName: path.Join("core", "core_module_search_types"),
				Data:         searchTemplate,
			})
			generator.GenerateFile(ctx, generator.FileRequest{
				OutputFile:   path.Join(moduleDir, e.Identifier, fmt.Sprintf("search_%s.go", e.Identifier)),
				TemplateName: path.Join("core", "core_module_search"),
				Data:         searchTemplate,
			})

		}

		// upsert
		primaryKey := field.ResolveFieldType(EntityPrimaryKey(e), e, nil)
		upsertTemplate := UpsertModuleTemplate{
			Package:          e.Identifier,
			ProjectName:      project.Identifier,
			EntityIdentifier: e.Identifier,
			EntityName:       helpers.ToCamelCase(e.Identifier),
			PrimaryKey:       primaryKey,
			Fields:           fields,
			Imports:          helpers.MapKeys(imports),
		}

		generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(moduleDir, e.Identifier, "types", "upsert.go"),
			TemplateName: path.Join("core", "core_module_upsert_types"),
			Data:         upsertTemplate,
		})
		generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(moduleDir, e.Identifier, "upsert.go"),
			TemplateName: path.Join("core", "core_module_upsert"),
			Data:         upsertTemplate,
		})

	}

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(projectDir, generator.CORE_DIR, "core.go"),
		TemplateName: path.Join("core", "core_main"),
		Data:         project,
		Funcs: template.FuncMap{
			"ToCamelCase": helpers.ToCamelCase,
		},
	})
}

func GetCustomQueryFields(cond entity.QueryCondition, project entity.Project) (inputFields map[string]field.Template, allFields map[string]field.Template, mapping map[string]field.Template) {
	inputFields = map[string]field.Template{}
	allFields = map[string]field.Template{}
	mapping = map[string]field.Template{}
	GetCustomQueryFieldsRecursive(cond, project, inputFields, allFields, mapping)
	return
}

func GetCustomQueryFieldsRecursive(cond entity.QueryCondition, project entity.Project, inputFields map[string]field.Template, allFields map[string]field.Template, mapping map[string]field.Template) {
	for _, comp := range cond.Comparisons {
		if comp.FieldOne.InputField {
			fullField, entity, found := FieldFromProject(project, comp.FieldTwo.ParentIdentifier, comp.FieldTwo.Identifier)
			if found {
				allFields[comp.FieldTwo.Identifier] = field.ResolveFieldType(fullField, entity, nil)
				fullField.Identifier = comp.FieldOne.Identifier
				inputFields[comp.FieldOne.Identifier] = field.ResolveFieldType(fullField, entity, nil)
				mapping[comp.FieldTwo.Identifier] = inputFields[comp.FieldOne.Identifier]
			}
		}
		if comp.FieldTwo.InputField {
			fullField, entity, found := FieldFromProject(project, comp.FieldOne.ParentIdentifier, comp.FieldOne.Identifier)
			if found {
				allFields[fullField.Identifier] = field.ResolveFieldType(fullField, entity, nil)
				fullField.Identifier = comp.FieldTwo.Identifier
				inputFields[comp.FieldTwo.Identifier] = field.ResolveFieldType(fullField, entity, nil)
				mapping[comp.FieldOne.Identifier] = inputFields[comp.FieldTwo.Identifier]
			}
		}
	}

	for _, c := range cond.Conditions {
		GetCustomQueryFieldsRecursive(c, project, inputFields, allFields, mapping)
	}
}

func FieldFromProject(project entity.Project, entityIdentifier, fieldIdentifier string) (entity.Field, entity.Entity, bool) {
	for _, e := range project.Entities {
		if e.Identifier == entityIdentifier {
			for _, f := range e.Fields {
				if f.Identifier == fieldIdentifier {
					return f, e, true
				}
			}
		}
	}
	return entity.Field{}, entity.Entity{}, false
}
