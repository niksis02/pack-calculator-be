package service_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/niksis02/pack-calculator-be/internal/service"
)

func TestGetPacks_ReturnsDefaults(t *testing.T) {
	t.Parallel()
	svc := service.NewPackService([]int{250, 500, 1000})
	cfg := svc.GetPacks()
	assert.Equal(t, []int{250, 500, 1000}, cfg.Packs)
}

func TestGetPacks_ReturnsCopy(t *testing.T) {
	t.Parallel()
	svc := service.NewPackService([]int{250})
	cfg := svc.GetPacks()
	cfg.Packs[0] = 999
	assert.Equal(t, []int{250}, svc.GetPacks().Packs, "mutating returned slice must not affect service state")
}

func TestSetPacks_ReplacesConfig(t *testing.T) {
	t.Parallel()
	svc := service.NewPackService([]int{250, 500})
	require.NoError(t, svc.SetPacks([]int{100, 200, 300}))
	assert.Equal(t, []int{100, 200, 300}, svc.GetPacks().Packs)
}

func TestSetPacks_RejectsEmpty(t *testing.T) {
	t.Parallel()
	svc := service.NewPackService([]int{250})
	assert.Error(t, svc.SetPacks([]int{}))
}

func TestSetPacks_RejectsZero(t *testing.T) {
	t.Parallel()
	svc := service.NewPackService([]int{250})
	assert.Error(t, svc.SetPacks([]int{0, 500}))
}

func TestSetPacks_RejectsNegative(t *testing.T) {
	t.Parallel()
	svc := service.NewPackService([]int{250})
	assert.Error(t, svc.SetPacks([]int{-1, 500}))
}

func TestCalculate_Basic(t *testing.T) {
	t.Parallel()
	svc := service.NewPackService([]int{250, 500, 1000, 2000, 5000})
	resp, err := svc.Calculate(251)
	require.NoError(t, err)
	assert.Equal(t, 500, resp.TotalItems)
}

func TestCalculate_ConcurrentSafety(t *testing.T) {
	t.Parallel()
	svc := service.NewPackService([]int{250, 500, 1000})

	var wg sync.WaitGroup
	const goroutines = 50

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			switch n % 3 {
			case 0:
				_ = svc.SetPacks([]int{100, 200, 500})
			case 1:
				_ = svc.GetPacks()
			default:
				_, _ = svc.Calculate(300)
			}
		}(i)
	}
	wg.Wait()
}
