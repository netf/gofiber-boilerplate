package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/netf/gofiber-boilerplate/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockTodoService struct {
	mock.Mock
}

func (m *MockTodoService) CreateTodo(todo *models.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoService) GetTodoByID(id uint) (*models.Todo, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Todo), args.Error(1)
}

func (m *MockTodoService) UpdateTodo(todo *models.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoService) DeleteTodo(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTodoService) ListTodos(page, pageSize int) ([]models.Todo, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]models.Todo), args.Get(1).(int64), args.Error(2)
}

func TestCreateTodo(t *testing.T) {
	mockService := new(MockTodoService)

	// Initialize the validator
	validate := validator.New()

	// Create the handler with both the mock service and the validator
	handler := &TodoHandler{
		service:  mockService,
		validate: validate,
	}

	app := fiber.New()
	app.Post("/todos", handler.CreateTodo)

	todo := models.Todo{Title: "Test Todo", Completed: false}
	mockService.On("CreateTodo", mock.AnythingOfType("*models.Todo")).Return(nil)

	body, _ := json.Marshal(todo)
	req := httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestGetTodoByID(t *testing.T) {
	mockService := new(MockTodoService)
	handler := NewTodoHandler(mockService)

	app := fiber.New()
	app.Get("/todos/:id", handler.GetTodoByID)

	todo := &models.Todo{ID: 1, Title: "Test Todo", Completed: false}
	mockService.On("GetTodoByID", uint(1)).Return(todo, nil)

	req := httptest.NewRequest("GET", "/todos/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestGetTodoByIDNotFound(t *testing.T) {
	mockService := new(MockTodoService)
	handler := NewTodoHandler(mockService)

	app := fiber.New()
	app.Get("/todos/:id", handler.GetTodoByID)

	mockService.On("GetTodoByID", uint(1)).Return(&models.Todo{}, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/todos/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestUpdateTodo(t *testing.T) {
	mockService := new(MockTodoService)
	handler := NewTodoHandler(mockService)

	app := fiber.New()
	app.Put("/todos/:id", handler.UpdateTodo)

	todo := models.Todo{ID: 1, Title: "Updated Todo", Completed: true}
	mockService.On("UpdateTodo", mock.AnythingOfType("*models.Todo")).Return(nil)

	body, _ := json.Marshal(todo)
	req := httptest.NewRequest("PUT", "/todos/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestDeleteTodo(t *testing.T) {
	mockService := new(MockTodoService)
	handler := NewTodoHandler(mockService)

	app := fiber.New()
	app.Delete("/todos/:id", handler.DeleteTodo)

	mockService.On("DeleteTodo", uint(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/todos/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestListTodos(t *testing.T) {
	mockService := new(MockTodoService)
	handler := NewTodoHandler(mockService)

	app := fiber.New()
	app.Get("/api/v1/todos", handler.ListTodos)

	testCases := []struct {
		name           string
		query          string
		expectedStatus int
		mockTodos      []models.Todo
		mockTotal      int64
		mockError      error
		setupMock      func(*MockTodoService)
	}{
		{
			name:           "Success - Default Pagination",
			query:          "",
			expectedStatus: fiber.StatusOK,
			mockTodos: []models.Todo{
				{ID: 1, Title: "Todo 1", Completed: false},
				{ID: 2, Title: "Todo 2", Completed: true},
			},
			mockTotal: 2,
			mockError: nil,
			setupMock: func(m *MockTodoService) {
				m.On("ListTodos", 1, 10).Return([]models.Todo{
					{ID: 1, Title: "Todo 1", Completed: false},
					{ID: 2, Title: "Todo 2", Completed: true},
				}, int64(2), nil)
			},
		},
		{
			name:           "Success - Custom Pagination",
			query:          "?page=2&page_size=5",
			expectedStatus: fiber.StatusOK,
			mockTodos: []models.Todo{
				{ID: 6, Title: "Todo 6", Completed: false},
				{ID: 7, Title: "Todo 7", Completed: true},
			},
			mockTotal: 7,
			mockError: nil,
			setupMock: func(m *MockTodoService) {
				m.On("ListTodos", 2, 5).Return([]models.Todo{
					{ID: 6, Title: "Todo 6", Completed: false},
					{ID: 7, Title: "Todo 7", Completed: true},
				}, int64(7), nil)
			},
		},
		{
			name:           "Error - Invalid Page",
			query:          "?page=0",
			expectedStatus: fiber.StatusBadRequest,
			setupMock:      func(m *MockTodoService) {}, // No mock setup needed for validation error
		},
		{
			name:           "Error - Invalid Page Size",
			query:          "?page_size=101",
			expectedStatus: fiber.StatusBadRequest,
			setupMock:      func(m *MockTodoService) {}, // No mock setup needed for validation error
		},
		{
			name:           "Error - Service Failure",
			query:          "",
			expectedStatus: fiber.StatusInternalServerError,
			setupMock: func(m *MockTodoService) {
				m.On("ListTodos", 1, 10).Return([]models.Todo{}, int64(0), errors.New("service error"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil
			mockService.Calls = nil

			tc.setupMock(mockService)

			req := httptest.NewRequest("GET", "/api/v1/todos"+tc.query, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err, "Failed to perform request")

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.expectedStatus == fiber.StatusOK {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err, "Failed to read response body")

				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err, "Failed to unmarshal response body")

				t.Logf("Response body: %s", string(body))
				t.Logf("Unmarshaled result: %+v", result)

				data, ok := result["data"].([]interface{})
				assert.True(t, ok, "Data should be a slice of interfaces")
				assert.NotEmpty(t, data, "Data should not be empty")
				assert.Len(t, data, len(tc.mockTodos), "Data length mismatch")

				for i, todoInterface := range data {
					todo, ok := todoInterface.(map[string]interface{})
					assert.True(t, ok, "Each item in data should be a map")
					if ok && i < len(tc.mockTodos) {
						assert.Equal(t, float64(tc.mockTodos[i].ID), todo["id"], "ID mismatch")
						assert.Equal(t, tc.mockTodos[i].Title, todo["title"], "Title mismatch")
						assert.Equal(t, tc.mockTodos[i].Completed, todo["completed"], "Completed status mismatch")
					}
				}

				// Check for total items in the response headers or body
				totalItems, ok := result["total_items"].(float64)
				if ok {
					assert.Equal(t, float64(tc.mockTotal), totalItems, "Total items mismatch")
				} else {
					t.Log("Total items not found in response body")
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestCreateTodoError(t *testing.T) {
	mockService := new(MockTodoService)
	handler := NewTodoHandler(mockService)

	app := fiber.New()
	app.Post("/todos", handler.CreateTodo)

	todo := models.Todo{Title: "Test Todo", Completed: false}
	mockService.On("CreateTodo", mock.AnythingOfType("*models.Todo")).Return(errors.New("database error"))

	body, _ := json.Marshal(todo)
	req := httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockService.AssertExpectations(t)
}
