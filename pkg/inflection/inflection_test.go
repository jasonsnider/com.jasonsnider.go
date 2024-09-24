package inflection

import (
	"testing"
)

func TestHumanize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"helloWorld", "hello World"},
		{"hello_world", "hello world"},
		{"hello-world", "hello world"},
		{"HelloWorld", "Hello World"},
		{"helloWorld_test", "hello World test"},
		{"hello-World_test", "hello World test"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := Humanize(test.input)
			if result != test.expected {
				t.Errorf("humanize(%q) = %q; want %q", test.input, result, test.expected)
			}
		})
	}
}
