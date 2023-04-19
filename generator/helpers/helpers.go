package helpers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/maykel/gpg/entity"
)

func ResolveTags(f entity.Field) string {
	return fmt.Sprintf("`json:\"%s\"`", f.Identifier)
}

func ToCamelCase(str string) string {
	if str == "id" {
		return "ID"
	}

	words := strings.Split(str, "_")
	key := strings.Title(strings.ToLower(words[0]))
	for _, word := range words[1:] {
		key += strings.Title(word)
	}
	key = strings.ReplaceAll(key, "Id", "ID")
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
