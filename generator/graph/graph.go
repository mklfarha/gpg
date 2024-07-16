package graph

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/core/repo"
	"github.com/maykel/gpg/generator/field"
)

type GraphEntityTemplate struct {
	Identifier       string
	EntityType       string
	EntityTypePlural string
	JSON             bool
	JSONMany         bool
	Required         bool
	ParentIdentifier string
	ParentEntityName string
	GraphGenType     string
	PrimaryKey       field.Template
	InFields         []field.Template
	OutFields        []field.Template
	Selects          []repo.SchemaSelectStatement
	CustomQueries    []entity.CustomQuery
	Search           bool
}

type GraphQueriesTemplate struct {
	ProjectName  string
	Project      entity.Project
	Entities     []GraphEntityTemplate
	JSONEntities []GraphEntityTemplate
	Enums        []field.Template
}

func GenerateGraph(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG][GraphQL] Generating GraphQL\n")
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	graphDir := path.Join(projectDir, generator.GRAPH_DIR)

	// config file
	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(graphDir, "gqlgen.yml"),
		TemplateName:    path.Join("graph", "graph_yaml"),
		DisableGoFormat: true,
	})

	// generate entities
	res, err := generateEntities(ctx, graphDir, project)
	if err != nil {
		return err
	}

	entityTemplates := res.EntityTemplates

	// generate queries and mutations
	err = generateQueries(ctx, graphDir, entityTemplates, project)
	if err != nil {
		return err
	}

	// install gqlgen
	cmd := exec.Command("go", "install", "github.com/99designs/gqlgen")
	cmd.Dir = graphDir
	err = cmd.Run()
	if err != nil {
		fmt.Printf("error running go get gqlgen\n")
	}

	// generate code with gqlgen
	fmt.Printf("----[GPG][GraphQL] GQLGEN generate\n")
	cmd = exec.Command("go", "run", "github.com/99designs/gqlgen", "generate")
	cmd.Dir = graphDir
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		fmt.Println("gqlgen result: " + out.String())
	}

	// overrides some of the generated file to add content
	err = overrideGqlgenFiles(ctx, graphDir, entityTemplates, project)
	if err != nil {
		return err
	}

	// generate mappers
	err = generateMapper(ctx, graphDir, project, res)
	if err != nil {
		return err
	}

	// generate server file
	fmt.Printf("----[GPG][GraphQL] Server\n")
	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(graphDir, "server.go"),
		TemplateName: path.Join("graph", "graph_server"),
		Data:         project,
	})

	return nil
}
