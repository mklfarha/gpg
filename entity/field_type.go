package entity

import "encoding/json"

type FieldType int

const (
	InvalideFieldType FieldType = iota
	UUIDFieldType
	IntFieldType
	FloatFieldType
	BooleanFieldType
	StringFieldType
	LargeStringFieldType
	ArrayFieldType
	JSONFieldType
	OptionsSingleFieldType
	OptionsManyFieldType
	DateFieldType
	DateTimeFieldType
)

var UsesRandomValues = []FieldType{
	UUIDFieldType,
	IntFieldType,
	FloatFieldType,
	BooleanFieldType,
	StringFieldType,
	LargeStringFieldType,
	OptionsSingleFieldType,
	OptionsManyFieldType,
	DateFieldType,
	DateTimeFieldType,
}

func FieldTypeFromString(in string) FieldType {
	switch in {
	case "uuid":
		return UUIDFieldType
	case "int":
		return IntFieldType
	case "float":
		return FloatFieldType
	case "boolean":
		return BooleanFieldType
	case "string":
		return StringFieldType
	case "large_string":
		return LargeStringFieldType
	case "array":
		return ArrayFieldType
	case "json":
		return JSONFieldType
	case "options_single":
		return OptionsSingleFieldType
	case "options_many":
		return OptionsManyFieldType
	case "date":
		return DateFieldType
	case "datetime":
		return DateTimeFieldType
	}
	return InvalideFieldType
}

func (ft FieldType) String() string {
	switch ft {
	case UUIDFieldType:
		return "uuid"
	case IntFieldType:
		return "int"
	case FloatFieldType:
		return "float"
	case BooleanFieldType:
		return "boolean"
	case StringFieldType:
		return "string"
	case LargeStringFieldType:
		return "large_string"
	case ArrayFieldType:
		return "array"
	case JSONFieldType:
		return "json"
	case OptionsSingleFieldType:
		return "options_single"
	case OptionsManyFieldType:
		return "options_many"
	case DateFieldType:
		return "date"
	case DateTimeFieldType:
		return "datetime"
	}
	return "invalid"
}

func (ft *FieldType) UnmarshalJSON(data []byte) error {
	var item interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case string:
		*ft = FieldTypeFromString(string(v))
	}
	return nil
}

func (ft *FieldType) MarshalJSON() ([]byte, error) {
	return json.Marshal(ft.String())
}
