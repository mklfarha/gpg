package entity

type Field struct {
	Identifier       string        `json:"identifier"`
	ParentIdentifier string        `json:"parent_identifier"`
	Render           Render        `json:"render"`
	Type             FieldType     `json:"type"`
	EntityRef        string        `json:"entity_ref"`
	JSONConfig       JSONConfig    `json:"json_config,omitempty"`
	ArrayConfig      ArrayConfig   `json:"array_config,omitempty"`
	OptionValues     []OptionValue `json:"values,omitempty"`
	Deprecated       bool          `json:"deprecated"`
	Required         bool          `json:"required"`
	Stored           bool          `json:"stored"`
	StorageConfig    StorageConfig `json:"storage_config"`
	Autogenerated    Autogenerated `json:"autogenerated"`
	Hidden           Hidden        `json:"hidden"`
	Validations      []Validation  `json:"validations"`
	InputField       bool          `json:"input_field"`
}

type Autogenerated struct {
	Type        AutogeneratedType `json:"type"`
	FuncName    string            `json:"func_name"`
	FailOnError bool              `json:"fail_on_error"`
}

type Hidden struct {
	API   []Operation `json:"api"`
	Admin []Operation `json:"admin"`
}

type OptionValue struct {
	Identifier string `json:"identifier"`
	Display    string `json:"display"`
}

type StorageConfig struct {
	PrimaryKey bool `json:"primary_key"`
	Index      bool `json:"index"`
	Search     bool `json:"search"`
	Unique     bool `json:"unique"`
}

type ArrayConfig struct {
	Type      FieldType `json:"type"`
	EntityRef string    `json:"entity_ref"`
}

func (f Field) HasNestedJsonFields() bool {
	if f.Type != JSONFieldType || f.JSONConfig.Reuse {
		return false
	}
	hasNestedJSONField := false
	for _, jf := range f.JSONConfig.Fields {
		if jf.Type == JSONFieldType && !jf.JSONConfig.Reuse && len(jf.JSONConfig.Fields) > 0 {
			hasNestedJSONField = true
		}
	}
	return hasNestedJSONField
}

func (f Field) NestedJsonFields() []Field {
	if f.Type != JSONFieldType || f.JSONConfig.Reuse {
		return []Field{}
	}
	nestedJSONField := []Field{}
	for _, jf := range f.JSONConfig.Fields {
		if jf.Type == JSONFieldType && !jf.JSONConfig.Reuse && len(jf.JSONConfig.Fields) > 0 {
			nestedJSONField = append(nestedJSONField, jf)
		}
	}
	return nestedJSONField
}
