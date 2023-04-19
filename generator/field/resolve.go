package field

import "github.com/maykel/gpg/entity"

func ResolveFieldType(f entity.Field, e entity.Entity, prefix *string) Template {
	switch f.Type {
	case entity.UUIDFieldType:
		return UUIDFieldTemplate(f, e)
	case entity.IntFieldType:
		return IntFieldTemplate(f, e)
	case entity.FloatFieldType:
		return FloatFieldTemplate(f, e)
	case entity.BooleanFieldType:
		return BooleanFieldTemplate(f, e)
	case entity.StringFieldType:
		return StringFieldTemplate(f, e)
	case entity.LargeStringFieldType:
		return StringFieldTemplate(f, e)
	case entity.JSONFieldType:
		return JSONFieldTemplate(f, e, prefix)
	case entity.OptionsSingleFieldType:
		return OptionsSingleFieldTemplate(f, e, prefix)
	case entity.OptionsManyFieldType:
		return OptionsManyFieldTemplate(f, e, prefix)
	case entity.DateFieldType:
		return DateFieldTemplate(f, e)
	case entity.DateTimeFieldType:
		return DatetimeFieldTemplate(f, e)
	}
	return Template{}
}

func ResolveFieldsAndImports(fields []entity.Field, e entity.Entity, prefix *string) ([]Template, map[string]any) {
	fieldTemplates := make([]Template, len(fields))
	imports := map[string]any{}
	for i, f := range fields {
		fieldTemplate := ResolveFieldType(f, e, prefix)
		if fieldTemplate.Import != nil {
			imports[*fieldTemplate.Import] = struct{}{}
		}
		fieldTemplates[i] = fieldTemplate
	}
	return fieldTemplates, imports
}
