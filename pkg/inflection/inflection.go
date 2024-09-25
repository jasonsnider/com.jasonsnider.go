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

func Slugify(input string) string {
	// One regex to replace spaces, hyphens, and underscores
	re := regexp.MustCompile(`[\s_-]`)
	slug := re.ReplaceAllString(input, `-`)

	// Remove all non-alphanumeric characters except hyphens
	reNonAlphaNum := regexp.MustCompile(`[^a-zA-Z0-9-]`)
	slug = reNonAlphaNum.ReplaceAllString(slug, ``)

	// Convert the entire string to lowercase
	return slug
}
