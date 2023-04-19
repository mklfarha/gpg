package core

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path"
	"text/template"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

func GenerateCoreRepository(ctx context.Context, rootPath string, project entity.Project) {
	fmt.Printf("--[GPG] Generating core repository\n")
	projectDir := generator.ProjectDir(ctx, rootPath, project)
	repoDir := path.Join(projectDir, generator.CORE_REPO_DIR)

	generator.GenerateFile(
		ctx,
		generator.FileRequest{
			OutputFile:   path.Join(repoDir, "sqlc.yaml"),
			TemplateName: path.Join("core", "repo", "repo_yaml"),
			Data: struct {
				ProjectName string
			}{
				ProjectName: project.Identifier,
			},
			DisableGoFormat: true,
		},
	)
	generateSQLSchemas(ctx, repoDir, project)
}

type RepoSchemaTemplate struct {
	Entities []RepoSchemaEntity
}

type RepoSchemaEntity struct {
	Name             string
	NameTitle        string
	PrimaryKey       string
	Fields           []RepoSchemaField
	Indexes          []RepoSchemaIndex
	Search           []RepoSchemaSearch
	SelectStatements []RepoSchemaSelectStatement
	CustomQueries    []entity.CustomQuery
}

type RepoSchemaField struct {
	Name     string
	Type     string
	Null     string
	HasComma bool
	Default  string
	Unique   string
}

type RepoSchemaIndex struct {
	Name      string
	FieldName string
	HasComma  bool
}

type RepoSchemaSearch struct {
	Name      string
	FieldName string
	IsLast    bool
}

type RepoSchemaSelectStatement struct {
	Name          string
	GraphName     string
	Fields        []RepoSchemaSelectStatementField
	IsPrimary     bool
	TimeFields    []field.Template
	SortSupported bool
}

type RepoSchemaSelectStatementField struct {
	Name   string
	Field  field.Template
	IsLast bool
}

func generateSQLSchemas(ctx context.Context, repoDir string, project entity.Project) {
	entities := make([]RepoSchemaEntity, len(project.Entities))
	for i, e := range project.Entities {
		fields, indexes, search := resolveSchemaFieldsAndIndexes(e)
		selects := ResolveSelectStatements(e)
		entityTemplate := RepoSchemaEntity{
			Name:             e.Identifier,
			NameTitle:        helpers.ToCamelCase(e.Identifier),
			PrimaryKey:       EntityPrimaryKey(e).Identifier,
			Fields:           fields,
			Indexes:          indexes,
			Search:           search,
			SelectStatements: selects,
			CustomQueries:    e.CustomQueries,
		}
		entities[i] = entityTemplate
	}
	tpl := RepoSchemaTemplate{
		Entities: entities,
	}
	fmt.Printf("----[GPG] Generating schema\n")
	sqlDir := path.Join(repoDir, generator.CORE_REPO_SQL_DIR)
	for _, e := range entities {
		generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(sqlDir, "schemas", fmt.Sprintf("%s.sql", e.Name)),
			TemplateName: path.Join("core", "repo", "repo_schema"),
			Data: struct {
				Entity RepoSchemaEntity
			}{
				Entity: e,
			},
			DisableGoFormat: true,
		})
	}

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(sqlDir, "schemas", ".skeema"),
		TemplateName:    path.Join("core", "repo", "repo_skeema"),
		Data:            project,
		DisableGoFormat: true,
	})

	fmt.Printf("----[GPG] Generating inserts\n")
	sqlQueriesDir := path.Join(repoDir, generator.CORE_REPO_SQL_QUERIES_DIR)
	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(sqlQueriesDir, "inserts.sql"),
		TemplateName:    path.Join("core", "repo", "repo_inserts"),
		Data:            tpl,
		DisableGoFormat: true,
	})
	fmt.Printf("----[GPG] Generating updates\n")
	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(sqlQueriesDir, "updates.sql"),
		TemplateName:    path.Join("core", "repo", "repo_updates"),
		Data:            tpl,
		DisableGoFormat: true,
	})

	fmt.Printf("----[GPG] Generating selects\n")
	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(sqlQueriesDir, "selects.sql"),
		TemplateName:    path.Join("core", "repo", "repo_selects"),
		Data:            tpl,
		DisableGoFormat: true,
		Funcs: template.FuncMap{
			"ToSQL": func(qc entity.QueryCondition) string {
				return qc.ToSQL()
			},
		},
	})

	fmt.Printf("----[GPG] SQLC Generate\n")
	cmd := exec.Command("go", "run", "github.com/kyleconroy/sqlc/cmd/sqlc", "generate")
	cmd.Dir = repoDir
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		fmt.Println("SQLC Result: " + out.String())
	}

	generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:   path.Join(repoDir, "repository.go"),
		TemplateName: path.Join("core", "repo", "repository"),
		Data: struct {
			ProjectName string
		}{
			ProjectName: project.Identifier,
		},
	})

	fmt.Printf("----[GPG][Skeema] Sync DB - Diff \n")
	cmd = exec.Command("go", "run", "github.com/skeema/skeema", "diff")
	cmd.Dir = path.Join(sqlDir, "schemas")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		fmt.Println("skeema Result: " + out.String())
	}

	fmt.Printf("----[GPG][Skeema] Sync DB - Push \n")
	cmd = exec.Command("go", "run", "github.com/skeema/skeema", "push", "--allow-unsafe")
	cmd.Dir = path.Join(sqlDir, "schemas")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		fmt.Println("skeema Result: " + out.String())
	}

}

func EntityPrimaryKey(e entity.Entity) entity.Field {
	for _, field := range e.Fields {
		if field.StorageConfig.PrimaryKey {
			return field
		}
	}
	return entity.Field{}
}

func resolveSchemaFieldsAndIndexes(e entity.Entity) ([]RepoSchemaField, []RepoSchemaIndex, []RepoSchemaSearch) {
	fields := []RepoSchemaField{}
	index := []RepoSchemaIndex{}
	search := []RepoSchemaSearch{}
	for _, f := range e.Fields {
		if f.Stored {
			ft := RepoSchemaField{
				Name: f.Identifier,
				Type: resolveRepoFieldType(f),
				Null: "NOT NULL",
			}

			if f.StorageConfig.Unique {
				ft.Unique = "UNIQUE"
			}

			switch f.Type {
			case entity.DateFieldType:
				ft.Default = "default '2022-02-02'"
			case entity.DateTimeFieldType:
				ft.Default = "default CURRENT_TIMESTAMP"
			}

			fields = append(fields, ft)

			if f.StorageConfig.Index && !f.StorageConfig.Unique {
				index = append(index, RepoSchemaIndex{
					Name:      f.Identifier,
					FieldName: f.Identifier,
				})
			}

			if f.StorageConfig.Search {
				search = append(search, RepoSchemaSearch{
					Name:      f.Identifier,
					FieldName: f.Identifier,
				})
			}
		}
	}

	for i := 0; i < len(index)-1; i++ {
		index[i].HasComma = true
	}
	for i := 0; i < len(fields)-1; i++ {
		fields[i].HasComma = true
	}

	if len(search) > 0 {
		search[len(search)-1].IsLast = true
	}

	return fields, index, search
}

func resolveRepoFieldType(f entity.Field) string {
	switch f.Type {

	case entity.UUIDFieldType:
		return "CHAR(36)"
	case entity.IntFieldType:
		return "INT"
	case entity.FloatFieldType:
		return "DOUBLE"
	case entity.BooleanFieldType:
		return "TINYINT(1)"
	case entity.StringFieldType:
		return "VARCHAR(255)"
	case entity.LargeStringFieldType:
		return "TEXT"
	case entity.JSONFieldType:
		return "JSON"
	case entity.OptionsSingleFieldType:
		return "INT"
	case entity.OptionsManyFieldType:
		return "JSON"
	case entity.DateFieldType:
		return "DATE"
	case entity.DateTimeFieldType:
		return "DATETIME"
	}
	return ""
}

func ResolveSearchFields(e entity.Entity) []field.Template {
	fields := []field.Template{}
	for _, f := range e.Fields {
		if f.StorageConfig.Search {
			fields = append(fields, field.ResolveFieldType(f, e, nil))
		}
	}
	return fields
}

func ResolveSelectStatements(e entity.Entity) []RepoSchemaSelectStatement {
	selects := []RepoSchemaSelectStatement{}
	primaryKey := EntityPrimaryKey(e)
	resolvedField := field.ResolveFieldType(primaryKey, e, nil)
	nameByID := fmt.Sprintf("%sBy%s", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(primaryKey.Identifier))
	selects = append(selects, RepoSchemaSelectStatement{
		Name:      nameByID,
		GraphName: nameByID,
		Fields: []RepoSchemaSelectStatementField{
			{
				Name:   primaryKey.Identifier,
				Field:  resolvedField,
				IsLast: true,
			},
		},
		IsPrimary:     true,
		SortSupported: false,
	})
	indexes := []string{}
	indexFields := map[string]entity.Field{}

	indexesIDs := []string{}
	indexIDsFields := map[string]entity.Field{}
	timeFields := []field.Template{}
	for _, f := range e.Fields {
		if f.Stored {
			if f.StorageConfig.Index &&
				f.Type != entity.DateTimeFieldType &&
				f.Type != entity.UUIDFieldType {
				indexes = append(indexes, f.Identifier)
				indexFields[f.Identifier] = f
			}

			if f.StorageConfig.Index &&
				f.Type == entity.UUIDFieldType {
				indexesIDs = append(indexesIDs, f.Identifier)
				indexIDsFields[f.Identifier] = f
			}

			if f.Type == entity.DateFieldType || f.Type == entity.DateTimeFieldType {
				resolvedField := field.ResolveFieldType(f, e, nil)
				timeFields = append(timeFields, resolvedField)
			}
		}
	}

	if len(indexes) == 0 {
		return selects
	}

	combinations := Combinations(indexes)
	for _, combination := range combinations {
		name := fmt.Sprintf("%sBy", helpers.ToCamelCase(e.Identifier))
		fields := []RepoSchemaSelectStatementField{}
		first := true
		for i, f := range combination {
			isLast := true
			if i < len(combination)-1 {
				isLast = false
			}
			resolvedField := field.ResolveFieldType(indexFields[f], e, nil)
			fields = append(fields, RepoSchemaSelectStatementField{
				Name:   f,
				Field:  resolvedField,
				IsLast: isLast,
			})
			if first {
				first = false
				name = fmt.Sprintf("%s%s", name, helpers.ToCamelCase(f))
			} else {
				name = fmt.Sprintf("%sAnd%s", name, helpers.ToCamelCase(f))
			}
		}

		selects = append(selects, RepoSchemaSelectStatement{
			Name:          name,
			GraphName:     name,
			Fields:        fields,
			TimeFields:    timeFields,
			SortSupported: true,
		})
	}

	for _, indexID := range indexesIDs {
		resolvedIDField := field.ResolveFieldType(indexIDsFields[indexID], e, nil)
		name := fmt.Sprintf("%sBy%s", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(resolvedIDField.Identifier))
		fields := []RepoSchemaSelectStatementField{
			{
				Name:   indexID,
				Field:  resolvedIDField,
				IsLast: true,
			},
		}
		selects = append(selects, RepoSchemaSelectStatement{
			Name:          name,
			GraphName:     name,
			Fields:        fields,
			TimeFields:    timeFields,
			SortSupported: true,
		})
		for _, combination := range combinations {
			name := fmt.Sprintf("%sBy%s", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(resolvedIDField.Identifier))
			fields := []RepoSchemaSelectStatementField{
				{
					Name:   indexID,
					Field:  resolvedIDField,
					IsLast: false,
				},
			}
			for i, f := range combination {
				isLast := true
				if i < len(combination)-1 {
					isLast = false
				}
				resolvedField := field.ResolveFieldType(indexFields[f], e, nil)
				fields = append(fields, RepoSchemaSelectStatementField{
					Name:   f,
					Field:  resolvedField,
					IsLast: isLast,
				})

				name = fmt.Sprintf("%sAnd%s", name, helpers.ToCamelCase(f))
			}

			selects = append(selects, RepoSchemaSelectStatement{
				Name:          name,
				GraphName:     name,
				Fields:        fields,
				TimeFields:    timeFields,
				SortSupported: true,
			})
		}
	}

	return selects

}

func Combinations(set []string) (subsets [][]string) {
	length := uint(len(set))

	// Go through all possible combinations of objects
	// from 1 (only first object in subset) to 2^length (all objects in subset)
	for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
		var subset []string

		for object := uint(0); object < length; object++ {
			// checks if object is contained in subset
			// by checking if bit 'object' is set in subsetBits
			if (subsetBits>>object)&1 == 1 {
				// add object to subset
				subset = append(subset, set[object])
			}
		}
		// add subset to subsets
		subsets = append(subsets, subset)
	}
	return subsets
}
