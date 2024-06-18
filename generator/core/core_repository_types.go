package core

import (
	"github.com/maykel/gpg/entity"
	"github.com/maykel/gpg/generator/field"
)

type RepoSchemaTemplate struct {
	Entities []RepoSchemaEntity
}

type RepoSchemaEntity struct {
	Name             string
	NameTitle        string
	PrimaryKey       string
	Fields           []RepoSchemaField
	Indexes          []RepoSchemaIndex
	Search           []RepoSchemaSearch
	SelectStatements []RepoSchemaSelectStatement
	CustomQueries    []entity.CustomQuery
}

type RepoSchemaField struct {
	Name     string
	Type     string
	Null     string
	HasComma bool
	Default  string
	Unique   string
}

type RepoSchemaIndex struct {
	Name      string
	FieldName string
	HasComma  bool
}

type RepoSchemaSearch struct {
	Name      string
	FieldName string
	IsLast    bool
}

type RepoSchemaSelectStatement struct {
	Name          string
	GraphName     string
	Fields        []RepoSchemaSelectStatementField
	IsPrimary     bool
	TimeFields    []field.Template
	SortSupported bool
}

type RepoSchemaSelectStatementField struct {
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
	case entity.JSONFieldType:
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
