package cache

import (
	"io"
	"os"
)

// ReadFile reads a file and returns the contents as a string.
// If an error occurs while opening or reading the file, it returns "error".
func ReadFile(filePath string) string {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "error"
	}
	defer file.Close()

	// Read the file contents
	contents, err := io.ReadAll(file)
	if err != nil {
		return "error"
	}

	// Return the contents as a string
	return string(contents)
}

// BustCssCache reads the contents of the css.txt file and returns it as a string.
// It uses the ReadFile function to read the file.
func BustCssCache() string {
	return ReadFile("web/assets/bust/css.txt")
}

// BustJsCache reads the contents of the js.txt file and returns it as a string.
// It uses the ReadFile function to read the file.
func BustJsCache() string {
	return ReadFile("web/assets/bust/js.txt")
}
