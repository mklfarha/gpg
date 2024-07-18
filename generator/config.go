package generator

import (
	"context"
	"fmt"
	"path"

	"github.com/maykel/gpg/entity"
)

func GenerateConfig(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG] Generating config\n")
	projectDir := ProjectDir(ctx, rootPath, project)
	configDir := path.Join(projectDir, CONFIG_DIR)
	GenerateFile(ctx, FileRequest{
		OutputFile:   path.Join(configDir, "config.go"),
		TemplateName: path.Join("config", "config"),
		Data:         project,
	})

	GenerateFile(ctx, FileRequest{
		OutputFile:      path.Join(configDir, "base.yaml"),
		TemplateName:    path.Join("config", "config_base"),
		Data:            project,
		DisableGoFormat: true,
	})
	GenerateFile(ctx, FileRequest{
		OutputFile:      path.Join(configDir, "cli.yaml"),
		TemplateName:    path.Join("config", "config_cli"),
		Data:            project,
		DisableGoFormat: true,
	})

	return nil
}
