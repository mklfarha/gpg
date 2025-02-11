package repo

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path"
	"text/template"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/helpers"
)

func generateRepositorySQLCode(ctx context.Context, repoDir string, project entity.Project) error {
	// generate sqlc yaml file
	err := generator.GenerateFile(
		ctx,
		generator.FileRequest{
			OutputFile:   path.Join(repoDir, "sqlc.yaml"),
			TemplateName: path.Join("core", "repo", "repo_yaml"),
			Data: struct {
				ProjectIdentifier string
				ProjectModule     string
				Fields            map[string]string
			}{
				ProjectIdentifier: project.Identifier,
				ProjectModule:     project.Module,
				Fields:            helpers.FieldsToCamelCase(project.Entities),
			},
			DisableGoFormat: true,
			Funcs: template.FuncMap{
				"StringContains": helpers.StringContains,
				"ToCamelCase":    helpers.ToCamelCase,
			},
		},
	)
	if err != nil {
		return err
	}

	fmt.Printf("----[GPG] SQLC Generate\n")
	cmd := exec.Command("go", "run", "github.com/sqlc-dev/sqlc/cmd/sqlc", "generate")
	cmd.Dir = repoDir
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		fmt.Println("SQLC Result: " + out.String())
		return err
	}

	return nil
}
