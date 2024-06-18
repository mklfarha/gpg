package helpers

import (
	"fmt"

	"github.com/maykel/gpg/entity"
)

func ResolveTags(f entity.Field) string {
	return fmt.Sprintf("`json:\"%s\"`", f.Identifier)
}

func ProjectHasUserEntity(project entity.Project) bool {
	for _, e := range project.Entities {
		if e.Identifier == "user" {
			return true
		}
	}
	return false
}

func EntityContainsOperation(ops []entity.Operation, op entity.Operation) bool {
	for _, o := range ops {
		if o == op {
			return true
		}
	}
	return false
}

func FieldFromProject(project entity.Project, entityIdentifier, fieldIdentifier string) (entity.Field, entity.Entity, bool) {
	for _, e := range project.Entities {
		if e.Identifier == entityIdentifier {
			for _, f := range e.Fields {
				if f.Identifier == fieldIdentifier {
					return f, e, true
				}
			}
		}
	}
	return entity.Field{}, entity.Entity{}, false
}

func EntityPrimaryKey(e entity.Entity) entity.Field {
	for _, field := range e.Fields {
		if field.StorageConfig.PrimaryKey {
			return field
		}
	}
	return entity.Field{}
}

func FieldsToCamelCase(entities []entity.Entity) map[string]string {
	res := map[string]string{}
	for _, e := range entities {
		for _, f := range e.Fields {
			_, found := res[f.Identifier]
			if !found {
				res[f.Identifier] = ToCamelCase(f.Identifier)
			}
		}
	}
	return res
}
