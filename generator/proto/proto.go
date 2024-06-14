package proto

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"text/template"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

type ProtoEntityTemplate struct {
	ProjectIdentifier string
	ParentIdentifier  string
	Identifier        string
	IdentifierPlural  string
	Name              string
	NamePlural        string
	Fields            []field.Template
	Search            bool
	Enums             map[string]ProtoEnumTemplate
	Imports           map[string]interface{}
}

type ProtoEnumTemplate struct {
	Field   field.Template
	Many    bool
	Options []string
}

type ProtoServiceTemplate struct {
	Identifier string
	Name       string
	Entities   []ProtoEntityTemplate
}

func Generate(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG][Proto] Generating Directory\n")
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	protoDir := path.Join(projectDir, generator.PROTO_DIR)

	pl := pluralize.NewClient()

	entityTemplates := []ProtoEntityTemplate{}
	entityNestedTemplates := map[string][]ProtoEntityTemplate{}

	//generate entities/models
	fmt.Printf("--[GPG][Proto] Generating Entities\n")
	for _, e := range project.Entities {
		template, nested, err := handleEntity(ctx, protoDir, project, e, e.Identifier)
		if err != nil {
			return err
		}
		entityTemplates = append(entityTemplates, template)
		nestedTemplates := []ProtoEntityTemplate{}
		for _, n := range nested {
			nestedTemplate, _, err := handleEntity(ctx, protoDir, project, entity.Entity{
				Identifier: pl.Singular(n.Identifier),
				Fields:     n.JSONConfig.Fields,
			}, e.Identifier)
			if err != nil {
				return err
			}
			nestedTemplates = append(nestedTemplates, nestedTemplate)
		}
		entityNestedTemplates[e.Identifier] = nestedTemplates
	}

	//generate project service definition
	fmt.Printf("--[GPG][Proto] Generating Service Definition\n")
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(protoDir, "proto", fmt.Sprintf("service_%s.proto", project.Identifier)),
		TemplateName: path.Join("proto", "service"),
		Data: ProtoServiceTemplate{
			Identifier: project.Identifier,
			Name:       helpers.ToCamelCase(project.Identifier),
			Entities:   entityTemplates,
		},
		DisableGoFormat: true,
		Funcs: template.FuncMap{
			"Inc": helpers.Inc,
		},
	})

	if err != nil {
		return err
	}

	fullDir := path.Join(protoDir, "gen")
	generator.CreateDir(fullDir)

	fmt.Printf("--[GPG][Proto] Generating Go code\n")
	// create gen.sh file
	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(protoDir, "proto", "gen.sh"),
		TemplateName: path.Join("proto", "gen"),
		Data: ProtoServiceTemplate{
			Identifier: project.Identifier,
			Name:       helpers.ToCamelCase(project.Identifier),
			Entities:   entityTemplates,
		},
		DisableGoFormat: true,
	})

	if err != nil {
		return err
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	// run bash file
	cmd := exec.Command("/bin/sh", "./gen.sh")
	cmd.Dir = path.Join(protoDir, "proto")
	cmd.Env = os.Environ()
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	} else {
		fmt.Println("--[GPG][Proto] Proto Go code generated! " + out.String())
	}

	fmt.Printf("--[GPG][Proto] Generating mappers\n")
	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(protoDir, "mapper", "json.go"),
		TemplateName:    path.Join("proto", "mapper_json"),
		DisableGoFormat: false,
	})
	if err != nil {
		return err
	}

	for _, et := range entityTemplates {
		dir := path.Join(protoDir, "mapper", et.Identifier)
		err := handleEntityMapper(ctx, dir, et)
		if err != nil {
			return err
		}

		nestedTemplates := entityNestedTemplates[et.Identifier]
		for _, net := range nestedTemplates {
			err := handleEntityMapper(ctx, dir, net)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func handleEntity(ctx context.Context, protoDir string, project entity.Project, e entity.Entity, parentIdentifier string) (ProtoEntityTemplate, []entity.Field, error) {
	fields := []field.Template{}
	searchable := false
	enums := map[string]ProtoEnumTemplate{}
	nested := []entity.Field{}
	entityTemplate := ProtoEntityTemplate{}
	var err error
	pl := pluralize.NewClient()
	imports := map[string]interface{}{}
	if len(e.Fields) > 0 {
		for _, f := range e.Fields {
			fieldTemplate := field.ResolveFieldType(f, e, nil)
			fields = append(fields, fieldTemplate)
			for _, field := range fields {
				if field.Enum {
					if _, found := enums[field.ProtoType]; !found {
						enums[field.ProtoType] = ProtoEnumTemplate{
							Field:   field,
							Many:    field.EnumMany,
							Options: field.ProtoEnumOptions,
						}
					}
				}
				if f.Type == entity.JSONFieldType {
					if len(f.JSONConfig.Fields) > 0 {
						nested = append(nested, f)
						imports[fmt.Sprintf("%s.proto", pl.Singular(f.Identifier))] = nil
					}
				}
				if f.Type == entity.DateFieldType || f.Type == entity.DateTimeFieldType {
					imports["google/protobuf/timestamp.proto"] = nil
				}
			}

			if f.StorageConfig.Search {
				searchable = true
			}
		}
		entityTemplate = ProtoEntityTemplate{
			ProjectIdentifier: project.Identifier,
			ParentIdentifier:  parentIdentifier,
			Identifier:        e.Identifier,
			IdentifierPlural:  pl.Plural(e.Identifier),
			Name:              helpers.ToCamelCase(e.Identifier),
			NamePlural:        pl.Plural(helpers.ToCamelCase(e.Identifier)),
			Fields:            fields,
			Search:            searchable,
			Enums:             enums,
			Imports:           imports,
		}
		err = generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:      path.Join(protoDir, "proto", fmt.Sprintf("%s.proto", e.Identifier)),
			TemplateName:    path.Join("proto", "model"),
			Data:            entityTemplate,
			DisableGoFormat: true,
			Funcs: template.FuncMap{
				"Inc": helpers.Inc,
			},
		})
	}
	return entityTemplate, nested, err
}

func handleEntityMapper(ctx context.Context, dir string, et ProtoEntityTemplate) error {
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(dir, fmt.Sprintf("%s.go", et.Identifier)),
		TemplateName:    path.Join("proto", "model_mapper"),
		Data:            et,
		DisableGoFormat: false,
		Funcs: template.FuncMap{
			"Inc": helpers.Inc,
		},
	})

	if err != nil {
		return err
	}
	return nil
}
