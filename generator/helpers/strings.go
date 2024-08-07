package helpers

import (
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToCamelCase(str string) string {
	if str == "id" {
		return "ID"
	}
	if str == "uuid" {
		return "UUID"
	}

	key := strcase.ToCamel(str)
	if strings.Contains(key, "_Id") {
		key = strings.ReplaceAll(key, "Id", "ID")
	}

	key = strings.ReplaceAll(key, "uuid", "UUID")
	key = strings.ReplaceAll(key, "Uuid", "UUID")
	key = strings.ReplaceAll(key, "Json", "JSON")
	key = strings.ReplaceAll(key, "Url", "URL")
	key = strings.ReplaceAll(key, "Https", "HTTPS")
	key = strings.ReplaceAll(key, "Http", "HTTP")

	return key

	/*words := strings.Split(str, "_")
	key := strings.Title(strings.ToLower(words[0]))
	for _, word := range words[1:] {
		key += strings.Title(word)
	}

	if strings.Contains(key, "_Id") {
		key = strings.ReplaceAll(key, "Id", "ID")
	}

	key = strings.ReplaceAll(key, "uuid", "UUID")
	key = strings.ReplaceAll(key, "Uuid", "UUID")
	return key*/
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func StringContains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}
