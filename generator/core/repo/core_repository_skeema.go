package repo

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
)

func executeSkeema(ctx context.Context, project entity.Project, sqlDir string) error {

	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(sqlDir, "schemas", ".skeema"),
		TemplateName:    path.Join("core", "repo", "repo_skeema"),
		Data:            project,
		DisableGoFormat: true,
	})
	if err != nil {
		return err
	}

	// try to sync db changes with skeema
	var out bytes.Buffer
	var stderr bytes.Buffer

	fmt.Printf("----[GPG][Skeema] Sync DB - Diff \n")
	cmd := exec.Command("go", "run", "github.com/skeema/skeema", "diff", "--allow-unsafe")
	cmd.Dir = path.Join(sqlDir, "schemas")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("skeema diff 1: " + fmt.Sprint(err) + ": " + stderr.String())
		fmt.Println("skeema diff 2: " + out.String())
		//return err - don't return error as the command sends exit (1) even if it succeeded
	}

	fmt.Printf("----[GPG][Skeema] Sync DB - Push \n")
	cmd = exec.Command("go", "run", "github.com/skeema/skeema", "push", "--allow-unsafe")
	cmd.Dir = path.Join(sqlDir, "schemas")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("skeema psuh 1: " + fmt.Sprint(err) + ": " + stderr.String())
		fmt.Println("skeema psuh 2: " + out.String())
		//return err - don't return error as the command sends exit (1) even if it succeeded
	}
	return nil
}
