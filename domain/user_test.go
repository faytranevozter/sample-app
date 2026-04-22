package domain

import (
	"testing"
)

func TestUserAllowedSort(t *testing.T) {
	expected := []string{"name", "email", "created_at", "updated_at"}
	if len(UserAllowedSort) != len(expected) {
		t.Errorf("UserAllowedSort length = %d, want %d", len(UserAllowedSort), len(expected))
	}
	for i, v := range expected {
		if UserAllowedSort[i] != v {
			t.Errorf("UserAllowedSort[%d] = %q, want %q", i, UserAllowedSort[i], v)
		}
	}
}
