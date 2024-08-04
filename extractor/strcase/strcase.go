package strcase

import (
	"regexp"
	"strings"
)

func ToCamel(varName string, initialCaseRegexp *regexp.Regexp) string {
	converted := initialCaseRegexp.ReplaceAllStringFunc(varName, func(s string) string {
		return strings.ToUpper(string(s[1]))
	})
	return strings.ToLower(string(converted[0])) + converted[1:]
}

var KebabRegexp, _ = regexp.Compile(`-\w`)
var SnakeRegexp, _ = regexp.Compile(`_\w`)
var SpaceRegexp, _ = regexp.Compile(` \w`)
