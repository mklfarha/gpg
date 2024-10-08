package aws

import (
	"context"
	"fmt"
	"os/exec"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
)

func GenerateAWSModule(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG] Generating AWS\n")
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	awsDir := path.Join(projectDir, generator.AWS_DIR)
	fmt.Printf("--[GPG] AWS Directory: %v\n", awsDir)

	cmd := exec.Command("go", "get", "github.com/aws/aws-sdk-go/...")
	cmd.Dir = projectDir
	err := cmd.Run()
	if err != nil {
		fmt.Printf("error running go get github.com/aws/aws-sdk-go/...: %v\n", err)
	}

	return generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(awsDir, "client.go"),
		TemplateName: path.Join("aws", "client"),
		Data:         project,
	})
}
