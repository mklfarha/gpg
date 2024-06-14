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
	Identifier        string
	IdentifierPlural  string
	EntityName        string
	EntityNamePlural  string
	Fields            []field.Template
	Search            bool
	Enums             map[string]ProtoEnumTemplate
	Imports           map[string]interface{}
}

type ProtoEnumTemplate struct {
	Field   field.Template
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

	//generate entities/models
	fmt.Printf("--[GPG][Proto] Generating Entities\n")
	for _, e := range project.Entities {
		template, nested, err := handleEntity(ctx, protoDir, project, e)
		if err != nil {
			return err
		}
		entityTemplates = append(entityTemplates, template)
		for _, n := range nested {
			_, _, err := handleEntity(ctx, protoDir, project, entity.Entity{
				Identifier: pl.Singular(n.Identifier),
				Fields:     n.JSONConfig.Fields,
			})
			if err != nil {
				return err
			}
		}
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
	// delete existing go files
	cmd := exec.Command("find", "./gen", "-name", "\"*.go\"", "-type", "f", "-delete")
	cmd.Dir = protoDir
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		fmt.Println("find delete: " + out.String())
	}

	cmd = exec.Command("go", "install", "google.golang.org/protobuf/cmd/protoc-gen-go@latest")
	cmd.Dir = projectDir
	err = cmd.Run()
	if err != nil {
		fmt.Printf("error running go get proto-gen-go\n")
	}

	cmd = exec.Command("go", "install", "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest")
	cmd.Dir = projectDir
	err = cmd.Run()
	if err != nil {
		fmt.Printf("error running go get proto-gen-go-grpc\n")
	}

	cmd = exec.Command("sh", "-c", "export PATH=\"$PATH:$(go env GOPATH)/bin\"")
	cmd.Dir = projectDir
	cmd.Stdout = &out
	err = cmd.Run()
	fmt.Println("export: " + out.String())
	if err != nil {
		fmt.Printf("error running export\n")
	}

	/*cmd = exec.Command("sh", "-c", "echo $(go env GOPATH)")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Run()
	fmt.Println("echo: " + out.String())*/

	// generate go code
	// protoc --go_out=./../go --go-grpc_out=./../go ./*.proto
	cmd = exec.Command("protoc", "--go_out=./idl/gen", "--go-grpc_out=./idl/gen", "./idl/proto/*.proto")
	cmd.Dir = projectDir
	cmd.Env = os.Environ()
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		fmt.Println("protoc out: " + out.String())
	}

	return err
}

func handleEntity(ctx context.Context, protoDir string, project entity.Project, e entity.Entity) (ProtoEntityTemplate, []entity.Field, error) {
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
							Options: field.ProtoEnumOptions,
						}
					}
				}
				if f.Type == entity.JSONFieldType {
					nested = append(nested, f)
					if len(f.JSONConfig.Fields) > 0 {
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
			Identifier:        e.Identifier,
			IdentifierPlural:  pl.Plural(e.Identifier),
			EntityName:        helpers.ToCamelCase(e.Identifier),
			EntityNamePlural:  pl.Plural(helpers.ToCamelCase(e.Identifier)),
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
