package entity

import (
	"fmt"
	"strings"
)

type CustomQuery struct {
	Name      string         `json:"name"`
	ExtraFrom string         `json:"extra_from"`
	Joins     []JoinQuery    `json:"joins"`
	Condition QueryCondition `json:"condition"`
	Order     string         `json:"order"`
	Group     string         `json:"group"`
}

type QueryCondition struct {
	Operator    QueryOperator     `json:"operator"`
	Comparisons []QueryComparison `json:"comparisons"`
	Conditions  []QueryCondition  `json:"conditions"`
}

type QueryComparison struct {
	Comparator QueryComparator `json:"comparator"`
	FieldOne   Field           `json:"field_one"`
	FieldTwo   Field           `json:"field_two"`
}

type JoinQuery struct {
	EntityIdentifier string         `json:"entity_identifier"`
	Condition        QueryCondition `json:"condition"`
}

func (qc QueryCondition) ToSQL() string {

	if qc.Operator == InvalidQueryOperator {
		return ""
	}

	if len(qc.Conditions) == 0 && len(qc.Comparisons) == 1 {
		return qc.Comparisons[0].ToSQL()
	}

	result := ""
	for index, comp := range qc.Comparisons {
		result += comp.ToSQL()
		if index < len(qc.Comparisons)-1 {
			result += fmt.Sprintf(" %s ", qc.Operator.String())
		}
	}

	if len(qc.Conditions) == 0 {
		return result
	}

	renderedSubConditions := []string{}
	for _, cond := range qc.Conditions {
		renderedSubConditions = append(renderedSubConditions, cond.ToSQL())
	}

	return strings.Join(renderedSubConditions, fmt.Sprintf(" %s ", qc.Operator.String()))
}

func (comp QueryComparison) ToSQL() string {
	field1 := "?"
	if !comp.FieldOne.InputField {
		field1 = fmt.Sprintf("%s.%s", comp.FieldOne.ParentIdentifier, comp.FieldOne.Identifier)
	}
	field2 := "?"
	if !comp.FieldTwo.InputField {
		field2 = fmt.Sprintf("%s.%s", comp.FieldTwo.ParentIdentifier, comp.FieldTwo.Identifier)
	}
	return fmt.Sprintf("(%s %s %s)", field1, comp.Comparator.ToSQL(), field2)
}
