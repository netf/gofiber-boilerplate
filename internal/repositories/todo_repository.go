package repositories

import (
	"github.com/netf/gofiber-boilerplate/internal/models"

	"gorm.io/gorm"
)

type TodoRepository interface {
	Create(todo *models.Todo) error
	GetByID(id uint) (*models.Todo, error)
	Update(todo *models.Todo) error
	Delete(id uint) error
	List(page, pageSize int) ([]models.Todo, int64, error)
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db}
}

func (r *todoRepository) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

func (r *todoRepository) GetByID(id uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.First(&todo, id).Error
	return &todo, err
}

func (r *todoRepository) Update(todo *models.Todo) error {
	return r.db.Save(todo).Error
}

func (r *todoRepository) Delete(id uint) error {
	return r.db.Delete(&models.Todo{}, id).Error
}

func (r *todoRepository) List(page, pageSize int) ([]models.Todo, int64, error) {
	var todos []models.Todo
	var total int64

	offset := (page - 1) * pageSize

	err := r.db.Model(&models.Todo{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(pageSize).Find(&todos).Error
	return todos, total, err
}
