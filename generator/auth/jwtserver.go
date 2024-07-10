package auth

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
)

func generateBasicJWTServer(ctx context.Context, authDir string, project entity.Project) error {

	if !project.HasJWTAuth() {
		return errors.New("invalid auth type")
	}

	if project.JWTAuth().Config.JWT == nil {
		return errors.New("missing JWT config")
	}

	jwtServerDir := path.Join(authDir, "jwtserver")
	fmt.Printf("--[GPG][AUTH] Generating JWT server\n")
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(jwtServerDir, "server.go"),
		TemplateName: path.Join("auth", "jwtserver", "jwt_server"),
		Data:         project,
	})
	if err != nil {
		return err
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(jwtServerDir, "types.go"),
		TemplateName: path.Join("auth", "jwtserver", "jwt_types"),
		Data:         project.JWTAuth().Config.JWT,
	})
	if err != nil {
		return err
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(jwtServerDir, "parse.go"),
		TemplateName: path.Join("auth", "jwtserver", "jwt_parse"),
		Data:         project,
	})
	if err != nil {
		return err
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(jwtServerDir, "refresh.go"),
		TemplateName: path.Join("auth", "jwtserver", "jwt_refresh"),
		Data:         project,
	})
	if err != nil {
		return err
	}

	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(jwtServerDir, "signin.go"),
		TemplateName: path.Join("auth", "jwtserver", "jwt_signin"),
		Data:         project,
	})
	if err != nil {
		return err
	}

	return generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(jwtServerDir, "validate.go"),
		TemplateName: path.Join("auth", "jwtserver", "jwt_validate"),
		Data:         project,
	})
}
