package handlers

import (
	"net/http"

	"github.com/netf/gofiber-boilerplate/config"
	"github.com/netf/gofiber-boilerplate/internal/auth"
	"github.com/netf/gofiber-boilerplate/internal/repositories"
	"github.com/netf/gofiber-boilerplate/internal/services"
	"github.com/netf/gofiber-boilerplate/internal/validators"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	// Initialize validator
	validators.InitValidator()

	// JWT Middleware
	authMiddleware := auth.JWTProtected(cfg.JWTSecret)

	// API group
	api := app.Group("/api/v1")

	// Unprotected routes
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "OK"})
	})

	// Todo routes
	registerTodoRoutes(api, db, authMiddleware)
}

func registerTodoRoutes(r fiber.Router, db *gorm.DB, authMiddleware fiber.Handler) {
	todoRepo := repositories.NewTodoRepository(db)
	todoService := services.NewTodoService(todoRepo)
	todoHandler := NewTodoHandler(todoService)

	todos := r.Group("/todos", authMiddleware)
	todos.Post("/", todoHandler.CreateTodo)
	todos.Get("/", todoHandler.ListTodos)
	todos.Get("/:id", todoHandler.GetTodoByID)
	todos.Put("/:id", todoHandler.UpdateTodo)
	todos.Delete("/:id", todoHandler.DeleteTodo)
}
