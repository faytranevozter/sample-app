package domain

import (
	"testing"
)

func TestArticleFilter_Query_Empty(t *testing.T) {
	f := ArticleFilter{}
	query := make(map[string]any)
	result := f.Query(query)

	if result == nil {
		t.Error("Query should not return nil")
	}
}

func TestArticleFilter_Query_WithAuthorName(t *testing.T) {
	authorName := "John Doe"
	f := ArticleFilter{
		AuthorName: &authorName,
	}
	query := make(map[string]any)
	result := f.Query(query)

	val, ok := result["author.name"]
	if !ok {
		t.Error("author.name should be set in query")
	}
	if *(val.(*string)) != "John Doe" {
		t.Errorf("author.name = %v, want %q", val, "John Doe")
	}
}

func TestArticleFilter_Query_PreservesExistingQuery(t *testing.T) {
	f := ArticleFilter{}
	query := map[string]any{
		"deleted_at": map[string]any{"$eq": nil},
	}
	result := f.Query(query)

	if _, ok := result["deleted_at"]; !ok {
		t.Error("existing deleted_at query should be preserved")
	}
}

func TestArticleAllowedSort(t *testing.T) {
	expected := []string{"title", "content", "author.name", "created_at", "updated_at"}
	if len(ArticleAllowedSort) != len(expected) {
		t.Errorf("ArticleAllowedSort length = %d, want %d", len(ArticleAllowedSort), len(expected))
	}
	for i, v := range expected {
		if ArticleAllowedSort[i] != v {
			t.Errorf("ArticleAllowedSort[%d] = %q, want %q", i, ArticleAllowedSort[i], v)
		}
	}
}
