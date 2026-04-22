package middleware

import (
	"testing"
)

func TestMd5String(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty string",
			input: "",
			want:  "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name:  "hello",
			input: "hello",
			want:  "5d41402abc4b2a76b9719d911017c592",
		},
		{
			name:  "url-like string",
			input: "/api/v1/users?page=1&limit=10",
			want:  md5String("/api/v1/users?page=1&limit=10"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := md5String(tt.input)
			if got != tt.want {
				t.Errorf("md5String(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestMd5String_Consistency(t *testing.T) {
	input := "test-input-value"
	first := md5String(input)
	second := md5String(input)

	if first != second {
		t.Errorf("md5String should be consistent, got %q and %q", first, second)
	}
}

func TestMd5String_Uniqueness(t *testing.T) {
	hash1 := md5String("input1")
	hash2 := md5String("input2")

	if hash1 == hash2 {
		t.Error("different inputs should produce different hashes")
	}
}
