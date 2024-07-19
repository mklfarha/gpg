package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/files"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/auth"
	gcli "github.com/maykel/gpg/generator/cli"
	"github.com/maykel/gpg/generator/core"
	"github.com/maykel/gpg/generator/core/events"
	"github.com/maykel/gpg/generator/core/repo"
	"github.com/maykel/gpg/generator/graph"
	"github.com/maykel/gpg/generator/monitoring"
	"github.com/maykel/gpg/generator/proto"
	"github.com/maykel/gpg/generator/web"
	"github.com/urfave/cli"
	"golang.org/x/sync/errgroup"
)

const (
	API_PROTOCOL_ALL      = "all"
	API_PROTOCOL_GRAPHQL  = "graphql"
	API_PROTOCOL_PROTOBUF = "protobuf"

	FLAG_PROTOCOL            = "protocol"
	FLAG_SELECT_COMBINATIONS = "enable_select_combinations"
	FLAG_SKIP_SKEEMA         = "skip_skeema"
)

var enableSelectCombinations bool
var skipSkeema bool

func main() {
	app := cli.NewApp()
	app.Name = "[GPG] Go Project Generator"

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        FLAG_SELECT_COMBINATIONS,
			Usage:       "Enable the generation of all select combination methods",
			Destination: &enableSelectCombinations,
		},
		&cli.BoolFlag{
			Name:        FLAG_SKIP_SKEEMA,
			Usage:       "Use flag to disable skeem sync with DB",
			Destination: &skipSkeema,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "genall", // Generates API/Auth and Web management tool
			Usage: "Provide a project config file and target directory",
			Action: func(c *cli.Context) error {
				configPath := c.Args().Get(0)
				targetDir := c.Args().Get(1)
				project, err := loadProject(configPath)
				if err != nil {
					panic(err)
				}
				generateAPI(targetDir, project)
				generateWeb(targetDir, project)
				return nil
			},
		},
		{
			Name:  "genweb", // Generates Web management tool
			Usage: "Provide a project config file and target directory",
			Action: func(c *cli.Context) error {
				configPath := c.Args().Get(0)
				targetDir := c.Args().Get(1)
				project, err := loadProject(configPath)
				if err != nil {
					panic(err)
				}

				generateWeb(targetDir, project)
				return nil
			},
		},
		{
			Name:  "genapi", // Generates API/Auth
			Usage: "Provide a project config file and target directory",
			Action: func(c *cli.Context) error {
				configPath := c.Args().Get(0)
				targetDir := c.Args().Get(1)
				project, err := loadProject(configPath)
				if err != nil {
					panic(err)
				}
				generateAPI(targetDir, project)
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func loadProject(configPath string) (entity.Project, error) {
	project := entity.Project{}

	// if empty try looking at the project level
	if configPath == "" {
		cliDir := files.AppDir()
		configPath = path.Join(cliDir, "gpg/project_config.json")
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println("reading error", err)
		return project, err
	}

	if err := json.Unmarshal(data, &project); err != nil {
		fmt.Println("unmarshal error", err)
	}
	res, _ := json.Marshal(project)
	fmt.Printf("%v \n", string(res))

	pl := pluralize.NewClient()
	for _, e := range project.Entities {
		for i, f := range e.Fields {
			f.ParentIdentifier = e.Identifier
			if f.Type == entity.JSONFieldType {
				if f.JSONConfig.Identifier == "" {
					e.Fields[i].JSONConfig.Identifier = pl.Singular(f.Identifier)
				}
			}
		}
	}

	fmt.Printf("--[GPG] Project Loaded \n")
	project.Identifier = strcase.ToSnake(project.Identifier)
	return project, nil
}

type APIGenerator struct {
	Name     string
	Func     func() error
	Blocking bool
}

func generateAPI(targetDir string, project entity.Project) {
	ctx := context.Background()

	project.DisableSelectCombinations = !enableSelectCombinations

	err := generator.GenerateProjectDirectories(ctx, targetDir, project)
	if err != nil {
		fmt.Sprintf("errors generating directories: %v", err)
		panic("error generating directories")
	}

	generators := []APIGenerator{
		{
			Name: "configuration",
			Func: func() error {
				return generator.GenerateConfig(ctx, targetDir, project)
			},
			Blocking: true,
		},
		{
			Name: "core entities",
			Func: func() error {
				return core.GenerateCoreEntities(ctx, targetDir, project)
			},
			Blocking: true,
		},
		{
			Name: "core repo",
			Func: func() error {
				return repo.GenerateCoreRepository(ctx, targetDir, project, skipSkeema)
			},
			Blocking: true,
		},
		{
			Name: "core events",
			Func: func() error {
				return events.GenerateCoreEvents(ctx, targetDir, project)
			},
			Blocking: true,
		},
		{
			Name: "core modules",
			Func: func() error {
				return core.GenerateCoreModules(ctx, targetDir, project)
			},
			Blocking: true,
		},
		{
			Name: "graphql",
			Func: func() error {
				fmt.Printf("protocoool: %v\n", project.API.Protocol)
				if project.API.Protocol == API_PROTOCOL_GRAPHQL || project.API.Protocol == API_PROTOCOL_ALL {
					return graph.GenerateGraph(ctx, targetDir, project)
				}
				return nil
			},
			Blocking: true,
		},
		{
			Name: "proto",
			Func: func() error {
				return proto.Generate(ctx, targetDir, project)
			},
			Blocking: true,
		},
		{
			Name: "monitoring",
			Func: func() error {
				return monitoring.GenerateMonitoring(ctx, targetDir, project)
			},
			Blocking: true,
		},
		{
			Name: "auth",
			Func: func() error {
				return auth.GenerateAuth(ctx, targetDir, project)
			},
			Blocking: true,
		},
		{
			Name: "api",
			Func: func() error {
				return generator.GenerateAPIModule(ctx, targetDir, project)
			},
			Blocking: true,
		},
		{
			Name: "custom",
			Func: func() error {
				return generator.GenerateCustom(ctx, targetDir, project)
			},
			Blocking: true,
		},
		{
			Name: "cli",
			Func: func() error {
				return gcli.GenerateCLIModule(ctx, targetDir, project)
			},
			Blocking: true,
		},
	}

	eg := errgroup.Group{}
	for _, g := range generators {
		eg.Go(g.Func)
	}

	if err := eg.Wait(); err != nil {
		log.Fatalf("error: %v", err)
	}

	err = generator.GoModTidy(ctx, targetDir, project)
	if err != nil {
		fmt.Printf("error running go mod tidy: %v", err)
	}

}

func generateWeb(targetDir string, project entity.Project) {
	ctx := context.Background()
	web.GenerateBaseWeb(ctx, targetDir, project)
}
