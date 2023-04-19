package generator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/maykel/gpg/entity"
)

var (
	API_DIR                   = "api"
	AUTH_DIR                  = "auth"
	CONFIG_DIR                = "config"
	GRAPH_DIR                 = "graph"
	CORE_DIR                  = "core"
	CORE_ENTITY_DIR           = path.Join(CORE_DIR, "entity")
	CORE_MODULE_DIR           = path.Join(CORE_DIR, "module")
	CORE_REPO_DIR             = path.Join(CORE_DIR, "repository")
	CORE_REPO_SQL_DIR         = "sql"
	CORE_REPO_SQL_QUERIES_DIR = path.Join(CORE_REPO_SQL_DIR, "queries")
	CORE_TOOLS                = "tools"
	WEB_DIR                   = "web"
	CUSTOM_DIR                = "custom"
)

var projectStructure = []string{
	API_DIR,
	AUTH_DIR,
	CONFIG_DIR,
	GRAPH_DIR,
	CORE_ENTITY_DIR,
	CORE_MODULE_DIR,
	CORE_REPO_DIR,
	CORE_TOOLS,
	WEB_DIR,
	CUSTOM_DIR,
}

func GenerateProjectDirectories(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG] Generating Directories \n")
	projectDir := ProjectDir(ctx, rootPath, project)
	for _, dir := range projectStructure {
		fullDir := path.Join(projectDir, dir)
		createDir(fullDir)
	}

	initModule(ctx, rootPath, project)

	return nil
}

func initModule(ctx context.Context, rootPath string, project entity.Project) {
	fmt.Printf("--[GPG] Init go module \n")
	projectDir := ProjectDir(ctx, rootPath, project)
	if !fileExists(path.Join(projectDir, "go.mod")) {
		cmd := exec.Command("go", "mod", "init", project.Identifier)
		cmd.Dir = projectDir
		err := cmd.Run()
		if err != nil {
			fmt.Printf("error running go mod init\n")
		}
	} else {
		fmt.Printf("----[GPG] go.mod already exists\n")
	}

	GenerateFile(ctx, FileRequest{
		OutputFile:   path.Join(projectDir, "tools", "tools.go"),
		TemplateName: "tools",
	})

	// go version
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("go", "version")
	cmd.Dir = projectDir
	cmd.Run()
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	fmt.Printf("--[GPG] go version: %v\n", out.String())

	// install sqlc
	cmd = exec.Command("go", "get", "github.com/kyleconroy/sqlc/cmd/sqlc")
	cmd.Dir = projectDir
	err := cmd.Run()
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err != nil {
		fmt.Printf("error running go insall sqlc -o: %v\n", out.String())
		fmt.Printf("error running go insall sqlc -err: %v\n", stderr.String())
	}

	// install sqlc
	cmd = exec.Command("go", "get", "golang.org/x/sync/errgroup")
	cmd.Dir = projectDir
	err = cmd.Run()
	if err != nil {
		fmt.Printf("error running go get errgroup\n")
	}

	// install gqlgen
	cmd = exec.Command("go", "get", "github.com/99designs/gqlgen")
	cmd.Dir = projectDir
	err = cmd.Run()
	if err != nil {
		fmt.Printf("error running go get gqlgen\n")
	}

	// go mod tidy
	GoModTidy(ctx, rootPath, project)

}

func GoModTidy(ctx context.Context, rootPath string, project entity.Project) {
	fmt.Printf("--[GPG] Go Mod Tidy \n")
	projectDir := ProjectDir(ctx, rootPath, project)
	// go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectDir
	err := cmd.Run()
	if err != nil {
		fmt.Printf("error running go mod tidy\n")
	}
}

func ProjectDir(ctx context.Context, rootPath string, project entity.Project) string {
	return path.Join(rootPath, project.Identifier)
}

func createDir(fullDir string) {
	err := os.MkdirAll(fullDir, 0o755)
	if err != nil {
		fmt.Printf("failed to create directory: %v", err)
	}
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
