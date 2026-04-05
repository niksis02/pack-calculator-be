package server

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"github.com/niksis02/pack-calculator-be/internal/handler"
	"github.com/niksis02/pack-calculator-be/internal/service"
)

// Server holds the fiber app and its configuration.
type Server struct {
	app          *fiber.App
	port         string
	allowOrigins string
}

// New creates a Server, wires up middlewares and routes, and returns it ready to run.
func New(svc *service.PackService, port, allowOrigins string) *Server {
	s := &Server{
		port:         port,
		allowOrigins: allowOrigins,
	}

	s.app = fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	s.setupMiddlewares()
	s.setupRoutes(handler.NewCalculateHandler(svc), handler.NewConfigHandler(svc))

	return s
}

func (s *Server) setupMiddlewares() {
	s.app.Use(recover.New())
	s.app.Use(logger.New())
	s.app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(s.allowOrigins, ","),
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
	}))
}

func (s *Server) setupRoutes(calcH *handler.CalculateHandler, cfgH *handler.ConfigHandler) {
	s.app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := s.app.Group("/api/v1")
	api.Get("/config/packs", cfgH.GetPacks)
	api.Post("/config/packs", cfgH.SetPacks)
	api.Post("/calculate", calcH.Calculate)
}

// Start begins listening on the configured port. Blocks until the server stops.
func (s *Server) Start() error {
	log.Printf("Starting server on :%s", s.port)
	return s.app.Listen("0.0.0.0:" + s.port)
}

// Shutdown gracefully stops the fiber app.
func (s *Server) Shutdown() error {
	log.Println("Shutting down server...")
	return s.app.Shutdown()
}

// Run starts the server and blocks, handling OS signals for graceful shutdown.
func (s *Server) Run() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-quit
		if err := s.Shutdown(); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
	}()

	if err := s.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
