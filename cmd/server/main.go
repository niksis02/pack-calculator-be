package main

import (
	"github.com/niksis02/pack-calculator-be/internal/config"
	"github.com/niksis02/pack-calculator-be/internal/server"
	"github.com/niksis02/pack-calculator-be/internal/service"
)

func main() {
	cfg := config.Load()

	svc := service.NewPackService([]int{250, 500, 1000, 2000, 5000})
	server.New(svc, cfg.Port, cfg.AllowOrigins).Run()
}
