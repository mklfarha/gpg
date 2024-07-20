package core

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

func generateEnums(ctx context.Context, project entity.Project, entitiesDir string, e entity.Entity) {
	for _, f := range e.Fields {
		if f.Type == entity.OptionsSingleFieldType || f.Type == entity.OptionsManyFieldType {
			entityDir := path.Join(entitiesDir, e.Identifier)
			generateEnum(ctx, project, entityDir, e, f, nil, e.Identifier)
		}
		if f.Type == entity.JSONFieldType {
			for _, jf := range f.JSONConfig.Fields {
				if jf.Type == entity.OptionsSingleFieldType || jf.Type == entity.OptionsManyFieldType {
					ft := field.ResolveFieldType(f, e, nil)
					entityDir := path.Join(entitiesDir, f.JSONConfig.Identifier)
					generateEnum(ctx, project, entityDir, e, jf, &ft, f.JSONConfig.Identifier)
				}
			}

			if f.HasNestedJsonFields() {
				generateEnums(ctx, project, entitiesDir, entity.Entity{
					Identifier: f.JSONConfig.Identifier,
					Fields:     f.JSONConfig.Fields,
				})
			}
		}
	}

}

func generateEnum(ctx context.Context,
	project entity.Project,
	entityDir string,
	e entity.Entity,
	f entity.Field,
	nestedEntity *field.Template,
	pkg string) {

	ft := field.ResolveFieldType(f, e, nestedEntity)

	name := ft.SingularIdentifier

	fmt.Printf("----[GPG] Generating enum: %s\n", name)

	values := make([]string, len(f.OptionValues))
	for i, v := range f.OptionValues {
		values[i] = fmt.Sprintf("%s_%s", strings.ToUpper(name), strings.ToUpper(v.Identifier))
	}

	enumTemplate := EnumTemplate{
		ProjectIdentifier: project.Identifier,
		Package:           pkg,
		EnumName:          helpers.ToCamelCase(name),
		EnumNameUpper:     strings.ToUpper(name),
		Values:            values,
		Options:           f.OptionValues,
	}

	generator.GenerateFile(
		ctx,
		generator.FileRequest{
			OutputFile:   path.Join(entityDir, fmt.Sprintf("%s.go", name)),
			TemplateName: path.Join("core", "enum"),
			Data:         enumTemplate,
		},
	)
}
