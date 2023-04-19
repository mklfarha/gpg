package generator

import (
	"context"
	"fmt"
	"path"

	"github.com/maykel/gpg/entity"
)

func GenerateAuth(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG] Generating auth\n")
	projectDir := ProjectDir(ctx, rootPath, project)
	authDir := path.Join(projectDir, AUTH_DIR)

	GenerateFile(ctx, FileRequest{
		OutputFile:   path.Join(authDir, "auth.go"),
		TemplateName: path.Join("auth", "auth"),
		Data:         project,
	})

	GenerateFile(ctx, FileRequest{
		OutputFile:   path.Join(authDir, "models.go"),
		TemplateName: path.Join("auth", "auth_models"),
		Data:         project,
	})

	GenerateFile(ctx, FileRequest{
		OutputFile:   path.Join(authDir, "parse.go"),
		TemplateName: path.Join("auth", "auth_parse"),
		Data:         project,
	})

	GenerateFile(ctx, FileRequest{
		OutputFile:   path.Join(authDir, "refresh.go"),
		TemplateName: path.Join("auth", "auth_refresh"),
		Data:         project,
	})

	GenerateFile(ctx, FileRequest{
		OutputFile:   path.Join(authDir, "signin.go"),
		TemplateName: path.Join("auth", "auth_signin"),
		Data:         project,
	})

	GenerateFile(ctx, FileRequest{
		OutputFile:   path.Join(authDir, "validate.go"),
		TemplateName: path.Join("auth", "auth_validate"),
		Data:         project,
	})

	return nil
}
