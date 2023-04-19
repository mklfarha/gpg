package entity

import "encoding/json"

type QueryOperator int

const (
	InvalidQueryOperator QueryOperator = iota
	ANDQueryOperator
	ORQueryOperator
)

func QueryOperatorFromString(in string) QueryOperator {
	switch in {
	case "AND":
		return ANDQueryOperator
	case "OR":
		return ORQueryOperator
	}
	return InvalidQueryOperator
}

func (ft QueryOperator) String() string {
	switch ft {
	case ANDQueryOperator:
		return "AND"
	case ORQueryOperator:
		return "OR"
	}
	return "invalid"
}

func (ft *QueryOperator) UnmarshalJSON(data []byte) error {
	var item interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case string:
		*ft = QueryOperatorFromString(string(v))
	}
	return nil
}

func (ft *QueryOperator) MarshalJSON() ([]byte, error) {
	return json.Marshal(ft.String())
}
