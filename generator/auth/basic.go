package auth

import (
	"context"
	"fmt"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
)

func generateBasicAuth(ctx context.Context, authDir string, project entity.Project) error {
	fmt.Printf("--[GPG] Generating basic auth\n")
	return generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(authDir, "basic.go"),
		TemplateName: path.Join("auth", "auth_basic"),
		Data:         project,
	})
}
