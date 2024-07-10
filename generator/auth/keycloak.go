package auth

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
)

func generateKeycloakClient(ctx context.Context, authDir string, project entity.Project) error {
	if !project.HasKeycloakAuth() {
		return errors.New("invalid auth type")
	}

	if project.KeycloakAuth().Config.Keycloak == nil {
		return errors.New("missing keycloak config")
	}

	kcServerDir := path.Join(authDir, "keycloak")
	fmt.Printf("--[GPG][AUTH] Generating keycloak client\n")
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(kcServerDir, "client.go"),
		TemplateName: path.Join("auth", "keycloak", "keycloak_client"),
		Data:         project,
	})
	if err != nil {
		return err
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(kcServerDir, "types.go"),
		TemplateName: path.Join("auth", "keycloak", "keycloak_types"),
		Data:         project,
	})
	if err != nil {
		return err
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(kcServerDir, "handle_http.go"),
		TemplateName: path.Join("auth", "keycloak", "keycloak_handle_http"),
		Data:         project,
	})
	if err != nil {
		return err
	}

	return nil
}
