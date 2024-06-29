package proto

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/field"
)

type ProtoEntityTemplate struct {
	ProjectIdentifier     string
	ParentIdentifier      string
	OrignalIdentifier     string
	FinalIdentifier       string
	FinalIdentifierPlural string
	Name                  string
	NamePlural            string
	Type                  string
	Fields                []field.Template
	PrimaryKey            field.Template
	Search                bool
	Enums                 map[string]ProtoEnumTemplate
	Imports               map[string]interface{}
	Declarations          []ProtoEntityDeclaration
}

type ProtoEntityDeclaration struct {
	Identifier  string
	Fields      []ProtoFieldDeclaration
	IsDependant bool
}

type ProtoFieldDeclaration struct {
	Identifier string
	Name       string
	Filtering  string
	IsEnum     bool
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

	err := os.RemoveAll(protoDir)
	if err != nil {
		fmt.Printf("ERROR: Deleting module directory\n")
	}

	fullDir := path.Join(protoDir, "gen")
	generator.CreateDir(fullDir)

	// generate proto files
	standaloneEntities, dependantEntities, err := generateProtoFiles(ctx, protoDir, project)
	if err != nil {
		return err
	}

	// generate base go code with protoc
	err = generateProtoc(ctx, protoDir, project, standaloneEntities)
	if err != nil {
		return err
	}

	// generate mappers to/from entity/proto
	err = generateMappers(ctx, protoDir, project, standaloneEntities, dependantEntities)
	if err != nil {
		return err
	}

	err = generateServer(ctx, protoDir, project, standaloneEntities, dependantEntities)
	if err != nil {
		return err
	}

	return err
}
