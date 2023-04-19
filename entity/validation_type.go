package entity

import "encoding/json"

type ValidationType int

const (
	InvalidValidationType ValidationType = iota
	RuleValidationType
)

func ValidationTypeFromString(in string) ValidationType {
	switch in {
	case "rule":
		return RuleValidationType
	}
	return InvalidValidationType
}

func (vt ValidationType) String() string {
	switch vt {
	case RuleValidationType:
		return "rule"
	}
	return "invalid"
}

func (vt *ValidationType) UnmarshalJSON(data []byte) error {
	var item interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case string:
		*vt = ValidationTypeFromString(string(v))
	}
	return nil
}

func (vt *ValidationType) MarshalJSON() ([]byte, error) {
	return json.Marshal(vt.String())
}
