package datastore

import (
	"context"
	"errors"

	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/gateway"
	"github.com/phamquanandpad/training-project/go/services/todo/internal/domain/model/todo"
	"gorm.io/gorm"
)

type TodoWriter struct{}

func NewTodoWriter() gateway.TodoCommandsGateway {
	return &TodoWriter{}
}

func (w *TodoWriter) CreateTodo(
	ctx context.Context,
	newTodo todo.NewTodo,
) (*todo.Todo, error) {
	tx, err := ExtractTodoDB(ctx)
	if err != nil {
		return nil, err
	}

	db := tx.WithContext(ctx)
	createdTodo := todo.Todo{
		UserID:      newTodo.UserID,
		Task:        newTodo.Task,
		Description: newTodo.Description,
		Status:      newTodo.Status,
	}

	if err := db.
		Create(&createdTodo).
		Error; err != nil {
		return nil, err
	}
	return &createdTodo, nil
}

func (w *TodoWriter) UpdateTodo(
	ctx context.Context,
	todoID todo.TodoID,
	updateTodo todo.UpdateTodo,
) (*todo.Todo, error) {
	tx, err := ExtractTodoDB(ctx)
	if err != nil {
		return nil, err
	}

	db := tx.WithContext(ctx)

	var todo todo.Todo
	if err := db.
		Where("id = ? AND deleted_at IS NULL", todoID).
		First(&todo).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	todo.Description = updateTodo.Description
	if updateTodo.Task != nil {
		todo.Task = *updateTodo.Task
	}
	if updateTodo.Status != nil {
		todo.Status = *updateTodo.Status
	}

	if err := db.Save(&todo).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

func (w *TodoWriter) SoftDeleteTodo(
	ctx context.Context,
	todoID todo.TodoID,
) error {
	tx, err := ExtractTodoDB(ctx)
	if err != nil {
		return err
	}

	db := tx.WithContext(ctx)

	if err := db.
		Where("id = ? AND deleted_at IS NULL", todoID).
		Delete(&todo.Todo{}).
		Error; err != nil {
		return err
	}
	return nil
}
