package core

import (
	"context"
	"fmt"
	"path"
	"text/template"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

func generateCustomQueries(ctx context.Context, req coreSubModuleRequest) error {
	fmt.Printf("--[GPG] Generating core module custom queries: %s\n", req.Entity.Identifier)
	for _, cq := range req.Entity.CustomQueries {
		imports := map[string]any{}
		inputFields, allFields, mapping := GetCustomQueryFields(cq.Condition, req.Project)
		for _, f := range allFields {
			if f.Import != nil {
				imports[*f.Import] = struct{}{}
			}
		}

		fetchTemplate := fetchModuleTemplate{
			Package:           req.Entity.Identifier,
			ProjectIdentifier: req.Project.Identifier,
			ProjectModule:     req.Project.Module,
			EntityIdentifier:  req.Entity.Identifier,
			EntityName:        helpers.ToCamelCase(req.Entity.Identifier),
			CustomQuery:       cq,
			Imports:           helpers.MapKeys(imports),
		}

		err := generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, "types", fmt.Sprintf("fetch_custom_%s.go", helpers.ToSnakeCase(cq.Name))),
			TemplateName: path.Join("core", "core_module_fetch_custom_types"),
			Data:         fetchTemplate,
			Funcs: template.FuncMap{
				"UniqueFields": func(cq entity.CustomQuery) map[string]field.Template {
					return inputFields
				},
			},
		})
		if err != nil {
			return err
		}
		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(req.ModuleDir, req.Entity.Identifier, fmt.Sprintf("fetch_custom_%s.go", helpers.ToSnakeCase(cq.Name))),
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
		if err != nil {
			return err
		}
	}
	return nil
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
			fullField, entity, found := helpers.FieldFromProject(project, comp.FieldTwo.ParentIdentifier, comp.FieldTwo.Identifier)
			if found {
				allFields[comp.FieldTwo.Identifier] = field.ResolveFieldType(fullField, entity, nil)
				fullField.Identifier = comp.FieldOne.Identifier
				inputFields[comp.FieldOne.Identifier] = field.ResolveFieldType(fullField, entity, nil)
				mapping[comp.FieldTwo.Identifier] = inputFields[comp.FieldOne.Identifier]
			}
		}
		if comp.FieldTwo.InputField {
			fullField, entity, found := helpers.FieldFromProject(project, comp.FieldOne.ParentIdentifier, comp.FieldOne.Identifier)
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
