package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/netf/gofiber-boilerplate/config"
	"github.com/netf/gofiber-boilerplate/internal/api/handlers"
	"github.com/netf/gofiber-boilerplate/internal/repositories"
	"github.com/netf/gofiber-boilerplate/internal/services"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB, cfg *config.Config) {
	todoRepo := repositories.NewTodoRepository(db)
	todoService := services.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	todoRoutes := router.Group("/todos")
	todoRoutes.Post("/", todoHandler.CreateTodo)
	todoRoutes.Get("/", todoHandler.ListTodos)
	todoRoutes.Get("/:id", todoHandler.GetTodoByID)
	todoRoutes.Put("/:id", todoHandler.UpdateTodo)
	todoRoutes.Delete("/:id", todoHandler.DeleteTodo)
}
