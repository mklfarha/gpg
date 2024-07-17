package repo

import (
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/field"
)

type SchemaTemplate struct {
	Entities []SchemaEntity
}

type SchemaEntity struct {
	Name             string
	NameTitle        string
	PrimaryKey       string
	Fields           []SchemaField
	Indexes          []SchemaIndex
	Search           []SchemaSearch
	SelectStatements []SchemaSelectStatement
	CustomQueries    []entity.CustomQuery
}

type SchemaField struct {
	Name     string
	Type     string
	Null     string
	HasComma bool
	Default  string
	Unique   string
}

type SchemaIndex struct {
	Name      string
	FieldName string
	HasComma  bool
}

type SchemaSearch struct {
	Name      string
	FieldName string
	IsLast    bool
}

type SchemaSelectStatement struct {
	Name             string
	Identifier       string
	EntityIdentifier string
	GraphName        string
	Fields           []SchemaSelectStatementField
	IsPrimary        bool
	TimeFields       []field.Template
	SortSupported    bool
}

type SchemaSelectStatementField struct {
	Name   string
	Field  field.Template
	IsLast bool
}

func resolveRepoFieldType(f entity.Field) string {
	switch f.Type {

	case entity.UUIDFieldType:
		return "CHAR(36)"
	case entity.IntFieldType:
		return "INT"
	case entity.FloatFieldType:
		return "DOUBLE"
	case entity.BooleanFieldType:
		return "TINYINT(1)"
	case entity.StringFieldType:
		return "VARCHAR(255)"
	case entity.LargeStringFieldType:
		return "TEXT"
	case entity.JSONFieldType, entity.ArrayFieldType:
		return "JSON"
	case entity.OptionsSingleFieldType:
		return "INT"
	case entity.OptionsManyFieldType:
		return "JSON"
	case entity.DateFieldType:
		return "DATE"
	case entity.DateTimeFieldType:
		return "DATETIME"
	}
	return ""
}
