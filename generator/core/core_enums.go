package core

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/helpers"
)

func generateEnums(ctx context.Context, entityDir string, e entity.Entity) {
	for _, f := range e.Fields {
		if f.Type == entity.OptionsSingleFieldType || f.Type == entity.OptionsManyFieldType {
			generateEnum(ctx, entityDir, e, f, nil)
		}
		if f.Type == entity.JSONFieldType {
			for _, jf := range f.JSONConfig.Fields {
				if jf.Type == entity.OptionsSingleFieldType || jf.Type == entity.OptionsManyFieldType {
					generateEnum(ctx, entityDir, e, jf, &f.Identifier)
				}
			}
		}
	}
}

func generateEnum(ctx context.Context,
	entityDir string,
	e entity.Entity,
	f entity.Field,
	prefix *string) {

	pl := pluralize.NewClient()
	name := pl.Singular(f.Identifier)
	if prefix != nil {
		name = fmt.Sprintf("%s_%s", pl.Singular(*prefix), name)
	}
	fmt.Printf("----[GPG] Generating enum: %s\n", name)

	values := make([]string, len(f.OptionValues))
	for i, v := range f.OptionValues {
		values[i] = fmt.Sprintf("%s_%s", strings.ToUpper(name), strings.ToUpper(v.Identifier))
	}

	enumTemplate := EnumTemplate{
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
