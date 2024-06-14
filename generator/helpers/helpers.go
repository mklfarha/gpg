package helpers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/maykel/gpg/entity"
)

func ResolveTags(f entity.Field) string {
	return fmt.Sprintf("`json:\"%s\"`", f.Identifier)
}

func ToCamelCase(str string) string {
	if str == "id" {
		return "ID"
	}
	if str == "uuid" {
		return "UUID"
	}

	words := strings.Split(str, "_")
	key := strings.Title(strings.ToLower(words[0]))
	for _, word := range words[1:] {
		key += strings.Title(word)
	}

	if strings.Contains(key, "_Id") {
		key = strings.ReplaceAll(key, "Id", "ID")
	}

	key = strings.ReplaceAll(key, "uuid", "UUID")
	key = strings.ReplaceAll(key, "Uuid", "UUID")
	return key
}

func MapKeys(in map[string]any) []string {
	keys := make([]string, len(in))
	i := 0
	for k := range in {
		keys[i] = k
		i++
	}
	return keys
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ProjectHasUserEntity(project entity.Project) bool {
	for _, e := range project.Entities {
		if e.Identifier == "user" {
			return true
		}
	}
	return false
}

func StringContains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}

func EntityContainsOperation(ops []entity.Operation, op entity.Operation) bool {
	for _, o := range ops {
		if o == op {
			return true
		}
	}
	return false
}

func Inc(i int) int {
	return i + 1
}

func ProtoEnumOptions(protoType string, opts []entity.OptionValue) []string {
	res := []string{strcase.ToScreamingSnake(fmt.Sprintf("%s_%s", protoType, "invalid"))}
	for _, opt := range opts {
		res = append(res, strcase.ToScreamingSnake(fmt.Sprintf("%s_%s", protoType, opt.Identifier)))
	}
	return res
}
