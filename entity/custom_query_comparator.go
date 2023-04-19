package entity

import "encoding/json"

type QueryComparator int

const (
	InvalidQueryComparator QueryComparator = iota
	EqualQueryComparator
	GreaterQueryComparator
	LesserQueryComparator
	LikeQueryComparator
	NotQueryComparator
)

func QueryComparatorFromString(in string) QueryComparator {
	switch in {
	case "EQUAL":
		return EqualQueryComparator
	case "GREATER":
		return GreaterQueryComparator
	case "LESSER":
		return LesserQueryComparator
	case "LIKE":
		return LikeQueryComparator
	case "NOT":
		return NotQueryComparator
	}
	return InvalidQueryComparator
}

func (ft QueryComparator) String() string {
	switch ft {
	case EqualQueryComparator:
		return "EQUAL"
	case GreaterQueryComparator:
		return "GREATER"
	case LesserQueryComparator:
		return "LESSER"
	case LikeQueryComparator:
		return "LIKE"
	case NotQueryComparator:
		return "NOT"
	}
	return "invalid"
}

func (ft *QueryComparator) UnmarshalJSON(data []byte) error {
	var item interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case string:
		*ft = QueryComparatorFromString(string(v))
	}
	return nil
}

func (ft *QueryComparator) MarshalJSON() ([]byte, error) {
	return json.Marshal(ft.String())
}

func (ft QueryComparator) ToSQL() string {
	switch ft {
	case EqualQueryComparator:
		return "="
	case GreaterQueryComparator:
		return ">"
	case LesserQueryComparator:
		return "<"
	case LikeQueryComparator:
		return "like"
	case NotQueryComparator:
		return "is not"
	}
	return "invalid"
}
