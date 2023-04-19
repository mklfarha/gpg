package entity

import "encoding/json"

type JSONConfigType int

const (
	InvalidJSONConfigType JSONConfigType = iota
	ManyJSONConfigType
	SingleJSONConfigType
)

func JSONConfigTypeFromString(in string) JSONConfigType {
	switch in {
	case "many":
		return ManyJSONConfigType
	case "single":
		return SingleJSONConfigType
	}
	return InvalidJSONConfigType
}

func (t JSONConfigType) String() string {
	switch t {
	case ManyJSONConfigType:
		return "many"
	case SingleJSONConfigType:
		return "single"
	}

	return "invalid"
}

func (t *JSONConfigType) UnmarshalJSON(data []byte) error {
	var item interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case string:
		*t = JSONConfigTypeFromString(string(v))
	}
	return nil
}

func (t *JSONConfigType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}
