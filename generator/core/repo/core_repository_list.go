package repo

import (
	"context"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
)

func generateRepositoryListCode(ctx context.Context, repoDir string, project entity.Project) error {
	listDir := path.Join(repoDir, "list")

	// list
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(listDir, "list.go"),
		TemplateName: path.Join("core", "repo", "list", "repo_list"),
		Data: struct {
			ProjectIdentifier string
		}{
			ProjectIdentifier: project.Identifier,
		},
	})
	if err != nil {
		return err
	}

	// list_fields
	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(listDir, "list_fields.go"),
		TemplateName: path.Join("core", "repo", "list", "repo_list_fields"),
		Data: struct {
			ProjectIdentifier string
		}{
			ProjectIdentifier: project.Identifier,
		},
	})
	if err != nil {
		return err
	}

	// types
	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(listDir, "types.go"),
		TemplateName: path.Join("core", "repo", "list", "repo_list_types"),
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
