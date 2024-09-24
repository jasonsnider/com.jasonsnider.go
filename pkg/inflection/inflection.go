package inflection

import (
	"regexp"
)

func Humanize(input string) string {
	// One regex to replace underscores, hyphens, and handle CamelCase
	re := regexp.MustCompile(`([a-z])([A-Z])|[_-]`)
	humanized := re.ReplaceAllString(input, `$1 $2`)

	// Convert the entire string to lowercase
	return humanized
}
