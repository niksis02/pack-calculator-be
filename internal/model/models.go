package model

// PackConfig holds the currently active pack sizes.
type PackConfig struct {
	Packs []int `json:"packs"`
}

// CalculateRequest is the body for POST /calculate.
type CalculateRequest struct {
	Items int `json:"items"`
}

// PackResult describes how many of a given pack size to ship.
type PackResult struct {
	Size  int `json:"size"`
	Count int `json:"count"`
}

// CalculateResponse is returned from POST /calculate.
type CalculateResponse struct {
	TotalItems int          `json:"total_items"`
	Packs      []PackResult `json:"packs"`
}
