package core

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/field"
	"github.com/maykel/gpg/generator/helpers"
)

func ResolveSelectStatements(project entity.Project, e entity.Entity) []RepoSchemaSelectStatement {
	selects := []RepoSchemaSelectStatement{}
	primaryKey := helpers.EntityPrimaryKey(e)
	resolvedField := field.ResolveFieldType(primaryKey, e, nil)
	nameByID := fmt.Sprintf("%sBy%s", helpers.ToCamelCase(e.Identifier), helpers.ToCamelCase(primaryKey.Identifier))
	selects = append(selects, RepoSchemaSelectStatement{
		Name:             nameByID,
		Identifier:       strcase.ToSnake(nameByID),
		EntityIdentifier: e.Identifier,
		GraphName:        nameByID,
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

	combinations := helpers.Combinations(indexes)
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
			Name:             name,
			Identifier:       strcase.ToSnake(name),
			EntityIdentifier: e.Identifier,
			GraphName:        name,
			Fields:           fields,
			TimeFields:       timeFields,
			SortSupported:    true,
		})
	}

	if project.Protocol == entity.ProjectProtocolProtobuf || project.DisableSelectCombinations {
		return selects
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
			Name:             name,
			Identifier:       strcase.ToSnake(name),
			EntityIdentifier: e.Identifier,
			GraphName:        name,
			Fields:           fields,
			TimeFields:       timeFields,
			SortSupported:    true,
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
				Name:             name,
				Identifier:       strcase.ToSnake(name),
				EntityIdentifier: e.Identifier,
				GraphName:        name,
				Fields:           fields,
				TimeFields:       timeFields,
				SortSupported:    true,
			})
		}
	}

	return selects

}
