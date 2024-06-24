package generator

import (
	"context"
	"fmt"
	"os/exec"
	"path"

	"github.com/maykel/gpg/entity"
)

func GenerateAPIModule(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG] Generating api\n")
	projectDir := ProjectDir(ctx, rootPath, project)
	apiDir := path.Join(projectDir, API_DIR)
	fmt.Printf("--[GPG] API Directory: %v\n", apiDir)

	GenerateFile(ctx, FileRequest{
		OutputFile:   path.Join(apiDir, "main.go"),
		TemplateName: path.Join("api", "api_main"),
		Data:         project,
	})

	GenerateFile(ctx, FileRequest{
		OutputFile:   path.Join(apiDir, "server.go"),
		TemplateName: path.Join("api", "api_server"),
		Data:         project,
	})

	cmd := exec.Command("go", "get", "github.com/aws/aws-sdk-go/...")
	cmd.Dir = projectDir
	err := cmd.Run()
	if err != nil {
		fmt.Printf("error running go get github.com/aws/aws-sdk-go/...: %v\n", err)
	}

	GenerateFile(ctx, FileRequest{
		OutputFile:   path.Join(apiDir, "aws_s3.go"),
		TemplateName: path.Join("api", "api_aws_s3"),
		Data:         project,
	})

	return nil
}
