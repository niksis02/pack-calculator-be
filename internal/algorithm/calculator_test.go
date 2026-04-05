package algorithm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/niksis02/pack-calculator-be/internal/algorithm"
)

func TestCalculate(t *testing.T) {
	t.Parallel()

	defaultPacks := []int{250, 500, 1000, 2000, 5000}

	tests := []struct {
		name      string
		packs     []int
		order     int
		wantTotal int
		wantPacks map[int]int // size → count
		wantErr   bool
	}{
		{
			name:      "order 1 ships smallest pack",
			packs:     defaultPacks,
			order:     1,
			wantTotal: 250,
			wantPacks: map[int]int{250: 1},
		},
		{
			name:      "order exactly 250 ships one 250-pack",
			packs:     defaultPacks,
			order:     250,
			wantTotal: 250,
			wantPacks: map[int]int{250: 1},
		},
		{
			name:      "order 251 ships one 500-pack (not two 250s)",
			packs:     defaultPacks,
			order:     251,
			wantTotal: 500,
			wantPacks: map[int]int{500: 1},
		},
		{
			name:      "order 501 ships 500+250",
			packs:     defaultPacks,
			order:     501,
			wantTotal: 750,
			wantPacks: map[int]int{500: 1, 250: 1},
		},
		{
			name:      "order 12001 ships 2×5000+2000+250",
			packs:     defaultPacks,
			order:     12001,
			wantTotal: 12250,
			wantPacks: map[int]int{5000: 2, 2000: 1, 250: 1},
		},
		{
			name:      "order exactly equals a large pack",
			packs:     defaultPacks,
			order:     5000,
			wantTotal: 5000,
			wantPacks: map[int]int{5000: 1},
		},
		{
			name:      "order 500 ships exactly one 500-pack",
			packs:     defaultPacks,
			order:     500,
			wantTotal: 500,
			wantPacks: map[int]int{500: 1},
		},
		{
			name:      "single pack size rounds up",
			packs:     []int{100},
			order:     1,
			wantTotal: 100,
			wantPacks: map[int]int{100: 1},
		},
		{
			name:      "single pack size exact match",
			packs:     []int{100},
			order:     300,
			wantTotal: 300,
			wantPacks: map[int]int{100: 3},
		},
		{
			name:      "non-standard pack sizes",
			packs:     []int{3, 5},
			order:     7,
			wantTotal: 8,
			wantPacks: map[int]int{3: 1, 5: 1},
		},
		{
			name:      "duplicate packs are deduplicated",
			packs:     []int{250, 250, 500},
			order:     251,
			wantTotal: 500,
			wantPacks: map[int]int{500: 1},
		},
		{
			name:    "empty packs returns error",
			packs:   []int{},
			order:   10,
			wantErr: true,
		},
		{
			name:    "nil packs returns error",
			packs:   nil,
			order:   10,
			wantErr: true,
		},
		{
			name:    "zero order returns error",
			packs:   defaultPacks,
			order:   0,
			wantErr: true,
		},
		{
			name:    "negative order returns error",
			packs:   defaultPacks,
			order:   -1,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp, err := algorithm.Calculate(tc.packs, tc.order)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantTotal, resp.TotalItems)

			// Verify total shipped equals sum of (size × count)
			var sum int
			got := make(map[int]int, len(resp.Packs))
			for _, p := range resp.Packs {
				got[p.Size] = p.Count
				sum += p.Size * p.Count
			}
			assert.Equal(t, resp.TotalItems, sum, "sum of packs must equal total_items")
			assert.GreaterOrEqual(t, resp.TotalItems, tc.order, "total must cover the order")
			assert.Equal(t, tc.wantPacks, got)
		})
	}
}
