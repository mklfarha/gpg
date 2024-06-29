package field

import (
	"github.com/maykel/gpg/entity"
)

func ResolveFieldType(f entity.Field, rootEntity entity.Entity, dependantEntity *Template) Template {
	switch f.Type {
	case entity.UUIDFieldType:
		return UUIDFieldTemplate(f, rootEntity)
	case entity.IntFieldType:
		return IntFieldTemplate(f, rootEntity)
	case entity.FloatFieldType:
		return FloatFieldTemplate(f, rootEntity)
	case entity.BooleanFieldType:
		return BooleanFieldTemplate(f, rootEntity)
	case entity.StringFieldType:
		return StringFieldTemplate(f, rootEntity)
	case entity.LargeStringFieldType:
		return StringFieldTemplate(f, rootEntity)
	case entity.JSONFieldType:
		dependant := false
		if dependantEntity != nil {
			dependant = true
		}
		return JSONFieldTemplate(f, rootEntity, dependant)
	case entity.OptionsSingleFieldType:
		return OptionsSingleFieldTemplate(f, rootEntity, dependantEntity)
	case entity.OptionsManyFieldType:
		return OptionsManyFieldTemplate(f, rootEntity, dependantEntity)
	case entity.DateFieldType:
		return DateFieldTemplate(f, rootEntity)
	case entity.DateTimeFieldType:
		return DatetimeFieldTemplate(f, rootEntity)
	}
	return Template{}
}

func ResolveFieldsAndImports(fields []entity.Field, e entity.Entity, dependantEntity *Template) ([]Template, map[string]any) {
	fieldTemplates := make([]Template, len(fields))
	imports := map[string]any{}
	for i, f := range fields {
		fieldTemplate := ResolveFieldType(f, e, dependantEntity)
		if fieldTemplate.Import != nil {
			imports[*fieldTemplate.Import] = struct{}{}
		}
		fieldTemplates[i] = fieldTemplate
	}
	return fieldTemplates, imports
}
