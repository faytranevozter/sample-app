package helpers

import (
	"net/url"
	"testing"
)

func TestGetLimitOffset(t *testing.T) {
	tests := []struct {
		name         string
		query        url.Values
		defaultLimit []int
		wantPage     int64
		wantLimit    int64
		wantOffset   int64
	}{
		{
			name:       "defaults when no params",
			query:      url.Values{},
			wantPage:   1,
			wantLimit:  10,
			wantOffset: 0,
		},
		{
			name:       "page 2 with default limit",
			query:      url.Values{"page": []string{"2"}},
			wantPage:   2,
			wantLimit:  10,
			wantOffset: 10,
		},
		{
			name:       "custom page and limit",
			query:      url.Values{"page": []string{"3"}, "limit": []string{"20"}},
			wantPage:   3,
			wantLimit:  20,
			wantOffset: 40,
		},
		{
			name:       "page 1 explicit",
			query:      url.Values{"page": []string{"1"}, "limit": []string{"5"}},
			wantPage:   1,
			wantLimit:  5,
			wantOffset: 0,
		},
		{
			name:         "custom default limit",
			query:        url.Values{},
			defaultLimit: []int{25},
			wantPage:     1,
			wantLimit:    25,
			wantOffset:   0,
		},
		{
			name:         "query limit overrides default limit",
			query:        url.Values{"limit": []string{"15"}},
			defaultLimit: []int{25},
			wantPage:     1,
			wantLimit:    15,
			wantOffset:   0,
		},
		{
			name:       "invalid page defaults to 1",
			query:      url.Values{"page": []string{"abc"}, "limit": []string{"10"}},
			wantPage:   1,
			wantLimit:  10,
			wantOffset: 0,
		},
		{
			name:       "invalid limit defaults to 10",
			query:      url.Values{"page": []string{"2"}, "limit": []string{"abc"}},
			wantPage:   2,
			wantLimit:  10,
			wantOffset: 10,
		},
		{
			name:       "zero page defaults to 1",
			query:      url.Values{"page": []string{"0"}},
			wantPage:   1,
			wantLimit:  10,
			wantOffset: 0,
		},
		{
			name:       "zero limit defaults to 10",
			query:      url.Values{"limit": []string{"0"}},
			wantPage:   1,
			wantLimit:  10,
			wantOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotLimit, gotOffset := GetLimitOffset(tt.query, tt.defaultLimit...)
			if gotPage != tt.wantPage {
				t.Errorf("page = %d, want %d", gotPage, tt.wantPage)
			}
			if gotLimit != tt.wantLimit {
				t.Errorf("limit = %d, want %d", gotLimit, tt.wantLimit)
			}
			if gotOffset != tt.wantOffset {
				t.Errorf("offset = %d, want %d", gotOffset, tt.wantOffset)
			}
		})
	}
}
