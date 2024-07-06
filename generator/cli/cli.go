package cli

import (
	"context"
	"fmt"
	"os/exec"
	"path"
	"text/template"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/helpers"
)

func GenerateCLIModule(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG] Generating CLI\n")
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	cliDir := path.Join(projectDir, generator.CLI_DIR)
	fmt.Printf("--[GPG] CLI Directory: %v\n", cliDir)

	cmd := exec.Command("go", "get", "github.com/manifoldco/promptui")
	cmd.Dir = projectDir
	err := cmd.Run()
	if err != nil {
		fmt.Printf("error running go get github.com/manifoldco/promptui: %v\n", err)
	}

	cmd = exec.Command("go", "get", "github.com/urfave/cli")
	cmd.Dir = projectDir
	err = cmd.Run()
	if err != nil {
		fmt.Printf("error running go get github.com/urfave/cli: %v\n", err)
	}

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(cliDir, "cli.go"),
		TemplateName: path.Join("cli", "cli"),
		Data:         project,
		Funcs: template.FuncMap{
			"ToCamelCase": helpers.ToCamelCase,
			"PrimaryKey": func(e entity.Entity) string {
				return helpers.ToCamelCase(helpers.EntityPrimaryKey(e).Identifier)
			},
		},
	})

	return nil
}
