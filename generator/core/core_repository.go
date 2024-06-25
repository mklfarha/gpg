package core

import (
	"context"
	"fmt"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
)

func GenerateCoreRepository(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG] Generating core repository\n")
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	repoDir := path.Join(projectDir, generator.CORE_REPO_DIR)

	sqlDir := path.Join(repoDir, generator.CORE_REPO_SQL_DIR)

	// generate sql files
	err := generateRepositorySQL(ctx, project, sqlDir)
	if err != nil {
		return err
	}

	// generate go code with SQLC
	err = generateRepositorySQLCode(ctx, repoDir, project)
	if err != nil {
		return err
	}

	// TODO: make this optional
	// ignoring error on purpose
	executeSkeema(ctx, project, sqlDir)

	// list
	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(repoDir, "list.go"),
		TemplateName: path.Join("core", "repo", "repo_list"),
		Data: struct {
			ProjectName string
		}{
			ProjectName: project.Identifier,
		},
	})
	if err != nil {
		return err
	}

	// new function to return generated code module
	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(repoDir, "repository.go"),
		TemplateName: path.Join("core", "repo", "repository"),
		Data: struct {
			ProjectName string
		}{
			ProjectName: project.Identifier,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
