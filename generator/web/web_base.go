package web

import (
	"bytes"
	"context"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"
	"os/exec"
	"path"
	"text/template"
	"unicode/utf8"

	"github.com/disintegration/letteravatar"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/core"
	"github.com/maykel/gpg/generator/helpers"
)

func GenerateBaseWeb(ctx context.Context, rootPath string, project entity.Project) error {
	fmt.Printf("--[GPG] Generating web\n")
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	webDir := path.Join(projectDir, generator.WEB_DIR)

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(webDir, "package.json"),
		TemplateName:    path.Join("web", "package"),
		Data:            project,
		DisableGoFormat: true,
	})

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(webDir, "public", "index.html"),
		TemplateName:    path.Join("web", "public", "index"),
		Data:            project,
		DisableGoFormat: true,
	})

	generateIcon(project, webDir, 192)
	generateIcon(project, webDir, 512)

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(webDir, "public", "manifest.json"),
		TemplateName:    path.Join("web", "public", "manifest"),
		Data:            project,
		DisableGoFormat: true,
	})

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(webDir, "src", "index.js"),
		TemplateName:    path.Join("web", "src", "index"),
		Data:            project,
		DisableGoFormat: true,
		Funcs: template.FuncMap{
			"ToCamelCase": helpers.ToCamelCase,
		},
	})

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(webDir, "src", "client.js"),
		TemplateName:    path.Join("web", "src", "client"),
		Data:            project,
		DisableGoFormat: true,
		Funcs: template.FuncMap{
			"ToCamelCase": helpers.ToCamelCase,
		},
	})

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(webDir, "src", "FetchUser.js"),
		TemplateName:    path.Join("web", "src", "user"),
		Data:            project,
		DisableGoFormat: true,
		Funcs: template.FuncMap{
			"ToCamelCase": helpers.ToCamelCase,
		},
	})

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(webDir, "src", "Utils.js"),
		TemplateName:    path.Join("web", "src", "utils"),
		Data:            project,
		DisableGoFormat: true,
		Funcs: template.FuncMap{
			"ToCamelCase": helpers.ToCamelCase,
		},
	})

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(webDir, "src", "App.js"),
		TemplateName:    path.Join("web", "src", "app"),
		Data:            project,
		DisableGoFormat: true,
		Funcs: template.FuncMap{
			"ToCamelCase": helpers.ToCamelCase,
		},
	})

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(webDir, "src", "App.css"),
		TemplateName:    path.Join("web", "src", "app_css"),
		Data:            project,
		DisableGoFormat: true,
	})

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(webDir, "src", "bootstrap.min.css"),
		TemplateName:    path.Join("web", "src", "bootstrap"),
		Data:            project,
		DisableGoFormat: true,
	})

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(webDir, "src", "components", "Header.js"),
		TemplateName:    path.Join("web", "src", "components", "header"),
		Data:            project,
		DisableGoFormat: true,
	})

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(webDir, "src", "components", "GPGModal.js"),
		TemplateName:    path.Join("web", "src", "components", "modal"),
		Data:            project,
		DisableGoFormat: true,
	})

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(webDir, "src", "components", "SearchEntity.js"),
		TemplateName:    path.Join("web", "src", "components", "search_entity"),
		Data:            project,
		DisableGoFormat: true,
		Funcs: template.FuncMap{
			"ToCamelCase": helpers.ToCamelCase,
			"HasSearch": func(e entity.Entity) bool {
				for _, f := range e.Fields {
					if f.StorageConfig.Search {
						return true
					}
				}
				return false
			},
			"SearchFields": func(e entity.Entity) string {
				searchFields := core.ResolveSearchFields(e)
				res := "(item) => { var res = ''; "
				first := true
				for _, sf := range searchFields {
					if sf.Type != "uuid.UUID" {
						if first {
							res += "res = item." + sf.Identifier + ";"
							first = false
						} else {
							res += "res = res.concat(' - ', item." + sf.Identifier + ");"
						}
					}
				}
				res += "return res;}"
				return res
			},
		},
	})

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(webDir, "src", "pages", "Login.js"),
		TemplateName:    path.Join("web", "src", "pages", "page_login"),
		Data:            project,
		DisableGoFormat: true,
	})

	GenerateEntities(ctx, rootPath, project)

	fmt.Printf("--[NPM] Install \n")
	cmd := exec.Command("npm", "install")
	cmd.Dir = webDir
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		fmt.Println("npm install result: " + out.String())
	}

	return nil
}

func generateIcon(project entity.Project, webDir string, size int) {
	firstLetter, _ := utf8.DecodeRuneInString(project.Render.Name)
	img, err := letteravatar.Draw(size, firstLetter, &letteravatar.Options{
		Palette: []color.Color{
			color.RGBA{76, 182, 172, 255},
		},
	},
	)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(path.Join(webDir, "public", fmt.Sprintf("logo%d.png", size)))
	if err != nil {
		log.Fatal(err)
	}

	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
	}
}
