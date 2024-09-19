package handlers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	apiUtils "github.com/netf/gofiber-boilerplate/internal/api/utils"
	"github.com/netf/gofiber-boilerplate/internal/models"
	"github.com/netf/gofiber-boilerplate/internal/services"
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
// @Success 201 {object} apiUtils.Response[models.Todo]
// @Failure 400 {object} apiUtils.ErrorResponse
// @Failure 500 {object} apiUtils.ErrorResponse
// @Router /todos [post]
// @Security ApiKeyAuth
func (h *TodoHandler) CreateTodo(c *fiber.Ctx) error {
	var todo models.Todo
	if err := c.BodyParser(&todo); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		errorResponse := apiUtils.CreateErrorResponse("Cannot parse JSON", fiber.StatusBadRequest)
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	if err := h.validate.Struct(&todo); err != nil {
		log.Warn().Err(err).Msg("Validation failed")
		errorResponse := apiUtils.CreateErrorResponse(err.Error(), fiber.StatusBadRequest)
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	if err := h.service.CreateTodo(&todo); err != nil {
		log.Error().Err(err).Msg("Failed to create todo")
		errorResponse := apiUtils.CreateErrorResponse("Failed to create todo", fiber.StatusInternalServerError)
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	response := apiUtils.CreateResponse[models.Todo](todo)
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetTodoByID retrieves a todo item by ID
// @Summary Get a todo by ID
// @Tags Todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} apiUtils.Response[models.Todo]
// @Failure 400 {object} apiUtils.ErrorResponse
// @Failure 404 {object} apiUtils.ErrorResponse
// @Failure 500 {object} apiUtils.ErrorResponse
// @Router /todos/{id} [get]
// @Security ApiKeyAuth
func (h *TodoHandler) GetTodoByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Warn().Msg("Invalid ID parameter")
		errorResponse := apiUtils.CreateErrorResponse("Invalid ID", fiber.StatusBadRequest)
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	todo, err := h.service.GetTodoByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Uint("id", uint(id)).Msg("Todo not found")
			errorResponse := apiUtils.CreateErrorResponse("Todo not found", fiber.StatusNotFound)
			return c.Status(fiber.StatusNotFound).JSON(errorResponse)
		}
		log.Error().Err(err).Msg("Failed to retrieve todo")
		errorResponse := apiUtils.CreateErrorResponse("Failed to retrieve todo", fiber.StatusInternalServerError)
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	response := apiUtils.CreateResponse[models.Todo](todo)
	return c.JSON(response)
}

// UpdateTodo updates an existing todo item
// @Summary Update a todo
// @Tags Todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body models.Todo true "Todo item"
// @Success 200 {object} apiUtils.Response[models.Todo]
// @Failure 400 {object} apiUtils.ErrorResponse
// @Failure 500 {object} apiUtils.ErrorResponse
// @Router /todos/{id} [put]
// @Security ApiKeyAuth
func (h *TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Warn().Msg("Invalid ID parameter")
		errorResponse := apiUtils.CreateErrorResponse("Invalid ID", fiber.StatusBadRequest)
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	var todo models.Todo
	if err := c.BodyParser(&todo); err != nil {
		log.Error().Err(err).Msg("Failed to parse request body")
		errorResponse := apiUtils.CreateErrorResponse("Cannot parse JSON", fiber.StatusBadRequest)
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	if err := h.validate.Struct(&todo); err != nil {
		log.Warn().Err(err).Msg("Validation failed")
		errorResponse := apiUtils.CreateErrorResponse(err.Error(), fiber.StatusBadRequest)
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	todo.ID = uint(id)
	if err := h.service.UpdateTodo(&todo); err != nil {
		log.Error().Err(err).Msg("Failed to update todo")
		errorResponse := apiUtils.CreateErrorResponse("Failed to update todo", fiber.StatusInternalServerError)
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	response := apiUtils.CreateResponse[models.Todo](todo)
	return c.JSON(response)
}

// DeleteTodo deletes a todo item
// @Summary Delete a todo
// @Tags Todos
// @Param id path int true "Todo ID"
// @Success 204 "No Content"
// @Failure 400 {object} apiUtils.ErrorResponse
// @Failure 500 {object} apiUtils.ErrorResponse
// @Router /todos/{id} [delete]
// @Security ApiKeyAuth
func (h *TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Warn().Msg("Invalid ID parameter")
		errorResponse := apiUtils.CreateErrorResponse("Invalid ID", fiber.StatusBadRequest)
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	if err := h.service.DeleteTodo(uint(id)); err != nil {
		log.Error().Err(err).Msg("Failed to delete todo")
		errorResponse := apiUtils.CreateErrorResponse("Failed to delete todo", fiber.StatusInternalServerError)
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListTodos retrieves all todo items with pagination
// @Summary Get all todos
// @Description Get a paginated list of todos
// @Tags Todos
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1)
// @Param page_size query int false "Page size" default(10) minimum(1) maximum(100)
// @Success 200 {object} apiUtils.Response[[]models.Todo]
// @Failure 400 {object} apiUtils.ErrorResponse
// @Failure 500 {object} apiUtils.ErrorResponse
// @Router /todos [get]
// @Security ApiKeyAuth
func (h *TodoHandler) ListTodos(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)

	// Validate page and page_size
	if page < 1 {
		errorResponse := apiUtils.CreateErrorResponse("Invalid page number", fiber.StatusBadRequest)
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}
	if pageSize < 1 || pageSize > 100 {
		errorResponse := apiUtils.CreateErrorResponse("Invalid page size", fiber.StatusBadRequest)
		return c.Status(fiber.StatusBadRequest).JSON(errorResponse)
	}

	todos, total, err := h.service.ListTodos(page, pageSize)
	if err != nil {
		errorResponse := apiUtils.CreateErrorResponse("Failed to fetch todos", fiber.StatusInternalServerError)
		return c.Status(fiber.StatusInternalServerError).JSON(errorResponse)
	}

	response := apiUtils.CreateResponse[[]models.Todo](todos, page, pageSize, int(total))
	return c.JSON(response)
}
