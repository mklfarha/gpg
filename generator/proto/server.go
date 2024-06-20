package proto

import (
	"context"
	"fmt"
	"path"
	"sort"

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
		se.DeclarationKeys = make(map[string][]string)
		for entityId, mapTypes := range se.Declarations {
			for k, _ := range mapTypes {
				se.DeclarationKeys[entityId] = append(se.DeclarationKeys[entityId], k)
			}
		}

		for _, dk := range se.DeclarationKeys {
			sort.Strings(dk)
		}

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

func getEntityDeclarations(e ProtoEntityTemplate, dependantEntities map[string][]ProtoEntityTemplate, prefix *string) map[string]map[string]string {
	finalRes := make(map[string]map[string]string)
	entityRes := make(map[string]string)
	for _, f := range e.Fields {
		finalIdentifier := f.Identifier
		if prefix != nil {
			finalIdentifier = *prefix + f.Identifier
		}
		switch f.InternalType {
		case entity.IntFieldType:
			entityRes[finalIdentifier] = "filtering.TypeInt"
		case entity.FloatFieldType:
			entityRes[finalIdentifier] = "filtering.TypeFloat"
		case entity.BooleanFieldType:
			entityRes[finalIdentifier] = "filtering.TypeBool"
		case entity.DateTimeFieldType, entity.DateFieldType:
			entityRes[finalIdentifier] = "filtering.TypeTimestamp"
		case entity.UUIDFieldType, entity.StringFieldType, entity.LargeStringFieldType:
			entityRes[finalIdentifier] = "filtering.TypeString"
		case entity.OptionsSingleFieldType, entity.OptionsManyFieldType:
			entityRes[finalIdentifier] = fmt.Sprintf("filtering.TypeEnum(pb.%s(0).Type())", f.ProtoType)
		case entity.JSONFieldType:
			nestedEntities := dependantEntities[e.Identifier]
			var nestedEntity ProtoEntityTemplate
			for _, ne := range nestedEntities {
				if ne.Identifier == f.Identifier {
					nestedEntity = ne
				}
			}

			parentPrefix := fmt.Sprintf("%s.", f.Identifier)
			nestedDeclarations := getEntityDeclarations(nestedEntity, dependantEntities, &parentPrefix)
			nestedRes := make(map[string]string)
			for ndi, nd := range nestedDeclarations[nestedEntity.Identifier] {
				nestedRes[ndi] = nd
			}
			finalRes[f.Identifier] = nestedRes
		}
	}

	finalRes[e.Identifier] = entityRes
	return finalRes
}
