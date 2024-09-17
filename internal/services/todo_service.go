package services

import (
	"github.com/netf/gofiber-boilerplate/internal/models"
	"github.com/netf/gofiber-boilerplate/internal/repositories"
)

type TodoService interface {
	CreateTodo(todo *models.Todo) error
	GetTodoByID(id uint) (*models.Todo, error)
	UpdateTodo(todo *models.Todo) error
	DeleteTodo(id uint) error
	ListTodos() ([]models.Todo, error)
}

type todoService struct {
	repo repositories.TodoRepository
}

func NewTodoService(repo repositories.TodoRepository) TodoService {
	return &todoService{repo}
}

func (s *todoService) CreateTodo(todo *models.Todo) error {
	return s.repo.Create(todo)
}

func (s *todoService) GetTodoByID(id uint) (*models.Todo, error) {
	return s.repo.GetByID(id)
}

func (s *todoService) UpdateTodo(todo *models.Todo) error {
	return s.repo.Update(todo)
}

func (s *todoService) DeleteTodo(id uint) error {
	return s.repo.Delete(id)
}

func (s *todoService) ListTodos() ([]models.Todo, error) {
	return s.repo.List()
}
