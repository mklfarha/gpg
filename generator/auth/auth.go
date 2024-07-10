package auth

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
)

func GenerateAuth(ctx context.Context, rootPath string, project entity.Project) error {
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	authDir := path.Join(projectDir, generator.AUTH_DIR)
	err := os.RemoveAll(authDir)
	if err != nil {
		fmt.Printf("ERROR: Deleting module directory\n")
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(authDir, "types.go"),
		TemplateName: path.Join("auth", "auth_types"),
		Data:         project,
	})
	if err != nil {
		return err
	}

	if project.HasBasicAuth() {
		err = generateBasicAuth(ctx, authDir, project)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("--[GPG][AUTH] Skip basic auth\n")
	}

	if project.HasJWTAuth() {
		err := generateBasicJWTServer(ctx, authDir, project)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("--[GPG][AUTH] Skip jwt auth \n")
	}

	if project.HasKeycloakAuth() {
		err := generateKeycloakClient(ctx, authDir, project)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("--[GPG][AUTH] Skipp keycloak auth \n")
	}

	return nil
}
