package repo

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
)

func GenerateCoreRepository(ctx context.Context, rootPath string, project entity.Project, skipSkeema bool) error {
	fmt.Printf("--[GPG] Generating core repository\n")
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	repoDir := path.Join(projectDir, generator.CORE_REPO_DIR)

	sqlDir := path.Join(repoDir, generator.CORE_REPO_SQL_DIR)

	err := os.RemoveAll(repoDir)
	if err != nil {
		fmt.Printf("ERROR: Deleting repo directory\n")
	}

	// generate sql files
	err = generateRepositorySQL(ctx, project, sqlDir)
	if err != nil {
		return err
	}

	// generate go code with SQLC
	err = generateRepositorySQLCode(ctx, repoDir, project)
	if err != nil {
		return err
	}

	if !skipSkeema {
		// ignoring error on purpose
		executeSkeema(ctx, project, sqlDir)
	}

	// list module
	err = generateRepositoryListCode(ctx, repoDir, project)
	if err != nil {
		return err
	}

	// new function to return generated code module
	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(repoDir, "repository.go"),
		TemplateName: path.Join("core", "repo", "repository"),
		Data: struct {
			ProjectIdentifier string
		}{
			ProjectIdentifier: project.Identifier,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
