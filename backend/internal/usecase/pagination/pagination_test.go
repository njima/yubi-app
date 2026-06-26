package pagination

import "testing"

func intPtr(v int) *int { return &v }

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		page      *int
		limit     *int
		wantPage  int
		wantLimit int
	}{
		{
			name:      "nil params use defaults",
			page:      nil,
			limit:     nil,
			wantPage:  1,
			wantLimit: DefaultLimit,
		},
		{
			name:      "zero page uses default",
			page:      intPtr(0),
			limit:     nil,
			wantPage:  1,
			wantLimit: DefaultLimit,
		},
		{
			name:      "negative page uses default",
			page:      intPtr(-1),
			limit:     nil,
			wantPage:  1,
			wantLimit: DefaultLimit,
		},
		{
			name:      "valid page is used",
			page:      intPtr(3),
			limit:     nil,
			wantPage:  3,
			wantLimit: DefaultLimit,
		},
		{
			name:      "zero limit uses default",
			page:      nil,
			limit:     intPtr(0),
			wantPage:  1,
			wantLimit: DefaultLimit,
		},
		{
			name:      "negative limit uses default",
			page:      nil,
			limit:     intPtr(-5),
			wantPage:  1,
			wantLimit: DefaultLimit,
		},
		{
			name:      "valid limit is used",
			page:      nil,
			limit:     intPtr(50),
			wantPage:  1,
			wantLimit: 50,
		},
		{
			name:      "limit exceeding max is capped",
			page:      nil,
			limit:     intPtr(2000),
			wantPage:  1,
			wantLimit: MaxLimit,
		},
		{
			name:      "limit at max boundary is kept",
			page:      nil,
			limit:     intPtr(MaxLimit),
			wantPage:  1,
			wantLimit: MaxLimit,
		},
		{
			name:      "valid page and limit",
			page:      intPtr(5),
			limit:     intPtr(100),
			wantPage:  5,
			wantLimit: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.page, tt.limit)
			if got.Page != tt.wantPage {
				t.Errorf("Page = %d, want %d", got.Page, tt.wantPage)
			}
			if got.Limit != tt.wantLimit {
				t.Errorf("Limit = %d, want %d", got.Limit, tt.wantLimit)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		limit      int
		wantPage   int
		wantLimit  int
		wantOffset int
	}{
		{
			name:       "valid page and limit",
			page:       3,
			limit:      10,
			wantPage:   3,
			wantLimit:  10,
			wantOffset: 20,
		},
		{
			name:       "zero page uses first page",
			page:       0,
			limit:      50,
			wantPage:   1,
			wantLimit:  50,
			wantOffset: 0,
		},
		{
			name:       "negative page uses first page",
			page:       -2,
			limit:      50,
			wantPage:   1,
			wantLimit:  50,
			wantOffset: 0,
		},
		{
			name:       "zero limit uses default",
			page:       2,
			limit:      0,
			wantPage:   2,
			wantLimit:  DefaultLimit,
			wantOffset: DefaultLimit,
		},
		{
			name:       "limit exceeding max is capped",
			page:       2,
			limit:      MaxLimit + 1,
			wantPage:   2,
			wantLimit:  MaxLimit,
			wantOffset: MaxLimit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.page, tt.limit)

			if got.Page != tt.wantPage {
				t.Errorf("Page = %d, want %d", got.Page, tt.wantPage)
			}
			if got.Limit != tt.wantLimit {
				t.Errorf("Limit = %d, want %d", got.Limit, tt.wantLimit)
			}
			if got.Offset != tt.wantOffset {
				t.Errorf("Offset = %d, want %d", got.Offset, tt.wantOffset)
			}
		})
	}
}
