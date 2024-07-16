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
	PROTO_DIR                 = "idl"
	CORE_DIR                  = "core"
	CORE_ENTITY_DIR           = path.Join(CORE_DIR, "entity")
	CORE_MODULE_DIR           = path.Join(CORE_DIR, "module")
	CORE_REPO_DIR             = path.Join(CORE_DIR, "repository")
	EVENTS_REPO_DIR           = path.Join(CORE_DIR, "events")
	CORE_REPO_SQL_DIR         = "sql"
	CORE_REPO_SQL_QUERIES_DIR = "queries"
	CORE_TOOLS                = "tools"
	WEB_DIR                   = "web"
	CUSTOM_DIR                = "custom"
	CLI_DIR                   = "cli"
	MONITORING_DIR            = "monitoring"
)

var projectStructure = []string{
	API_DIR,
	AUTH_DIR,
	CONFIG_DIR,
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
		CreateDir(fullDir)
	}

	initModule(ctx, rootPath, project)

	return nil
}

func initModule(ctx context.Context, rootPath string, project entity.Project) {
	fmt.Printf("--[GPG] Init go module \n")
	projectDir := ProjectDir(ctx, rootPath, project)
	var out bytes.Buffer
	var stderr bytes.Buffer
	if !fileExists(path.Join(projectDir, "go.mod")) {
		cmd := exec.Command("go", "mod", "init", project.Identifier)
		cmd.Dir = projectDir
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("error running go mod init: %v | %v | %v\n", err, out.String(), stderr.String())
		}
	} else {
		fmt.Printf("----[GPG] go.mod already exists\n")
	}

	GenerateFile(ctx, FileRequest{
		OutputFile:   path.Join(projectDir, "tools", "tools.go"),
		TemplateName: "tools",
	})

	// go version

	cmd := exec.Command("go", "version")
	cmd.Dir = projectDir
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Run()
	fmt.Printf("--[GPG] go version: %v\n", out.String())

	// install sqlc
	cmd = exec.Command("go", "get", "github.com/sqlc-dev/sqlc/cmd/sqlc")
	cmd.Dir = projectDir
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("error running go install sqlc: %v | %v | %v\n", err, out.String(), stderr.String())
	}

	// install sqlc
	cmd = exec.Command("go", "install", "golang.org/x/sync/errgroup")
	cmd.Dir = projectDir
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("error running go get errgroup: %v | %v | %v\n", err, out.String(), stderr.String())
	}

	// install gqlgen
	cmd = exec.Command("go", "install", "github.com/99designs/gqlgen")
	cmd.Dir = projectDir
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("error running go get gqlgen:%v | %v | %v\n", err, out.String(), stderr.String())
	}

	// go mod tidy
	GoModTidy(ctx, rootPath, project)

}

func GoModTidy(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG] Go Mod Tidy \n")
	projectDir := ProjectDir(ctx, rootPath, project)

	var out bytes.Buffer
	var stderr bytes.Buffer
	// go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectDir
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("error running go mod tidy: %v | %v | %v\n", err, out.String(), stderr.String())
		return err
	}
	return nil
}

func ProjectDir(ctx context.Context, rootPath string, project entity.Project) string {
	return path.Join(rootPath, project.Identifier)
}

func CreateDir(fullDir string) {
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
