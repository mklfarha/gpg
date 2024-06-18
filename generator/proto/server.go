package proto

import (
	"context"
	"fmt"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/helpers"
)

func generateServer(ctx context.Context, protoDir string, project entity.Project, standaloneEntities []ProtoEntityTemplate, dependantEntities map[string][]ProtoEntityTemplate) error {
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(protoDir, "server", "server.go"),
		TemplateName: path.Join("proto", "server"),
		Data: ProtoServiceTemplate{
			Identifier: project.Identifier,
			Name:       helpers.ToCamelCase(project.Identifier),
		},
		DisableGoFormat: false,
	})
	if err != nil {
		return err
	}

	for _, se := range standaloneEntities {
		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:      path.Join(protoDir, "server", fmt.Sprintf("create_%s.go", se.Identifier)),
			TemplateName:    path.Join("proto", "server_create_entity"),
			Data:            se,
			DisableGoFormat: false,
		})
		if err != nil {
			return err
		}
		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:      path.Join(protoDir, "server", fmt.Sprintf("update_%s.go", se.Identifier)),
			TemplateName:    path.Join("proto", "server_update_entity"),
			Data:            se,
			DisableGoFormat: false,
		})
		if err != nil {
			return err
		}
		se.Declarations = getEntityDeclarations(se, dependantEntities, nil)
		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:      path.Join(protoDir, "server", fmt.Sprintf("list_%s.go", se.Identifier)),
			TemplateName:    path.Join("proto", "server_list_entity"),
			Data:            se,
			DisableGoFormat: false,
		})
		if err != nil {
			return err
		}
	}

	return err
}

func getEntityDeclarations(e ProtoEntityTemplate, dependantEntities map[string][]ProtoEntityTemplate, prefix *string) map[string]string {
	res := make(map[string]string)
	for _, f := range e.Fields {
		finalIdentifier := f.Identifier
		if prefix != nil {
			finalIdentifier = *prefix + f.Identifier
		}

		switch f.InternalType {
		case entity.IntFieldType:
			res[finalIdentifier] = "filtering.TypeInt"
		case entity.FloatFieldType:
			res[finalIdentifier] = "filtering.TypeFloat"
		case entity.BooleanFieldType:
			res[finalIdentifier] = "filtering.TypeBool"
		case entity.DateTimeFieldType, entity.DateFieldType:
			res[finalIdentifier] = "filtering.TypeTimestamp"
		case entity.UUIDFieldType, entity.StringFieldType, entity.LargeStringFieldType:
			res[finalIdentifier] = "filtering.TypeString"
		case entity.OptionsSingleFieldType, entity.OptionsManyFieldType:
			res[finalIdentifier] = "filtering.TypeString" // fix this
		case entity.JSONFieldType:

			nestedEntities := dependantEntities[e.Identifier]
			var nestedEntity *ProtoEntityTemplate
			for _, ne := range nestedEntities {
				if ne.Identifier == f.Identifier {
					nestedEntity = &ne
				}
			}
			if nestedEntity != nil {
				parentPrefix := fmt.Sprintf("%s.", f.Identifier)
				nestedDeclarations := getEntityDeclarations(*nestedEntity, dependantEntities, &parentPrefix)
				for ndi, nd := range nestedDeclarations {
					res[ndi] = nd
				}
			}
		}
	}
	fmt.Printf("res: %v\n", res)
	return res
}
