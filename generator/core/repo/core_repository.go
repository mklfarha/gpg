package repo

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
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
			ProjectModule     string
		}{
			ProjectIdentifier: project.Identifier,
			ProjectModule:     project.Module,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func SyncSchema(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG] Generating core repository\n")
	projectDir := generator.ProjectDir(ctx, rootPath, project)

	var out bytes.Buffer
	var stderr bytes.Buffer
	if !generator.FileExists(path.Join(projectDir, "go.mod")) {
		cmd := exec.Command("go", "mod", "init", project.Identifier)
		cmd.Dir = projectDir
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("error running go mod init: %v | %v | %v\n", err, out.String(), stderr.String())
		}
	} else {
		fmt.Printf("----[GPG] go.mod already exists\n")
	}

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

	// install sqlc
	cmd := exec.Command("go", "get", "github.com/skeema/skeema")
	cmd.Dir = projectDir
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("error running go install skeema: %v | %v | %v\n", err, out.String(), stderr.String())
	}

	// ignoring error on purpose
	executeSkeema(ctx, project, sqlDir)

	return err
}
