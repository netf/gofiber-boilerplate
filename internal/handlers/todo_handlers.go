package handlers

import (
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/netf/gofiber-boilerplate/internal/models"
	"github.com/netf/gofiber-boilerplate/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Package handlers contains the HTTP handlers for the API
// @title Fiber API
// @version 1.0
// @description This is a sample server Fiber server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1

type TodoHandler struct {
	service  services.TodoService
	validate *validator.Validate
}

func NewTodoHandler(service services.TodoService) *TodoHandler {
	return &TodoHandler{
		service:  service,
		validate: validator.New(),
	}
}

// CreateTodo creates a new todo item
// @Summary Create a new todo
// @Tags Todos
// @Accept json
// @Produce json
// @Param todo body models.Todo true "Todo item"
// @Success 201 {object} models.Todo
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /todos [post]
// @Security ApiKeyAuth
func (h *TodoHandler) CreateTodo(c *fiber.Ctx) error {
	var todo models.Todo
	if err := c.BodyParser(&todo); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Validate the request
	if err := h.validate.Struct(&todo); err != nil {
		log.Warn().Err(err).Msg("Validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.service.CreateTodo(&todo); err != nil {
		log.Error().Err(err).Msg("Failed to create todo")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create todo",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(todo)
}

// GetTodoByID retrieves a todo item by ID
// @Summary Get a todo by ID
// @Tags Todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 400 {object} fiber.Map
// @Failure 404 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /todos/{id} [get]
// @Security ApiKeyAuth
func (h *TodoHandler) GetTodoByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Warn().Msg("Invalid ID parameter")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	todo, err := h.service.GetTodoByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Uint("id", uint(id)).Msg("Todo not found")
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Todo not found",
			})
		}
		log.Error().Err(err).Msg("Failed to retrieve todo")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve todo",
		})
	}

	return c.JSON(todo)
}

// UpdateTodo updates an existing todo item
// @Summary Update a todo
// @Tags Todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body models.Todo true "Todo item"
// @Success 200 {object} models.Todo
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /todos/{id} [put]
// @Security ApiKeyAuth
func (h *TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Warn().Msg("Invalid ID parameter")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var todo models.Todo
	if err := c.BodyParser(&todo); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Validate the request
	if err := h.validate.Struct(&todo); err != nil {
		log.Warn().Err(err).Msg("Validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	todo.ID = uint(id)
	if err := h.service.UpdateTodo(&todo); err != nil {
		log.Error().Err(err).Msg("Failed to update todo")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update todo",
		})
	}

	return c.JSON(todo)
}

// DeleteTodo deletes a todo item
// @Summary Delete a todo
// @Tags Todos
// @Param id path int true "Todo ID"
// @Success 204
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /todos/{id} [delete]
// @Security ApiKeyAuth
func (h *TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Warn().Msg("Invalid ID parameter")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := h.service.DeleteTodo(uint(id)); err != nil {
		log.Error().Err(err).Msg("Failed to delete todo")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete todo",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListTodos retrieves all todo items
// @Summary Get all todos
// @Tags Todos
// @Produce json
// @Success 200 {array} models.Todo
// @Failure 500 {object} fiber.Map
// @Router /todos [get]
// @Security ApiKeyAuth
func (h *TodoHandler) ListTodos(c *fiber.Ctx) error {
	todos, err := h.service.ListTodos()
	if err != nil {
		log.Error().Err(err).Msg("Failed to list todos")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list todos",
		})
	}

	return c.JSON(todos)
}
