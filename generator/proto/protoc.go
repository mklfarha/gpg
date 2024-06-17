package proto

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/helpers"
)

func generateProtoc(ctx context.Context, protoDir string, project entity.Project, standaloneEntities []ProtoEntityTemplate) error {
	fmt.Printf("--[GPG][Proto] Generating Go code\n")
	// create gen.sh file
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(protoDir, "proto", "gen.sh"),
		TemplateName: path.Join("proto", "gen"),
		Data: ProtoServiceTemplate{
			Identifier: project.Identifier,
			Name:       helpers.ToCamelCase(project.Identifier),
			Entities:   standaloneEntities,
		},
		DisableGoFormat: true,
	})

	if err != nil {
		return err
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	// run bash file
	cmd := exec.Command("/bin/sh", "./gen.sh")
	cmd.Dir = path.Join(protoDir, "proto")
	cmd.Env = os.Environ()
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	} else {
		fmt.Println("--[GPG][Proto] Proto Go code generated! " + out.String())
	}
	return err
}
