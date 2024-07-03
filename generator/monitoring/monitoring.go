package monitoring

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
)

func GenerateMonitoring(ctx context.Context, rootPath string, project entity.Project) error {
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	monitoringDir := path.Join(projectDir, generator.MONITORING_DIR)

	err := os.RemoveAll(monitoringDir)
	if err != nil {
		fmt.Printf("ERROR: Deleting monitoring directory\n")
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(monitoringDir, "types.go"),
		TemplateName: path.Join("monitoring", "types"),
	})
	if err != nil {
		return err
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(monitoringDir, "monitoring.go"),
		TemplateName: path.Join("monitoring", "monitoring"),
	})
	if err != nil {
		return err
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(monitoringDir, "metrics.go"),
		TemplateName: path.Join("monitoring", "metrics"),
	})
	if err != nil {
		return err
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(monitoringDir, "metrics_datadog.go"),
		TemplateName: path.Join("monitoring", "metrics_datadog"),
	})
	if err != nil {
		return err
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(monitoringDir, "logging.go"),
		TemplateName: path.Join("monitoring", "logging"),
	})
	if err != nil {
		return err
	}

	return nil
}
