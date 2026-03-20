package service

import (
	"testing"
)

func TestStrictSanitization(t *testing.T) {
	t.Run("Sanitize HTML Input", func(t *testing.T) {
		input := "<script>alert('xss')</script>"
		expected := ""
		result := strictSanitization(input)
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("Sanitize Safe Input", func(t *testing.T) {
		input := "This is a safe string."
		expected := "This is a safe string."
		result := strictSanitization(input)
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("Sanitize Input with HTML Tags", func(t *testing.T) {
		input := "<b>Bold Text</b> and <i>Italic Text</i>"
		expected := "Bold Text and Italic Text"
		result := strictSanitization(input)
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("Sanitize Empty Input", func(t *testing.T) {
		input := ""
		expected := ""
		result := strictSanitization(input)
		if result != expected {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})
}
