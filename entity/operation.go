package entity

import "encoding/json"

type Operation int

const (
	InvalidOperation Operation = iota
	SelectOperation
	UpsertOperation
	DeleteOperation
)

func OperationFromString(in string) Operation {
	switch in {
	case "select":
		return SelectOperation
	case "upsert":
		return UpsertOperation
	case "delete":
		return DeleteOperation
	}
	return InvalidOperation
}

func (o Operation) String() string {
	switch o {
	case SelectOperation:
		return "select"
	case UpsertOperation:
		return "upsert"
	case DeleteOperation:
		return "delete"
	}
	return "invalid"
}

func (o *Operation) UnmarshalJSON(data []byte) error {
	var item interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case string:
		*o = OperationFromString(string(v))
	}
	return nil
}

func (o *Operation) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}
