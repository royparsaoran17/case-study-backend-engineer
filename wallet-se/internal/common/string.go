package common

import (
	"regexp"
	"strings"
)

var space = regexp.MustCompile(`\s+`)

func ToSnakeCase(str string) string {
	snake := space.ReplaceAllString(str, " ")
	snake = strings.ReplaceAll(snake, " ", "_")
	return strings.ToLower(snake)
}
