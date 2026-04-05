package algorithm

import (
	"errors"
	"sort"

	"github.com/niksis02/pack-calculator-be/internal/model"
)

// ErrNoPackSizes is returned when the pack list is empty.
var ErrNoPackSizes = errors.New("at least one pack size is required")

// ErrInvalidOrder is returned for non-positive order quantities.
var ErrInvalidOrder = errors.New("order quantity must be >= 1")

const unreachable = -1

// Calculate returns the optimal packing for the given order quantity.
//
// Constraints:
//  1. Only whole packs may be used.
//  2. Total items shipped >= order (primary: minimize total shipped).
//  3. Minimize number of packs (secondary tie-breaker).
//
// Algorithm: unbounded DP over exact pack sums up to order + max(packSizes).
// dp[q] stores the minimum number of packs required to produce EXACTLY q items.
// The first reachable q >= order gives the minimum total; pack count is already
// minimized by the DP itself.
func Calculate(packSizes []int, order int) (model.CalculateResponse, error) {
	if len(packSizes) == 0 {
		return model.CalculateResponse{}, ErrNoPackSizes
	}
	if order < 1 {
		return model.CalculateResponse{}, ErrInvalidOrder
	}

	packs := dedupSorted(packSizes)
	maxPack := packs[len(packs)-1]
	upper := order + maxPack // tightest possible upper bound

	// dp[q] = minimum number of packs to make exactly q items; unreachable = -1
	dp := make([]int, upper+1)
	for i := range dp {
		dp[i] = unreachable
	}
	dp[0] = 0

	for q := 1; q <= upper; q++ {
		for _, p := range packs {
			if p > q {
				// packs are sorted ascending; none of the remaining can help
				break
			}
			prev := dp[q-p]
			if prev == unreachable {
				continue
			}
			candidate := prev + 1
			if dp[q] == unreachable || candidate < dp[q] {
				dp[q] = candidate
			}
		}
	}

	// Find the first reachable quantity >= order.
	// Because total == q in an exact-coverage model, the first valid q gives
	// the minimum total shipped; dp[q] (pack count) is already minimised.
	bestQ := unreachable
	for q := order; q <= upper; q++ {
		if dp[q] != unreachable {
			bestQ = q
			break
		}
	}

	if bestQ == unreachable {
		// Unreachable in theory if packs are valid — the single largest pack
		// always covers [order, order+maxPack].
		return model.CalculateResponse{}, errors.New("no valid packing found")
	}

	packCounts := reconstructPacks(dp, packs, bestQ)

	// Sort result descending by size for a natural display order.
	sort.Slice(packCounts, func(i, j int) bool {
		return packCounts[i].Size > packCounts[j].Size
	})

	return model.CalculateResponse{
		TotalItems: bestQ,
		Packs:      packCounts,
	}, nil
}

// reconstructPacks walks the dp table backwards to identify which packs were used.
// It always picks the smallest valid pack at each step (packs are sorted ascending),
// which is deterministic and produces a canonical decomposition.
func reconstructPacks(dp []int, packs []int, q int) []model.PackResult {
	counts := make(map[int]int)
	for q > 0 {
		for _, p := range packs {
			if p > q {
				continue
			}
			prev := dp[q-p]
			if prev == unreachable {
				continue
			}
			if dp[q] == prev+1 {
				counts[p]++
				q -= p
				break
			}
		}
	}
	result := make([]model.PackResult, 0, len(counts))
	for size, count := range counts {
		result = append(result, model.PackResult{Size: size, Count: count})
	}
	return result
}

// dedupSorted returns a sorted, deduplicated copy of the input slice,
// filtering out any non-positive values.
func dedupSorted(in []int) []int {
	seen := make(map[int]struct{}, len(in))
	out := make([]int, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; !ok && v > 0 {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	sort.Ints(out)
	return out
}
