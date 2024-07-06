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
	"strings"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/files"
	"github.com/maykel/gpg/generator"
	gcli "github.com/maykel/gpg/generator/cli"
	"github.com/maykel/gpg/generator/core"
	"github.com/maykel/gpg/generator/graph"
	"github.com/maykel/gpg/generator/helpers"
	"github.com/maykel/gpg/generator/monitoring"
	"github.com/maykel/gpg/generator/proto"
	"github.com/maykel/gpg/generator/web"
	"github.com/urfave/cli"
)

const (
	API_PROTOCOL_ALL      = "all"
	API_PROTOCOL_GRAPHQL  = "graphql"
	API_PROTOCOL_PROTOBUF = "protobuf"

	FLAG_PROTOCOL            = "protocol"
	FLAG_SELECT_COMBINATIONS = "enable_select_combinations"
	FLAG_SKIP_SKEEMA         = "skip_skeema"
)

var protocol string
var enableSelectCombinations bool
var skipSkeema bool

func main() {
	app := cli.NewApp()
	app.Name = "[GPG] Go Project Generator"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        FLAG_PROTOCOL,
			Value:       "all",
			Usage:       "API protocol to generate",
			Destination: &protocol,
		},
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

	for _, e := range project.Entities {
		for _, f := range e.Fields {
			f.ParentIdentifier = e.Identifier
		}
	}

	fmt.Printf("--[GPG] Project Loaded \n")
	project.Identifier = strings.ReplaceAll(project.Identifier, "-", "_")
	if helpers.ProjectHasUserEntity(project) {
		// in the future add a flag
		project.Auth.Enabled = true
	}
	return project, nil
}

func generateAPI(targetDir string, project entity.Project) {
	ctx := context.Background()
	switch protocol {
	case API_PROTOCOL_ALL:
		project.Protocol = entity.ProjectProtocolAll
	case API_PROTOCOL_GRAPHQL:
		project.Protocol = entity.ProjectProtocolGraphQL
	case API_PROTOCOL_PROTOBUF:
		project.Protocol = entity.ProjectProtocolProtobuf
	}

	project.DisableSelectCombinations = !enableSelectCombinations

	generator.GenerateProjectDirectories(ctx, targetDir, project)
	generator.GenerateConfig(ctx, targetDir, project)
	core.GenerateCoreEntities(ctx, targetDir, project)
	core.GenerateCoreRepository(ctx, targetDir, project, skipSkeema)
	err := core.GenerateCoreModules(ctx, targetDir, project)
	if err != nil {
		fmt.Printf("ERROR: Error generating core modules: %v", err)
	}
	switch protocol {
	case API_PROTOCOL_ALL:
		graph.GenerateGraph(ctx, targetDir, project)
		proto.Generate(ctx, targetDir, project)
	case API_PROTOCOL_GRAPHQL:
		graph.GenerateGraph(ctx, targetDir, project)
	case API_PROTOCOL_PROTOBUF:
		proto.Generate(ctx, targetDir, project)
	}
	monitoring.GenerateMonitoring(ctx, targetDir, project)
	generator.GenerateAuth(ctx, targetDir, project)
	generator.GenerateAPIModule(ctx, targetDir, project)
	generator.GenerateCustom(ctx, targetDir, project)
	generator.GoModTidy(context.Background(), targetDir, project)
	gcli.GenerateCLIModule(ctx, targetDir, project)
}

func generateWeb(targetDir string, project entity.Project) {
	ctx := context.Background()
	web.GenerateBaseWeb(ctx, targetDir, project)
}
