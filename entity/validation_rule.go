package entity

import "encoding/json"

type ValidationRule int

const (
	InvalidValidationRule ValidationRule = iota
	NotNilValidationRule
	UniqueValidationRule
)

func ValidationRuleFromString(in string) ValidationRule {
	switch in {
	case "not_nil":
		return NotNilValidationRule
	case "unique":
		return UniqueValidationRule
	}
	return InvalidValidationRule
}

func (vt ValidationRule) String() string {
	switch vt {
	case NotNilValidationRule:
		return "not_nil"
	case UniqueValidationRule:
		return "unique"
	}

	return "invalid"
}

func (vt *ValidationRule) UnmarshalJSON(data []byte) error {
	var item interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case string:
		*vt = ValidationRuleFromString(string(v))
	}
	return nil
}

func (vt *ValidationRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(vt.String())
}
