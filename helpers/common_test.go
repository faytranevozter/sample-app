package helpers

import (
	"strings"
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid simple email", "test@example.com", true},
		{"valid with subdomain", "user@mail.example.com", true},
		{"valid with plus", "user+tag@example.com", true},
		{"valid with dots", "first.last@example.com", true},
		{"empty string", "", false},
		{"missing @", "testexample.com", false},
		{"missing domain", "test@", false},
		{"missing local", "@example.com", false},
		{"double @", "test@@example.com", false},
		{"spaces", "test @example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidEmail(tt.email)
			if got != tt.want {
				t.Errorf("IsValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestToJSON(t *testing.T) {
	tests := []struct {
		name   string
		input  any
		indent string
		want   string
	}{
		{
			name:   "simple map",
			input:  map[string]string{"key": "value"},
			indent: "\t",
			want:   "{\n\t\"key\": \"value\"\n}",
		},
		{
			name:   "integer",
			input:  42,
			indent: "  ",
			want:   "42",
		},
		{
			name:   "string",
			input:  "hello",
			indent: "\t",
			want:   "\"hello\"",
		},
		{
			name:   "nil",
			input:  nil,
			indent: "\t",
			want:   "null",
		},
		{
			name:   "empty indent",
			input:  map[string]int{"a": 1},
			indent: "",
			want:   "{\n\"a\": 1\n}",
		},
		{
			name:   "slice",
			input:  []int{1, 2, 3},
			indent: "\t",
			want:   "[\n\t1,\n\t2,\n\t3\n]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToJSON(tt.input, tt.indent)
			if got != tt.want {
				t.Errorf("ToJSON(%v, %q) = %q, want %q", tt.input, tt.indent, got, tt.want)
			}
		})
	}
}

func TestDump(t *testing.T) {
	// Dump writes to stdout; just verify it doesn't panic
	t.Run("single value", func(t *testing.T) {
		Dump("hello")
	})

	t.Run("multiple values", func(t *testing.T) {
		Dump("hello", 42, map[string]string{"key": "value"})
	})

	t.Run("nil value", func(t *testing.T) {
		Dump(nil)
	})
}

func TestToJSON_Struct(t *testing.T) {
	type sample struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	got := ToJSON(sample{Name: "test", Value: 1}, "\t")
	if !strings.Contains(got, `"name": "test"`) {
		t.Errorf("ToJSON struct should contain name field, got %q", got)
	}
	if !strings.Contains(got, `"value": 1`) {
		t.Errorf("ToJSON struct should contain value field, got %q", got)
	}
}
