package util

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	patternAlphaNumeric = regexp.MustCompile("[^a-zA-Z0-9]+")
	patternLowerUpper   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func Slug(seperator string, input ...string) string {
	result := make([]string, len(input))
	for i, s := range input {
		if s == "" {
			continue
		}

		s := patternAlphaNumeric.ReplaceAllString(s, seperator)
		s = patternLowerUpper.ReplaceAllString(s, fmt.Sprintf("$1%s$2", seperator))
		result[i] = s
	}

	return strings.Join(result, seperator)
}
