package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/niksis02/pack-calculator-be/internal/handler"
	"github.com/niksis02/pack-calculator-be/internal/model"
	"github.com/niksis02/pack-calculator-be/internal/service"
)

func newTestApp() *fiber.App {
	svc := service.NewPackService([]int{250, 500, 1000, 2000, 5000})
	calcH := handler.NewCalculateHandler(svc)
	cfgH := handler.NewConfigHandler(svc)

	app := fiber.New()
	api := app.Group("/api/v1")
	api.Get("/config/packs", cfgH.GetPacks)
	api.Post("/config/packs", cfgH.SetPacks)
	api.Post("/calculate", calcH.Calculate)
	return app
}

func TestGetPacks_ReturnsDefaults(t *testing.T) {
	t.Parallel()
	app := newTestApp()

	req := httptest.NewRequest("GET", "/api/v1/config/packs", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body model.PackConfig
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, []int{250, 500, 1000, 2000, 5000}, body.Packs)
}

func TestSetPacks_UpdatesConfig(t *testing.T) {
	t.Parallel()
	app := newTestApp()

	payload, _ := json.Marshal(model.PackConfig{Packs: []int{100, 500}})
	req := httptest.NewRequest("POST", "/api/v1/config/packs", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var body model.PackConfig
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, []int{100, 500}, body.Packs)
}

func TestSetPacks_EmptyArray_Returns400(t *testing.T) {
	t.Parallel()
	app := newTestApp()

	payload, _ := json.Marshal(model.PackConfig{Packs: []int{}})
	req := httptest.NewRequest("POST", "/api/v1/config/packs", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestSetPacks_NegativeValue_Returns400(t *testing.T) {
	t.Parallel()
	app := newTestApp()

	payload, _ := json.Marshal(model.PackConfig{Packs: []int{-1, 500}})
	req := httptest.NewRequest("POST", "/api/v1/config/packs", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestSetPacks_BadJSON_Returns400(t *testing.T) {
	t.Parallel()
	app := newTestApp()

	req := httptest.NewRequest("POST", "/api/v1/config/packs", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestCalculate_ReturnsOptimalResult(t *testing.T) {
	t.Parallel()
	app := newTestApp()

	payload, _ := json.Marshal(model.CalculateRequest{Items: 251})
	req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result model.CalculateResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(t, 500, result.TotalItems)
	require.Len(t, result.Packs, 1)
	assert.Equal(t, 500, result.Packs[0].Size)
	assert.Equal(t, 1, result.Packs[0].Count)
}

func TestCalculate_ZeroItems_Returns400(t *testing.T) {
	t.Parallel()
	app := newTestApp()

	payload, _ := json.Marshal(model.CalculateRequest{Items: 0})
	req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestCalculate_BadJSON_Returns400(t *testing.T) {
	t.Parallel()
	app := newTestApp()

	req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBufferString("{bad}"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestCalculate_AfterConfigUpdate(t *testing.T) {
	t.Parallel()
	app := newTestApp()

	// Update pack config
	cfgPayload, _ := json.Marshal(model.PackConfig{Packs: []int{100, 500}})
	cfgReq := httptest.NewRequest("POST", "/api/v1/config/packs", bytes.NewReader(cfgPayload))
	cfgReq.Header.Set("Content-Type", "application/json")
	cfgResp, err := app.Test(cfgReq)
	require.NoError(t, err)
	assert.Equal(t, 200, cfgResp.StatusCode)

	// Calculate against updated config
	calcPayload, _ := json.Marshal(model.CalculateRequest{Items: 101})
	calcReq := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewReader(calcPayload))
	calcReq.Header.Set("Content-Type", "application/json")
	calcResp, err := app.Test(calcReq)
	require.NoError(t, err)
	assert.Equal(t, 200, calcResp.StatusCode)

	var result model.CalculateResponse
	require.NoError(t, json.NewDecoder(calcResp.Body).Decode(&result))
	// packs=[100,500], order=101 → 2×100=200 is less total than 1×500=500 (primary objective)
	assert.Equal(t, 200, result.TotalItems)
}
