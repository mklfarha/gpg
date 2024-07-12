package repo

import (
	"context"
	"fmt"
	"path"
	"text/template"

	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator"
	"github.com/maykel/gpg/generator/helpers"
)

func generateRepositorySQL(ctx context.Context, project entity.Project, sqlDir string) error {
	entities := make([]SchemaEntity, len(project.Entities))
	for i, e := range project.Entities {
		fields, indexes, search := resolveSchemaFieldsAndIndexes(e)
		selects := ResolveSelectStatements(project, e)
		entityTemplate := SchemaEntity{
			Name:             e.Identifier,
			NameTitle:        helpers.ToCamelCase(e.Identifier),
			PrimaryKey:       helpers.EntityPrimaryKey(e).Identifier,
			Fields:           fields,
			Indexes:          indexes,
			Search:           search,
			SelectStatements: selects,
			CustomQueries:    e.CustomQueries,
		}
		entities[i] = entityTemplate
	}
	tpl := SchemaTemplate{
		Entities: entities,
	}
	fmt.Printf("----[GPG] Generating SQL files\n")
	fmt.Printf("----[GPG] Generating SQL schema\n")
	for _, e := range entities {
		err := generator.GenerateFile(ctx, generator.FileRequest{
			OutputFile:   path.Join(sqlDir, "schemas", fmt.Sprintf("%s.sql", e.Name)),
			TemplateName: path.Join("core", "repo", "repo_schema"),
			Data: struct {
				Entity SchemaEntity
			}{
				Entity: e,
			},
			DisableGoFormat: true,
		})
		if err != nil {
			return err
		}
	}

	// queries directory
	sqlQueriesDir := path.Join(sqlDir, generator.CORE_REPO_SQL_QUERIES_DIR)

	fmt.Printf("----[GPG] Generating inserts\n")
	err := generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(sqlQueriesDir, "inserts.sql"),
		TemplateName:    path.Join("core", "repo", "repo_inserts"),
		Data:            tpl,
		DisableGoFormat: true,
	})
	if err != nil {
		return err
	}

	fmt.Printf("----[GPG] Generating updates\n")
	err = generator.GenerateFile(ctx, generator.FileRequest{
		OutputFile:      path.Join(sqlQueriesDir, "updates.sql"),
		TemplateName:    path.Join("core", "repo", "repo_updates"),
		Data:            tpl,
		DisableGoFormat: true,
	})
	if err != nil {
		return err
	}

	fmt.Printf("----[GPG] Generating selects\n")
	return generator.GenerateFile(ctx, generator.FileRequest{
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
}

func resolveSchemaFieldsAndIndexes(e entity.Entity) ([]SchemaField, []SchemaIndex, []SchemaSearch) {
	fields := []SchemaField{}
	index := []SchemaIndex{}
	search := []SchemaSearch{}
	for _, f := range e.Fields {
		if f.Stored {
			ft := SchemaField{
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
				index = append(index, SchemaIndex{
					Name:      f.Identifier,
					FieldName: f.Identifier,
				})
			}

			if f.StorageConfig.Search {
				search = append(search, SchemaSearch{
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
