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

func generateEnums(ctx context.Context, project entity.Project, entityDir string, e entity.Entity) {
	for _, f := range e.Fields {
		if f.Type == entity.OptionsSingleFieldType || f.Type == entity.OptionsManyFieldType {
			generateEnum(ctx, project, entityDir, e, f, nil)
		}
		if f.Type == entity.JSONFieldType {
			for _, jf := range f.JSONConfig.Fields {
				if jf.Type == entity.OptionsSingleFieldType || jf.Type == entity.OptionsManyFieldType {
					ft := field.ResolveFieldType(f, e, nil)
					generateEnum(ctx, project, entityDir, e, jf, &ft)
				}
			}
		}
	}
}

func generateEnum(ctx context.Context,
	project entity.Project,
	entityDir string,
	e entity.Entity,
	f entity.Field,
	nestedEntity *field.Template) {

	ft := field.ResolveFieldType(f, e, nestedEntity)

	name := ft.SingularIdentifier
	if nestedEntity != nil {
		name = fmt.Sprintf("%s_%s", nestedEntity.SingularIdentifier, name)
	}
	fmt.Printf("----[GPG] Generating enum: %s\n", name)

	values := make([]string, len(f.OptionValues))
	for i, v := range f.OptionValues {
		values[i] = fmt.Sprintf("%s_%s", strings.ToUpper(name), strings.ToUpper(v.Identifier))
	}

	enumTemplate := EnumTemplate{
		ProjectName:   project.Identifier,
		Package:       e.Identifier,
		EnumName:      helpers.ToCamelCase(name),
		EnumNameUpper: strings.ToUpper(name),
		Values:        values,
		Options:       f.OptionValues,
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
