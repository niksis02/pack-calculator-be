package service

import (
	"errors"
	"sync"

	"github.com/niksis02/pack-calculator-be/internal/algorithm"
	"github.com/niksis02/pack-calculator-be/internal/model"
)

// PackService manages pack configuration and calculation.
// It is safe for concurrent use.
type PackService struct {
	mu    sync.RWMutex
	packs []int
}

// NewPackService returns a PackService seeded with the given default pack sizes.
func NewPackService(defaultPacks []int) *PackService {
	cp := make([]int, len(defaultPacks))
	copy(cp, defaultPacks)
	return &PackService{packs: cp}
}

// GetPacks returns the current pack configuration (thread-safe).
func (s *PackService) GetPacks() model.PackConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cp := make([]int, len(s.packs))
	copy(cp, s.packs)
	return model.PackConfig{Packs: cp}
}

// SetPacks validates and replaces the pack configuration (thread-safe).
// All values must be positive integers; at least one must be provided.
func (s *PackService) SetPacks(packs []int) error {
	if len(packs) == 0 {
		return errors.New("pack list must not be empty")
	}
	for _, p := range packs {
		if p <= 0 {
			return errors.New("all pack sizes must be positive integers")
		}
	}
	cp := make([]int, len(packs))
	copy(cp, packs)
	s.mu.Lock()
	s.packs = cp
	s.mu.Unlock()
	return nil
}

// Calculate computes the optimal packing for an order against the current config.
func (s *PackService) Calculate(order int) (model.CalculateResponse, error) {
	s.mu.RLock()
	packs := make([]int, len(s.packs))
	copy(packs, s.packs)
	s.mu.RUnlock()

	return algorithm.Calculate(packs, order)
}
