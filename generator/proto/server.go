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
	fmt.Printf("--[GPG][Proto] Generating server.go\n")
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(protoDir, "server", "server.go"),
		TemplateName: path.Join("proto", "server"),
		Data: ProtoServiceTemplate{
			Identifier: project.Identifier,
			Module:     project.Module,
			Name:       helpers.ToCamelCase(project.Identifier),
			AuthImport: project.AuthImport(),
		},
		DisableGoFormat: false,
	})
	if err != nil {
		return err
	}

	fmt.Printf("--[GPG][Proto] Generating auth.go\n")
	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(protoDir, "server", "auth.go"),
		TemplateName: path.Join("proto", "server_auth"),
		Data: ProtoServiceTemplate{
			Identifier: project.Identifier,
			Module:     project.Module,
			Name:       helpers.ToCamelCase(project.Identifier),
			AuthImport: project.AuthImport(),
		},
		DisableGoFormat: false,
	})
	if err != nil {
		return err
	}

	for _, se := range standaloneEntities {
		fmt.Printf("--[GPG][Proto] Generating create: %v\n", se.FinalIdentifier)
		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:      path.Join(protoDir, "server", fmt.Sprintf("create_%s.go", se.FinalIdentifier)),
			TemplateName:    path.Join("proto", "server_create_entity"),
			Data:            se,
			DisableGoFormat: false,
		})
		if err != nil {
			return err
		}

		fmt.Printf("--[GPG][Proto] Generating update: %v\n", se.FinalIdentifier)
		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:      path.Join(protoDir, "server", fmt.Sprintf("update_%s.go", se.FinalIdentifier)),
			TemplateName:    path.Join("proto", "server_update_entity"),
			Data:            se,
			DisableGoFormat: false,
		})
		if err != nil {
			return err
		}

		se.Declarations = getEntityDeclarations(se, dependantEntities, nil)
		fmt.Printf("--[GPG][Proto] Generating list: %v\n", se.FinalIdentifier)
		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:      path.Join(protoDir, "server", fmt.Sprintf("list_%s.go", se.FinalIdentifier)),
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

func getEntityDeclarations(e ProtoEntityTemplate, dependantEntities map[string][]ProtoEntityTemplate, nestedEntity *ProtoEntityTemplate) []ProtoEntityDeclaration {
	finalRes := []ProtoEntityDeclaration{}
	isDependant := false
	if nestedEntity != nil {
		isDependant = true
	}
	entityRes := ProtoEntityDeclaration{
		Identifier:  e.FinalIdentifier,
		IsDependant: isDependant,
		Fields:      []ProtoFieldDeclaration{},
	}
	for _, f := range e.Fields {
		finalIdentifier := f.Identifier
		if nestedEntity != nil {
			finalIdentifier = fmt.Sprintf("%s.%s", nestedEntity.OrignalIdentifier, f.Identifier)
		}
		switch f.InternalType {
		case entity.IntFieldType:
			entityRes.Fields = append(entityRes.Fields, ProtoFieldDeclaration{
				Name:      finalIdentifier,
				Filtering: "filtering.TypeInt",
				IsEnum:    false,
			})
		case entity.FloatFieldType:
			entityRes.Fields = append(entityRes.Fields, ProtoFieldDeclaration{
				Name:      finalIdentifier,
				Filtering: "filtering.TypeFloat",
				IsEnum:    false,
			})
		case entity.BooleanFieldType:
			entityRes.Fields = append(entityRes.Fields, ProtoFieldDeclaration{
				Name:      finalIdentifier,
				Filtering: "filtering.TypeBool",
				IsEnum:    false,
			})
		case entity.DateTimeFieldType, entity.DateFieldType:
			entityRes.Fields = append(entityRes.Fields, ProtoFieldDeclaration{
				Name:      finalIdentifier,
				Filtering: "filtering.TypeTimestamp",
				IsEnum:    false,
			})
		case entity.UUIDFieldType, entity.StringFieldType, entity.LargeStringFieldType:
			entityRes.Fields = append(entityRes.Fields, ProtoFieldDeclaration{
				Name:      finalIdentifier,
				Filtering: "filtering.TypeString",
				IsEnum:    false,
			})
		case entity.ArrayFieldType:
			filtering := ""
			switch f.ArrayInternalType {
			case entity.UUIDFieldType, entity.StringFieldType:
				filtering = "filtering.TypeString"
			case entity.IntFieldType:
				filtering = "filtering.TypeInt"
			case entity.FloatFieldType:
				filtering = "filtering.TypeFloat"
			case entity.BooleanFieldType:
				filtering = "filtering.TypeBool"
			case entity.DateTimeFieldType, entity.DateFieldType:
				filtering = "filtering.TypeTimestamp"
			}
			entityRes.Fields = append(entityRes.Fields, ProtoFieldDeclaration{
				Name:      finalIdentifier,
				Filtering: filtering,
				IsEnum:    false,
			})
		case entity.OptionsSingleFieldType, entity.OptionsManyFieldType:
			enumType := f.ProtoType
			entityRes.Fields = append(entityRes.Fields, ProtoFieldDeclaration{
				Name:      finalIdentifier,
				Filtering: fmt.Sprintf("pb.%s(0).Type()", enumType),
				IsEnum:    true,
			})
		case entity.JSONFieldType:
			nestedEntities, found := dependantEntities[e.OrignalIdentifier]
			if found {
				var nestedEntity ProtoEntityTemplate
				for _, ne := range nestedEntities {
					if ne.OrignalIdentifier == f.SingularIdentifier {
						nestedEntity = ne
					}
				}

				nestedEntityDeclarations := getEntityDeclarations(nestedEntity, dependantEntities, &nestedEntity)
				finalRes = append(finalRes, nestedEntityDeclarations...)
			} else {
				fmt.Printf("nested entity not found: %v", f.Identifier)
			}
		}
	}

	finalRes = append(finalRes, entityRes)
	return finalRes
}
