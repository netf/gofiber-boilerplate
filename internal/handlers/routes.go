package handlers

import (
	"net/http"

	"github.com/netf/gofiber-boilerplate/config"
	"github.com/netf/gofiber-boilerplate/internal/auth"
	"github.com/netf/gofiber-boilerplate/internal/repositories"
	"github.com/netf/gofiber-boilerplate/internal/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB, cfg *config.Config) {
	// JWT Middleware
	authMiddleware := auth.JWTProtected(cfg.JWTSecret)

	// Unprotected routes
	router.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(fiber.Map{"status": "OK"})
	})

	// Todo routes
	registerTodoRoutes(router, db, authMiddleware)
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
